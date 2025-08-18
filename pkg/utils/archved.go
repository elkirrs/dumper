package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ArchivedLocalFile(dbName, file, sourceDir, targetDir string) error {

	actualFile := fmt.Sprintf("%s/%s", sourceDir, filepath.Base(file))

	if err := createDir(targetDir); err != nil {
		return fmt.Errorf("couldn't create a directory %s: %v", targetDir, err)
	}

	absActualFile, err := filepath.Abs(actualFile)
	if err != nil {
		return fmt.Errorf("couldn't get an absolute path for %s: %v", actualFile, err)
	}

	patterns := []string{
		fmt.Sprintf("%s/%s*.sql", sourceDir, dbName),
		fmt.Sprintf("%s/%s*", sourceDir, dbName),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return fmt.Errorf("error when searching for files using a template %s: %v", pattern, err)
		}

		for _, file := range matches {
			absPath, err := filepath.Abs(file)
			if err != nil {
				return fmt.Errorf("couldn't get an absolute path for %s: %v", file, err)
			}

			if absPath == absActualFile {
				continue
			}

			base := filepath.Base(file)
			dest := filepath.Join(targetDir, base)
			err = os.Rename(file, dest)
			if err != nil {
				return fmt.Errorf("couldn't move the file %s -> %s: %v", file, dest, err)
			}
			fmt.Printf("The %s file moved to %s\n", base, targetDir)
		}
	}

	return nil
}

func createDir(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("couldn't create a directory %s: %v", path, err)
	}
	return nil
}
