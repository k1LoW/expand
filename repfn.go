package expand

import (
	"github.com/buildkite/interpolate"
)

type repFn func(in string) (string, error)

func InterpolateRepFn(mapping func(string) (string, bool)) repFn {
	mapper := Mapper{mapping: mapping}
	return func(in string) (string, error) {
		replace, err := interpolate.Interpolate(mapper, in)
		if err != nil {
			return in, err
		}
		return replace, nil
	}
}
