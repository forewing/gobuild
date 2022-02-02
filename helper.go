package gobuild

import (
	"path/filepath"
)

func evalSymlinksAbs(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(absPath)
}
