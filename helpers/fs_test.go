package helpers

import (
"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var (
	filesToBeRemoved = []string{
		"/test/file",
		"/test/file.ext",
		"file.ext",
	}
	files = []string{
		"/test/exists",
		"/test/exists.ext",
		"exists.ext",
	}
	filesNotExist = []string {
		"/text/no-file",
		"no-file",
		"no-file.ext",
	}
)

func init() {
	appFS = afero.NewMemMapFs()
	for _, f := range append(files, filesToBeRemoved...) {
		afero.WriteFile(appFS, f, []byte(f), 0644)
	}
}

func TestEnsureDirectory(t *testing.T) {
	var folders = []string {
		"/test",
		"/test/deeper",
		"folder",
		"folder/deeper/than/previous",
	}
	for _, f := range folders {
		EnsureDirectory(f)

		ok, err := afero.IsDir(appFS, f)
		assert.True(t, ok)
		assert.Nil(t, err)
	}
}

func TestRemoveFile(t *testing.T) {
	for _, f := range filesToBeRemoved {
		info, err := appFS.Stat(f)
		assert.NotNil(t, info)
		assert.Nil(t, err)

		RemoveFile(f)

		info, err = appFS.Stat(f)
		assert.Nil(t, info)
		assert.NotNil(t, err)
	}
}

func TestGetTempFile(t *testing.T) {
	name, err := GetTempFile()

	info, e := appFS.Stat(name)

	assert.Nil(t, err)
	assert.Nil(t, e)
	assert.NotNil(t, info)
	assert.Contains(t, name, "belit")
}

func TestFileExists_Exists(t *testing.T) {
	for _, f := range files {
		err := FileExists(f)

		assert.Nil(t, err)
	}
}

func TestFileExists_NotExists(t *testing.T) {
	for _, f := range filesNotExist {
		err := FileExists(f)

		assert.NotNil(t, err)
	}
}

func TestGetFileContents_Exists(t *testing.T) {
	for _, f := range files {
		buffer, err := GetFileContents(f)

		assert.Nil(t, err)
		assert.Equal(t, []byte(f), buffer)
	}
}

func TestGetFileContents_NotExists(t *testing.T) {
	for _, f := range filesNotExist {
		_, err := GetFileContents(f)

		assert.NotNil(t, err)
	}
}
