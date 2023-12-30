package expand

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/buildkite/interpolate"
)

type repFn func(in string) (string, error)

var envPlaceholderRe = regexp.MustCompile(`\${.+}`)

func InterpolateRepFn(mapping func(string) (string, bool)) repFn {
	const repToken = "__NOT_INTERPOLATE_START__"
	mapper := Mapper{mapping: mapping}
	return func(in string) (string, error) {
		if !envPlaceholderRe.MatchString(in) {
			return in, nil
		}
		r := strings.NewReplacer("${", "${", "$", repToken)
		rr := strings.NewReplacer(repToken, "$")

		replace, err := interpolate.Interpolate(mapper, r.Replace(in))
		if err != nil {
			return in, err
		}
		return rr.Replace(replace), nil
	}
}

func ExprRepFn(delimStart, delimEnd string, env any) repFn {
	const strDQuote = `"`
	return func(in string) (string, error) {
		if !strings.Contains(in, delimStart) {
			return in, nil
		}

		if strings.Count(in, strDQuote) >= 2 {
			oldnew := []string{}
			dds := fmt.Sprintf("%s%s", strDQuote, delimStart)
			dde := fmt.Sprintf("%s%s", delimEnd, strDQuote)
			matches := substrWithDelims(dds, dde, in)
			for _, m := range matches {
				oldnew = append(oldnew, m[0], fmt.Sprintf("%s%s%s", delimStart, m[1], delimEnd))
			}
			rep := strings.NewReplacer(oldnew...)
			in = rep.Replace(in)
		}

		matches := substrWithDelims(delimStart, delimEnd, in)
		oldnew := []string{}
		for _, m := range matches {
			o, err := expr.Eval(m[1], env)
			if err != nil {
				return in, err
			}
			var s string
			switch v := o.(type) {
			case string:
				// Stringify only one expression.
				stat := getNumberStat(v)
				if strings.TrimSpace(in) == m[0] && stat.isNum {
					s = fmt.Sprintf("%q", v)
				} else if strings.TrimSpace(in) == m[0] && (v == "true" || v == "false") {
					s = fmt.Sprintf("%q", v)
				} else {
					s = v
				}
			case int64:
				s = strconv.Itoa(int(v))
			case uint64:
				s = strconv.Itoa(int(v))
			case float64:
				s = strconv.FormatFloat(v, 'f', -1, 64)
			case int:
				s = strconv.Itoa(v)
			case bool:
				s = strconv.FormatBool(v)
			case map[string]any, []any:
				bytes, err := json.Marshal(v)
				if err != nil {
					return in, err
				} else {
					s = string(bytes)
				}
			default:
				s = ""
			}
			oldnew = append(oldnew, m[0], s)
		}
		rep := strings.NewReplacer(oldnew...)
		return rep.Replace(in), nil
	}
}

func substrWithDelims(delimStart, delimEnd, in string) [][]string {
	matches := [][]string{}
	i := 0
	for {
		in = in[i:]
		m, c := trySubstr(delimStart, delimEnd, in)
		if c < 0 {
			break
		}
		matches = append(matches, m)
		i = c
	}
	return matches
}

func trySubstr(delimStart, delimEnd, in string) ([]string, int) {
	if delimStart == delimEnd {
		if strings.Count(in, delimStart) < 2 {
			return nil, -1
		}
	}
	si := strings.Index(in, delimStart)
	if si < 0 {
		return nil, -1
	}
	e := strings.Index(in[si+len(delimStart):], delimEnd)
	if e < 0 {
		return nil, -1
	}
	se := e + si + len(delimStart)
	if si >= se {
		return nil, -1
	}
	wd := in[si : se+len(delimEnd)]
	id := strings.TrimSuffix(strings.TrimPrefix(wd, delimStart), delimEnd)
	return []string{wd, id}, se + len(delimEnd)
}
