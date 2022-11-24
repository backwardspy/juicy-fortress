package stages

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mholt/archiver/v4"
)

func extract(source string, destination string, decapitatePaths bool) error {
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open %v: %v", source, err)
	}

	format, input, err := archiver.Identify(source, file)
	if err != nil {
		return fmt.Errorf("unsupported archive type for file %v: %v", source, err)
	}

	if ex, ok := format.(archiver.Extractor); ok {
		return ex.Extract(context.Background(), input, nil, func(ctx context.Context, f archiver.File) error {
			if f.IsDir() {
				return nil
			}

			var outputPath string
			if decapitatePaths {
				outputPath = path.Join(destination, decapitate(f.NameInArchive))
			} else {
				outputPath = path.Join(destination, f.NameInArchive)
			}

			// create symlinks
			if f.LinkTarget != "" {
				fmt.Printf("%v is a link to %v\n", f.NameInArchive, f.LinkTarget)
				os.Symlink(f.LinkTarget, outputPath)
				return nil
			}

			reader, err := f.Open()
			if err != nil {
				return err
			}
			defer reader.Close()

			fmt.Printf("extract: %v -> %v\n", f.NameInArchive, outputPath)

			writer, err := safeCreateFile(outputPath, f.Mode())
			if err != nil {
				return fmt.Errorf("failed to create %v: %v", outputPath, err)
			}
			defer writer.Close()

			if _, err := io.Copy(writer, reader); err != nil {
				return fmt.Errorf("failed to write %v: %v", outputPath, err)
			}

			return nil
		})
	} else {
		return fmt.Errorf("expected an extractable file")
	}
}
