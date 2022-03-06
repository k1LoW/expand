package expand

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/token"
)

// ReplaceYAML replaces the values of YAML (string) using replacefunc.
func ReplaceYAML(s string, replacefunc func(s string) string) string {
	tokens := lexer.Tokenize(s)
	if len(tokens) == 0 {
		return ""
	}
	texts := []string{}
	for _, tk := range tokens {
		lines := strings.Split(tk.Origin, "\n")
		expand := false
		quote := false
		if tk.NextType() != token.MappingValueType {
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
func ExpandYAML(s string, mapping func(string) string) string {
	replacefunc := func(in string) string {
		return os.Expand(in, mapping)
	}
	return ReplaceYAML(s, replacefunc)
}

// ExpandYAML replaces ${var} or $var in the values of YAML ([]byte) based on the mapping function.
func ExpandYAMLBytes(b []byte, mapping func(string) string) []byte {
	return []byte(ExpandYAML(string(b), mapping))
}

// ExpandenvYAML replaces ${var} or $var in the values of YAML (string) according to the values
// of the current environment variables.
func ExpandenvYAML(s string) string {
	return ExpandYAML(s, os.Getenv)
}

// ExpandenvYAML replaces ${var} or $var in the values of YAML ([]byte) according to the values
// of the current environment variables.
func ExpandenvYAMLBytes(b []byte) []byte {
	return ExpandYAMLBytes(b, os.Getenv)
}
