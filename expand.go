package expand

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/buildkite/interpolate"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

type Mapper struct {
	mapping func(string) (string, bool)
}

// Implement Env
var _ interpolate.Env = Mapper{}

func (m Mapper) Get(key string) (string, bool) {
	if m.mapping == nil {
		return "", false
	}
	return m.mapping(key)
}

// ReplaceYAML replaces the tokens of YAML (string) using replacefunc.
func ReplaceYAML(s string, replacefunc func(s string) string, replaceMapKey bool) string {
	tokens := lexer.Tokenize(s)
	if len(tokens) == 0 {
		return ""
	}
	texts := []string{}
	for _, tk := range tokens {
		lines := strings.Split(tk.Origin, "\n")
		expand := false
		quote := false
		if replaceMapKey || tk.NextType() != token.MappingValueType {
			switch tk.Type {
			case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
				expand = true
				if len(lines) == 1 {
					quote = true
				} else if len(lines) == 2 && strings.Trim(lines[1], " ") == "" {
					quote = true
				}
			}
		}
		if len(lines) == 1 {
			line := lines[0]
			if expand && line != "" {
				line = replacefunc(line)
				if quote && token.IsNeedQuoted(line) {
					old := strings.Trim(line, " ")
					new := strconv.Quote(old)
					line = strings.Replace(line, old, new, 1)
				}
			}
			if len(texts) == 0 {
				texts = append(texts, line)
			} else {
				text := texts[len(texts)-1]
				texts[len(texts)-1] = text + line
			}
		} else {
			for idx, src := range lines {
				line := src
				if expand && line != "" {
					line = replacefunc(line)
					if quote && token.IsNeedQuoted(line) {
						old := strings.Trim(line, " ")
						new := strconv.Quote(old)
						line = strings.Replace(line, old, new, 1)
					}
				}
				if idx == 0 {
					if len(texts) == 0 {
						texts = append(texts, line)
					} else {
						text := texts[len(texts)-1]
						texts[len(texts)-1] = text + line
					}
				} else {
					texts = append(texts, line)
				}
			}
		}
	}
	return fmt.Sprintf("%s\n", strings.Join(texts, "\n"))
}

// ExpandYAML replaces ${var} or $var in the values of YAML (string) based on the mapping function.
func ExpandYAML(s string, mapping func(string) (string, bool)) string {
	mapper := Mapper{mapping: mapping}
	replacefunc := func(in string) string {
		replace, err := interpolate.Interpolate(mapper, in)
		if err != nil {
			return in
		}
		return replace
	}
	return ReplaceYAML(s, replacefunc, false)
}

// ExpandYAML replaces ${var} or $var in the values of YAML ([]byte) based on the mapping function.
func ExpandYAMLBytes(b []byte, mapping func(string) (string, bool)) []byte {
	return []byte(ExpandYAML(string(b), mapping))
}

// ExpandenvYAML replaces ${var} or $var in the values of YAML (string) according to the values
// of the current environment variables.
func ExpandenvYAML(s string) string {
	return ExpandYAML(s, os.LookupEnv)
}

// ExpandenvYAML replaces ${var} or $var in the values of YAML ([]byte) according to the values
// of the current environment variables.
func ExpandenvYAMLBytes(b []byte) []byte {
	return ExpandYAMLBytes(b, os.LookupEnv)
}
