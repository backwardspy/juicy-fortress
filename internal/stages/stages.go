package stages

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/backwardspy/juicy-fortress/internal/platform"
)

func DownloadDwarfFortress(dfdir string) error {
	if info, _ := os.Stat(dfdir); info != nil {
		fmt.Printf("folder \"%v\" already exists. delete it to redownload Dwarf Fortress.\n", dfdir)
		return nil
	}

	url, err := url.Parse(platform.DwarfFortressURL)
	if err != nil {
		return fmt.Errorf("failed to parse download url (%v): %v", platform.DwarfFortressURL, err)
	}

	cacheDir, err := cacheDir()
	if err != nil {
		return err
	}

	archive := path.Join(cacheDir, path.Base(url.Path))

	if _, err := os.Stat(archive); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("downloading %v into %v...\n", url, archive)
		err = download(url, archive)
		if err != nil {
			return fmt.Errorf("failed to download %v into %v: %v", url, archive, err)
		}
	} else {
		fmt.Printf("%v already exists\n", archive)
	}

	fmt.Printf("extracting %v into %v...\n", archive, dfdir)
	err = extract(archive, dfdir, true)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	return nil
}

func InstallDFHack(dfdir string) error {
	sentinelPath := path.Join(dfdir, platform.DFHackSentinel)
	if info, _ := os.Stat(sentinelPath); info != nil {
		fmt.Printf("path \"%v\" already exists. delete it to reinstall DFHack.\n", sentinelPath)
		return nil
	}

	url, err := url.Parse(platform.DFHackURL)
	if err != nil {
		return fmt.Errorf("failed to parse download url (%v): %v", platform.DFHackURL, err)
	}

	cacheDir, err := cacheDir()
	if err != nil {
		return err
	}

	archive := path.Join(cacheDir, path.Base(url.Path))

	if _, err := os.Stat(archive); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("downloading %v into %v...\n", url, archive)
		err = download(url, archive)
		if err != nil {
			return fmt.Errorf("failed to download %v into %v: %v", url, archive, err)
		}
	} else {
		fmt.Printf("%v already exists\n", archive)
	}

	fmt.Printf("extracting %v into %v...\n", archive, dfdir)
	err = extract(archive, dfdir, false)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	return nil
}

func InstallTWBT(dfdir string) error {
	sentinelPath := path.Join(dfdir, platform.TWBTLibDst)
	if info, _ := os.Stat(sentinelPath); info != nil {
		fmt.Printf("path \"%v\" already exists. delete it to reinstall TWBT.\n", sentinelPath)
		return nil
	}

	url, err := url.Parse(platform.TWBTUrl)
	if err != nil {
		return fmt.Errorf("failed to parse download url (%v): %v", platform.TWBTUrl, err)
	}

	cacheDir, err := cacheDir()
	if err != nil {
		return err
	}

	archive := path.Join(cacheDir, path.Base(url.Path))

	if _, err := os.Stat(archive); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("downloading %v into %v...\n", url, archive)
		err = download(url, archive)
		if err != nil {
			return fmt.Errorf("failed to download %v into %v: %v", url, archive, err)
		}
	} else {
		fmt.Printf("%v already exists\n", archive)
	}

	// extract into cacheDir because we only want some of the files.
	zipDir := path.Join(cacheDir, "twbt")
	if err := os.MkdirAll(zipDir, 0700); err != nil {
		return fmt.Errorf("failed to make directory %v: %v", zipDir, err)
	}

	fmt.Printf("extracting %v into %v...\n", archive, zipDir)
	err = extract(archive, zipDir, false)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	// cp the files we want into the right places
	cp := copyWrapper{}
	cp.Copy(path.Join(zipDir, platform.TWBTLibSrc), path.Join(dfdir, platform.TWBTLibDst))
	cp.Copy(path.Join(zipDir, "overrides.txt"), path.Join(dfdir, "data/init/overrides.txt"))

	for _, artFile := range []string{"shadows.png", "transparent1px.png", "white1px.png"} {
		cp.Copy(path.Join(zipDir, artFile), path.Join(dfdir, "data/art", artFile))
	}

	if cp.err != nil {
		return cp.err
	}

	return nil
}

func InstallSpacefox(dfdir string) error {
	fmt.Printf("TODO: install spacefox into %v\n", dfdir)
	return nil
}
