package totorow

import (
	"testing"
)

func TestGetKeyAndPath(t *testing.T) {
	for name, c := range map[string]struct {
		url, expectKey, expectPath string
		expectError                error
	}{
		"notHasSlash": {
			url:         "test",
			expectError: imageURLInvalid,
		},
		"notHasKey": {
			url:         "/test",
			expectError: imageURLInvalid,
		},
		"notHasPath": {
			url:         "test/",
			expectError: imageURLInvalid,
		},
		"oneSlash": {
			url:         "/",
			expectError: imageURLInvalid,
		},
		"blank": {
			url:         "",
			expectError: imageURLInvalid,
		},
		"normal": {
			url:        "key/path.png",
			expectKey:  "key",
			expectPath: "path.png",
		},
		"pathWithSlash": {
			url:        "key/dir/path.png",
			expectKey:  "key",
			expectPath: "dir/path.png",
		},
	} {
		c := c
		t.Run(name, func(t *testing.T) {
			key, path, err := getKeyAndPath(c.url)
			if err != c.expectError {
				t.Errorf("check error failed: expect[%v], got[%v]", c.expectError, err)
			}
			if key != c.expectKey {
				t.Errorf("check key failed: expect[%v], got[%v]", c.expectKey, key)
			}
			if path != c.expectPath {
				t.Errorf("check path failed: expect[%v], got[%v]", c.expectPath, path)
			}
		})
	}
}
