package model

import (
	"errors"
	"os"
	"path/filepath"
)

const defaultRootModelDir = "models"

func existsPath(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func resolveModelDir() (string, error) {
	fullPath, err := filepath.Abs(defaultRootModelDir)
	if err != nil {
		return "", err
	}

	exists, err := existsPath(fullPath)
	if err != nil {
		return "", err
	}

	if !exists {
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return "", err
		}
	}

	return fullPath, nil

}

func resolveModelPath(path string) (string, error) {

	return "", nil
}
