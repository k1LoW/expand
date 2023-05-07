package expand

import (
	"fmt"
	"os"
	"testing"

	"github.com/goccy/go-yaml"
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
		{
			`key: ${KEY}
port: ${PORT}`,
			map[string]string{
				"KEY":  "hello\nworld",
				"PORT": "2202",
			},
			`key: "hello\nworld"
port: 2202`,
		},
		{
			`port: ${PORT}
key: ${KEY}`,
			map[string]string{
				"KEY":  "hello\nworld",
				"PORT": "2202",
			},
			`port: 2202
key: "hello\nworld"`,
		},
		{
			`key: '${KEY}'
port: ${PORT}`,
			map[string]string{
				"KEY":  "hello\nworld",
				"PORT": "2202",
			},
			`key: 'hello
world'
port: 2202`,
		},
		{
			`port: ${PORT}
key: '${KEY}'`,
			map[string]string{
				"KEY":  "hello\nworld",
				"PORT": "2202",
			},
			`port: 2202
key: 'hello
world'`,
		},
		{
			`key: "${KEY}"
port: ${PORT}`,
			map[string]string{
				"KEY":  "hello\nworld",
				"PORT": "2202",
			},
			`key: "hello\nworld"
port: 2202`,
		},
		{
			`port: ${PORT}
key: "${KEY}"`,
			map[string]string{
				"KEY":  "hello\nworld",
				"PORT": "2202",
			},
			`port: 2202
key: "hello\nworld"`,
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
		var v interface{}
		if err := yaml.Unmarshal([]byte(got), &v); err != nil {
			t.Errorf("%s", err)
		}
		if err := yaml.Unmarshal([]byte(tt.want), &v); err != nil {
			t.Errorf("%s", err)
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
		var v interface{}
		if err := yaml.Unmarshal([]byte(got), &v); err != nil {
			t.Errorf("%s", err)
		}
		if err := yaml.Unmarshal([]byte(tt.want), &v); err != nil {
			t.Errorf("%s", err)
		}
	}
}

func TestReplaceYAMLWithExprRepFn(t *testing.T) {
	const (
		delimStart = "{{"
		delimEnd   = "}}"
	)

	tests := []struct {
		env           interface{}
		replaceMapKey bool
		in            string
		want          string
	}{
		{
			map[string]interface{}{
				"hello": "world",
			},
			false,
			`v: "{{ hello }}"`,
			`v: world`,
		},
		{
			map[string]interface{}{
				"hello": 3,
			},
			false,
			`v: "{{ hello }}"`,
			`v: 3`,
		},
		{
			map[string]interface{}{
				"hello": map[string]interface{}{
					"foo": "bar",
				},
			},
			false,
			`v: "{{ hello }}"`,
			`v: {"foo":"bar"}`,
		},
		{
			map[string]interface{}{
				"hello": map[string]interface{}{
					"foo": "bar",
				},
			},
			false,
			`v:   "{{ hello }}"`,
			`v:   {"foo":"bar"}`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := ReplaceYAML(tt.in, ExprRepFn(delimStart, delimEnd, tt.env), tt.replaceMapKey)
			if err != nil {
				t.Error(err)
			}
			if got != tt.want {
				t.Errorf("got %#v\nwant %#v", got, tt.want)
			}
		})
	}
}
