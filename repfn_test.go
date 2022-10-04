package expand

import (
	"fmt"
	"os"
	"testing"
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
