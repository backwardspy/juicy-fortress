package stages

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	cp "github.com/otiai10/copy"
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
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to make cache dir (%v): %v", cacheDir, err)
	}

	return cacheDir, nil
}

type copyWrapper struct {
	err error
}

func newCopyWrapper() copyWrapper {
	return copyWrapper{
		err: nil,
	}
}

func (cw *copyWrapper) Copy(from string, to string, verbose bool) {
	if cw.err != nil {
		return
	}

	if verbose {
        fmt.Printf("copy: %v to %v\n", from, to)
	}

	cp.Copy(from, to)
}
