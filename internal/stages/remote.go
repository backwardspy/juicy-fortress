package stages

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func download(url *url.URL, outputPath string) error {
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create %v: %v", outputPath, err)
	}
	defer out.Close()

	resp, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("failed to GET %v: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status from %v: %v", url, resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save downloaded file to %v: %v", outputPath, err)
	}

	return nil
}
