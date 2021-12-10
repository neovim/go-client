package plugin

import (
	"fmt"
	"testing"
)

const (
	originalManifest = "call remote#host#RegisterPlugin('P', '0', [\n\\ {'type': 'function', 'name': 'Foo', 'sync': 1, 'opts': {}},\n\\ ])\n"
	updatedManifest  = "call remote#host#RegisterPlugin('P', '0', [\n\\ {'type': 'function', 'name': 'Bar', 'sync': 1, 'opts': {}},\n\\ ])\n"
)

func TestReplaceManifest(t *testing.T) {
	t.Parallel()

	var replaceManifestTests = []struct {
		name     string
		original string
		expected string
	}{
		{
			name:     "Original at beginning of file",
			original: fmt.Sprintf("%sline A\nline B\n", originalManifest),
			expected: fmt.Sprintf("%sline A\nline B\n", updatedManifest),
		},
		{
			name:     "Original in middle of file",
			original: fmt.Sprintf("line A\n%sline B\n", originalManifest),
			expected: fmt.Sprintf("line A\n%sline B\n", updatedManifest),
		},
		{
			name:     "Original at end of file",
			original: fmt.Sprintf("line A\nline B\n%s", originalManifest),
			expected: fmt.Sprintf("line A\nline B\n%s", updatedManifest),
		},
		{
			name:     "Original at end of file, no trailing \\n",
			original: fmt.Sprintf("line A\nline B\n%s", originalManifest[:len(originalManifest)-1]),
			expected: fmt.Sprintf("line A\nline B\n%s", updatedManifest),
		},
		{
			name:     "No manifest",
			original: "line A\nline B\n",
			expected: fmt.Sprintf("line A\nline B\n%s", updatedManifest),
		},
		{
			name:     "Empty file",
			original: "",
			expected: updatedManifest,
		},
		{
			name:     "No manifest, no trailing \\n",
			original: "line A\nline B",
			expected: fmt.Sprintf("line A\nline B\n%s", updatedManifest),
		},
		{
			name:     "Extra \\ ])` in file", // ensure non-greedy match trailing ])
			original: fmt.Sprintf("line A\n%sline B\n\\ ])\nline C\n", originalManifest),
			expected: fmt.Sprintf("line A\n%sline B\n\\ ])\nline C\n", updatedManifest),
		},
	}
	for _, tt := range replaceManifestTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := string(replaceManifest("P", []byte(tt.original), []byte(updatedManifest)))
			if actual != tt.expected {
				t.Fatalf("%s\n got = %q\nwant = %q", tt.name, actual, tt.expected)
			}

			// Replace should be idempotent.
			actual = string(replaceManifest("P", []byte(tt.expected), []byte(updatedManifest)))
			if actual != tt.expected {
				t.Fatalf("%s (no change expected)\n got = %q\nwant = %q", tt.name, actual, tt.expected)
			}
		})
	}
}
