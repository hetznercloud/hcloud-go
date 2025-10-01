package envutil

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:unparam
func writeTmpFile(t *testing.T, tmpDir, filename, content string) string {
	filepath := path.Join(tmpDir, filename)

	err := os.WriteFile(filepath, []byte(content), 0644)
	require.NoError(t, err)

	return filepath
}

func TestLookupEnvWithFile(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(t *testing.T, tmpDir string)
		want  func(t *testing.T, value string, err error)
	}{
		{
			name:  "without any environment",
			setup: func(_ *testing.T, _ string) {},
			want: func(t *testing.T, value string, err error) {
				assert.NoError(t, err)
				assert.Empty(t, value)
			},
		},
		{
			name: "value from environment",
			setup: func(t *testing.T, tmpDir string) {
				t.Setenv("CONFIG", "value")

				// Test for precedence
				filepath := writeTmpFile(t, tmpDir, "config", "content")
				t.Setenv("CONFIG_FILE", filepath)
			},
			want: func(t *testing.T, value string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "value", value)
			},
		},
		{
			name: "empty value from environment",
			setup: func(t *testing.T, tmpDir string) {
				t.Setenv("CONFIG", "")

				// Test for precedence
				filepath := writeTmpFile(t, tmpDir, "config", "content")
				t.Setenv("CONFIG_FILE", filepath)
			},
			want: func(t *testing.T, value string, err error) {
				assert.NoError(t, err)
				assert.Empty(t, value)
			},
		},
		{
			name: "value from file",
			setup: func(t *testing.T, tmpDir string) {
				// The extra spaces ensure that the value is sanitized
				filepath := writeTmpFile(t, tmpDir, "config", "content  ")
				t.Setenv("CONFIG_FILE", filepath)
			},
			want: func(t *testing.T, value string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "content", value)
			},
		},
		{
			name: "empty value from file",
			setup: func(t *testing.T, tmpDir string) {
				filepath := writeTmpFile(t, tmpDir, "config", "")
				t.Setenv("CONFIG_FILE", filepath)
			},
			want: func(t *testing.T, value string, err error) {
				assert.NoError(t, err)
				assert.Empty(t, value)
			},
		},
		{
			name: "missing file",
			setup: func(t *testing.T, _ string) {
				t.Setenv("CONFIG_FILE", "/tmp/this-file-does-not-exits")
			},
			want: func(t *testing.T, value string, err error) {
				assert.Error(t, err, "failed to read CONFIG_FILE: open /tmp/this-file-does-not-exits: no such file or directory")
				assert.Empty(t, value)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testCase.setup(t, tmpDir)
			value, err := LookupEnvWithFile("CONFIG")
			testCase.want(t, value, err)
		})
	}
}
