package plugin

import (
	"fmt"
	"testing"
)

const (
	originalManifest = "call remote#host#RegisterPlugin('P', '0', [\n\\ {'type': 'function', 'name': 'Foo', 'sync': 1, 'opts': {}},\n\\ ])\n"
	updatedManifest  = "call remote#host#RegisterPlugin('P', '0', [\n\\ {'type': 'function', 'name': 'Bar', 'sync': 1, 'opts': {}},\n\\ ])\n"
)

var replaceManifestTests = []struct {
	what, original, expected string
}{
	{
		"Original at beginning of file",
		fmt.Sprintf("%sline A\nline B\n", originalManifest),
		fmt.Sprintf("%sline A\nline B\n", updatedManifest),
	},
	{
		"Original in middle of file",
		fmt.Sprintf("line A\n%sline B\n", originalManifest),
		fmt.Sprintf("line A\n%sline B\n", updatedManifest),
	},
	{
		"Original at end of file",
		fmt.Sprintf("line A\nline B\n%s", originalManifest),
		fmt.Sprintf("line A\nline B\n%s", updatedManifest),
	},
	{
		"Original at end of file, no trailing \\n",
		fmt.Sprintf("line A\nline B\n%s", originalManifest[:len(originalManifest)-1]),
		fmt.Sprintf("line A\nline B\n%s", updatedManifest),
	},
	{
		"No manifest",
		"line A\nline B\n",
		fmt.Sprintf("line A\nline B\n%s", updatedManifest),
	},
	{
		"Empty file",
		"",
		updatedManifest,
	},
	{
		"No manifest, no trailing \\n",
		"line A\nline B",
		fmt.Sprintf("line A\nline B\n%s", updatedManifest),
	},
	{
		"Extra \\ ])` in file", // ensure non-greedy match trailing ])
		fmt.Sprintf("line A\n%sline B\n\\ ])\nline C\n", originalManifest),
		fmt.Sprintf("line A\n%sline B\n\\ ])\nline C\n", updatedManifest),
	},
}

func TestReplaceManifest(t *testing.T) {
	for _, tt := range replaceManifestTests {
		actual := string(replaceManifest("P", []byte(tt.original), []byte(updatedManifest)))
		if actual != tt.expected {
			t.Errorf("%s\n got = %q\nwant = %q", tt.what, actual, tt.expected)
			continue
		}
		// Replace should be idempotent.
		actual = string(replaceManifest("P", []byte(tt.expected), []byte(updatedManifest)))
		if actual != tt.expected {
			t.Errorf("%s (no change expected)\n got = %q\nwant = %q", tt.what, actual, tt.expected)
		}
	}
}
