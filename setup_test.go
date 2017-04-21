package totorow

import (
	"reflect"
	"testing"

	"github.com/mholt/caddy"
)

func TestParse(t *testing.T) {
	for name, c := range map[string]struct {
		input       string
		shouldError bool
		expect      Totorow
	}{
		"noArg": {
			input:       "totorow",
			shouldError: true,
		},
		"oneArg": {
			input: "totorow repo.xml",
			expect: Totorow{
				RepoConfig: "repo.xml",
				BaseURL:    "/",
			},
		},
		"twoArg": {
			input: "totorow repo.xml /test",
			expect: Totorow{
				RepoConfig: "repo.xml",
				BaseURL:    "/test",
			},
		},
		"tooManyArg": {
			input:       "totorow repo.xml /test arg",
			shouldError: true,
		},
	} {
		c := c
		t.Run(name, func(t *testing.T) {
			var got Totorow
			err := parse(caddy.NewTestController("http", c.input), &got)
			if c.shouldError && err == nil {
				t.Error("should error, but not")
			}
			if !c.shouldError && err != nil {
				t.Errorf("got unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, c.expect) {
				t.Errorf("don't match expected: expected[%#v], got[%#v]", c.expect, got)
			}
		})
	}
}
