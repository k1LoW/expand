package expand

import (
	"encoding/json"
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

// ReplaceYAML replaces the tokens of YAML (string) using repFn.
func ReplaceYAML(s string, repFn func(s string) (string, error), opts ...Option) (string, error) {
	c := &config{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return "", err
		}
	}

	var err error
	tokens := lexer.Tokenize(s)
	if len(tokens) == 0 {
		return "", nil
	}
	texts := []string{}
	for _, tk := range tokens {
		lines := strings.Split(tk.Origin, "\n")
		isMapKey := tk.NextType() == token.MappingValueType
		nte := false // Need to expand
		qt := false  // Quote target
		if c.replaceMapKey || !isMapKey {
			switch tk.Type {
			case token.StringType, token.SingleQuoteType, token.DoubleQuoteType:
				nte = true
				if len(lines) == 1 {
					qt = true
				} else if len(lines) == 2 && strings.Trim(lines[1], " ") == "" {
					if tk.Prev != nil && tk.Prev.Type == token.LiteralType && token.Type(tk.Prev.Indicator) == token.Type(token.BlockScalarIndicator) {
						// Block scalars does not quote
						qt = false
					} else {
						qt = true
					}
				}
			}
		}
		if len(lines) == 1 {
			line := lines[0]
			if nte && line != "" {
				line, err = repFn(line)
				if err != nil {
					return "", err
				}
				if isNeedQuoted(qt, isMapKey, c, line) {
					line = quoteLine(line)
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
				if nte && line != "" {
					line, err = repFn(line)
					if err != nil {
						return "", err
					}
					if isNeedQuoted(qt, isMapKey, c, line) {
						line = quoteLine(line)
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

	if strings.HasSuffix(s, "\n") && !strings.HasSuffix(tokens[len(tokens)-1].Value, "\n") {
		return fmt.Sprintf("%s\n", strings.Join(texts, "\n")), nil
	}
	return strings.Join(texts, "\n"), nil
}

// ExpandYAML replaces ${var} or $var in the values of YAML (string) based on the mapping function.
func ExpandYAML(s string, mapping func(string) (string, bool)) string {
	repFn := InterpolateRepFn(mapping)
	rep, _ := ReplaceYAML(s, repFn)
	return rep
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

func quoteOnce(s string) string {
	u, err := strconv.Unquote(s)
	if err != nil {
		return strconv.Quote(s)
	}
	return strconv.Quote(u)
}

func quoteLine(line string) string {
	old := strings.Trim(line, " ")
	new := quoteOnce(old)
	// Avoid duplicate quotes heuristically.
	switch {
	case strings.HasPrefix(new, `"'`) && strings.HasSuffix(new, `'"`):
		// no quote
		return line
	case strings.HasPrefix(new, `"\"`) && strings.HasSuffix(new, `\""`):
		new = fmt.Sprintf(`"%s"`, strings.TrimSuffix(strings.TrimPrefix(new, `"\"`), `\""`))
		return strings.Replace(line, old, new, 1)
	default:
		return strings.Replace(line, old, new, 1)
	}
}

func isNeedQuoted(quoteTarget bool, isMapKey bool, c *config, line string) bool {
	if quoteTarget && token.IsNeedQuoted(line) ||
		// If there is a line break in the result of the conversion of what was one line, quote it.
		strings.Contains(line, "\n") {
		if c.quoteCollection {
			return true
		}
		if isJSONString(line) && !isMapKey {
			// Not quoting to be interpreted as inline YAML
			return false
		}
		return true
	}
	return false
}

func isJSONString(line string) bool {
	if !strings.Contains(line, "{") && !strings.Contains(line, "[") {
		return false
	}
	var v any
	if err := json.Unmarshal([]byte(strings.Trim(line, " ")), &v); err == nil {
		switch v.(type) {
		case []any, map[any]any, map[string]any:
			return true
		}
	}
	return false
}
