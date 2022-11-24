package stages

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"
)

func decapitate(filePath string) string {
	parts := strings.Split(filePath, string(os.PathSeparator))
	return path.Join(parts[1:]...)
}

func safeCreateFile(filePath string, mode fs.FileMode) (*os.File, error) {
	if err := os.MkdirAll(path.Dir(filePath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %v: %v", path.Dir(filePath), err)
	}

	return os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
}

func cacheDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get cache dir: %v", err)
	}

	cacheDir = path.Join(cacheDir, "juicy-fortress")
	if err := os.Mkdir(cacheDir, 0755); !errors.Is(err, os.ErrExist) {
		return "", fmt.Errorf("failed to make cache dir (%v): %v", cacheDir, err)
	}

	return cacheDir, nil
}

type copyWrapper struct {
	err error
}

func (cw *copyWrapper) Copy(from string, to string) {
	if cw.err != nil {
		return
	}

	fmt.Printf("copying %v to %v\n", from, to)

	info, err := os.Stat(from)
	if err != nil {
		cw.err = err
		return
	}

	if !info.Mode().IsRegular() {
		cw.err = fmt.Errorf("%s is not a regular file", from)
		return
	}

	fromFile, err := os.Open(from)
	if err != nil {
		cw.err = err
		return
	}
	defer fromFile.Close()

	fromState, err := fromFile.Stat()
	if err != nil {
		cw.err = err
		return
	}

	toFile, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fromState.Mode())
	if err != nil {
		cw.err = err
		return
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, fromFile)
	cw.err = err
}
