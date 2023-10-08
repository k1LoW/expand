package expand

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInterpolateRepFn(t *testing.T) {
	tests := []struct {
		envs map[string]string
		in   string
		want string
	}{
		{map[string]string{}, "hello ${UNDEFINED}", "hello "},
		{map[string]string{}, "hello ${UNDEFINED:-world}", "hello world"},
		{map[string]string{"UNDEFINED": "space"}, "hello ${UNDEFINED:-world}", "hello space"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			for k, v := range tt.envs {
				t.Setenv(k, v)
			}
			repFn := InterpolateRepFn(os.LookupEnv)
			got, err := repFn(tt.in)
			if err != nil {
				t.Error(err)
			}
			if got != tt.want {
				t.Errorf("got %v\nwant %v", got, tt.want)
			}
		})
	}
}

func TestSubstrWithDelims(t *testing.T) {
	tests := []struct {
		delimStart string
		delimEnd   string
		in         string
		want       [][]string
	}{
		{"{{", "}}", " {{ hello }} {{ value }} ", [][]string{{"{{ hello }}", " hello "}, {"{{ value }}", " value "}}},
		{"{{", "}}", "{{ {{ hello }} {{ value }}", [][]string{{"{{ {{ hello }}", " {{ hello "}, {"{{ value }}", " value "}}},
		{`"{{`, `}}"`, ` "{{ hello }}" "{{ value }}" `, [][]string{{`"{{ hello }}"`, " hello "}, {`"{{ value }}"`, " value "}}},
		{`"{{`, `}}"`, `"{{ {{ hello }}" {{ value }}`, [][]string{{`"{{ {{ hello }}"`, " {{ hello "}}},
		{`"{{`, `}}"`, `"{{ hello }}-{{ value }}"`, [][]string{{`"{{ hello }}-{{ value }}"`, " hello }}-{{ value "}}},
		{"{{", "}}", "{{ hello", [][]string{}},
		{"{{", "}}", "hello }}", [][]string{}},
		{"%%", "%%", " {{ hello }} {{ value }} ", [][]string{}},
		{"%%", "%%", " %% hello %% %% value %% ", [][]string{{"%% hello %%", " hello "}, {"%% value %%", " value "}}},
		{"%%", "%%", "%% %% hello %% %% value %%", [][]string{{"%% %%", " "}, {"%% %%", " "}}},
	}
	for _, tt := range tests {
		got := substrWithDelims(tt.delimStart, tt.delimEnd, tt.in)
		if diff := cmp.Diff(got, tt.want, nil); diff != "" {
			t.Errorf("%s", diff)
		}
	}
}

func TestExprRepFn(t *testing.T) {
	tests := []struct {
		delimStart string
		delimEnd   string
		env        interface{}
		in         string
		want       string
	}{
		{
			"{{",
			"}}",
			map[string]interface{}{
				"hello": "world",
			},
			" {{ hello }} {{ value }} ",
			" world  ",
		},
		{
			"{{",
			"}}",
			map[string]interface{}{
				"hello": "world",
			},
			` "{{ hello }}" "{{ value }}" `,
			" world  ",
		},
		{
			"{{",
			"}}",
			map[string]interface{}{
				"hello": "world",
				"value": "one",
			},
			` "{{ hello }}-{{ value }}" `,
			` world-one `,
		},
		{
			"{{",
			"}}",
			map[string]interface{}{
				"hello": "world\nworld",
			},
			`"{{ hello }}"`,
			`world
world`,
		},
		{
			"{{",
			"}}",
			map[string]interface{}{
				"hello": -3,
			},
			`"{{ hello }}"`,
			`-3`,
		},
	}
	for _, tt := range tests {
		repFn := ExprRepFn(tt.delimStart, tt.delimEnd, tt.env)
		got, err := repFn(tt.in)
		if err != nil {
			t.Error(err)
		}
		if got != tt.want {
			t.Errorf("got %#v\nwant %#v", got, tt.want)
		}
	}
}
