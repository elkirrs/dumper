package utils_test

import (
	"dumper/pkg/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchivedLocalFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	targetDir := filepath.Join(tmpDir, "target")

	assert.NoError(t, os.MkdirAll(sourceDir, 0755))

	file1 := filepath.Join(sourceDir, "db1.sql")
	file2 := filepath.Join(sourceDir, "db1_extra.sql")
	file3 := filepath.Join(sourceDir, "db2.sql")
	for _, f := range []string{file1, file2, file3} {
		_, err := os.Create(f)
		assert.NoError(t, err)
	}

	err := utils.ArchivedLocalFile("db1", file1, sourceDir, targetDir)
	assert.NoError(t, err)

	assert.FileExists(t, file1)
	assert.FileExists(t, filepath.Join(targetDir, "db1_extra.sql"))
	assert.FileExists(t, file3)
}

func TestArchivedLocalFile_CreateDirError(t *testing.T) {
	sourceDir := t.TempDir()
	file := filepath.Join(sourceDir, "file.sql")
	_, err := os.Create(file)
	assert.NoError(t, err)

	targetDir := "/root/forbidden"

	err = utils.ArchivedLocalFile("db1", file, sourceDir, targetDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't create a directory")
}

func TestArchivedLocalFile_RenameError(t *testing.T) {
	tmpDir := t.TempDir()
	sourceDir := filepath.Join(tmpDir, "source")
	targetDir := filepath.Join(tmpDir, "target")
	assert.NoError(t, os.MkdirAll(sourceDir, 0755))
	assert.NoError(t, os.MkdirAll(targetDir, 0755))

	file := filepath.Join(sourceDir, "db1.sql")
	_, err := os.Create(file)
	assert.NoError(t, err)

	conflict := filepath.Join(targetDir, "db1_extra.sql")
	assert.NoError(t, os.MkdirAll(conflict, 0755))

	fileToMove := filepath.Join(sourceDir, "db1_extra.sql")
	_, err = os.Create(fileToMove)
	assert.NoError(t, err)

	err = utils.ArchivedLocalFile("db1", file, sourceDir, targetDir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't move the file")
}
