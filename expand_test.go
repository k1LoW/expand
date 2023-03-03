package expand

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExpandYAML(t *testing.T) {
	tests := []struct {
		in   string
		envs map[string]string
		want string
	}{
		{
			`key: value
key2: value2
`,
			map[string]string{},
			`key: value
key2: value2
`},
		{
			`key: value
key2: value2`,
			map[string]string{},
			`key: value
key2: value2`},
		{
			`default: "hello ${UNDEFINED:-world}"
multi: |

  hello world

  hello ${WORLD}
`,
			map[string]string{
				"WORLD": ": world :world",
			},
			`default: "hello world"
multi: |

  hello world

  hello : world :world
`},
		{
			`coverage:
  acceptable: ${COVERAGE_ACCEPTABLE}
  badge:
    path: ${COVERAGE_BADGE_PATH}
comment:
  enable: ${COMMENT_ENABLE}
`,
			map[string]string{
				"COVERAGE_ACCEPTABLE": "60%",
				"COVERAGE_BADGE_PATH": "path/to/coverage.svg",
				"COMMENT_ENABLE":      "true",
			},
			`coverage:
  acceptable: 60%
  badge:
    path: path/to/coverage.svg
comment:
  enable: true
`},
		{
			`key: ${VALUE}
${KEY}: value
`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			`key: envvalue
${KEY}: value
`},
		{
			`string: |
         hello world
         hello ${WORLD}
array:
  - hello world
  - |
    hello ${WORLD}
`,
			map[string]string{
				"WORLD": "world",
			},
			`string: |
         hello world
         hello world
array:
  - hello world
  - |
    hello world
`},
		{
			`string: hello ${WORLD}
multi: |

  hello world

  hello ${WORLD}
`,
			map[string]string{
				"WORLD": ": world :world",
			},
			`string: "hello : world :world"
multi: |

  hello world

  hello : world :world
`},
		{
			`test: |
  current.url == 'https://example.com/#about'
`,
			map[string]string{},
			`test: |
  current.url == 'https://example.com/#about'
`},
		{
			`key: "hello$"`,
			map[string]string{},
			`key: "hello$"`,
		},
		{
			`key: "hello$ ${WORLD}"`,
			map[string]string{
				"WORLD": "world",
			},
			`key: "hello$ world"`,
		},
	}
	for _, tt := range tests {
		for k, v := range tt.envs {
			t.Setenv(k, v)
		}
		mapper := func(in string) (string, bool) {
			return os.LookupEnv(in)
		}
		got := ExpandYAML(tt.in, mapper)
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}

func TestReplaceYAML(t *testing.T) {
	tests := []struct {
		in            string
		envs          map[string]string
		replaceMapKey bool
		want          string
	}{
		{
			`key: ${VALUE}
${KEY}: value
`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			false,
			`key: envvalue
${KEY}: value
`,
		},
		{
			`key: ${VALUE}
${KEY}: value
`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			true,
			`key: envvalue
envkey: value
`,
		},
		{
			`key: "${VALUE}-${VALUE}"
`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			true,
			`key: "envvalue-envvalue"
`,
		},
		{
			`key: "${VALUE}\n${VALUE}"
`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			true,
			`key: "envvalue\nenvvalue"
`,
		},
		{
			`key: "${VALUE}\n${VALUE}"`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			true,
			`key: "envvalue\nenvvalue"`,
		},
		{
			`key: |
  ${VALUE}
  ${VALUE}
`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			true,
			`key: |
  envvalue
  envvalue
`,
		},
		{
			`key: |
  ${VALUE}
  ${VALUE}`,
			map[string]string{
				"KEY":   "envkey",
				"VALUE": "envvalue",
			},
			true,
			`key: |
  envvalue
  envvalue`,
		},
	}
	repFn := func(in string) (string, error) {
		return os.Expand(in, os.Getenv), nil
	}
	for _, tt := range tests {
		for k, v := range tt.envs {
			t.Setenv(k, v)
		}
		got, err := ReplaceYAML(tt.in, repFn, tt.replaceMapKey)
		if err != nil {
			t.Error(err)
		}
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}
