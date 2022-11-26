package stages

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/backwardspy/juicy-fortress/internal/platform"
)

func DownloadDwarfFortress(dfdir string, verbose bool) error {
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
	err = extract(archive, dfdir, platform.DwarfFortressDecapitate, verbose)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	return nil
}

func InstallDFHack(dfdir string, verbose bool) error {
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
	err = extract(archive, dfdir, false, verbose)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	return nil
}

func InstallTWBT(dfdir string, verbose bool) error {
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
	err = extract(archive, zipDir, false, verbose)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	// cp the files we want into the right places
	cp := newCopyWrapper()
	cp.Copy(path.Join(zipDir, platform.TWBTLibSrc), path.Join(dfdir, platform.TWBTLibDst), verbose)
	cp.Copy(path.Join(zipDir, "overrides.txt"), path.Join(dfdir, "data/init/overrides.txt"), verbose)

	for _, artFile := range []string{"shadows.png", "transparent1px.png", "white1px.png"} {
		cp.Copy(path.Join(zipDir, artFile), path.Join(dfdir, "data/art", artFile), verbose)
	}

	if cp.err != nil {
		return cp.err
	}

	return nil
}

func InstallSpacefox(dfdir string, verbose bool) error {
	sentinelPath := path.Join(dfdir, platform.SpacefoxSentinel)
	if info, _ := os.Stat(sentinelPath); info != nil {
		fmt.Printf("path \"%v\" already exists. delete it to reinstall Spacefox.\n", sentinelPath)
		return nil
	}

	url, err := url.Parse(platform.SpacefoxURL)
	if err != nil {
		return fmt.Errorf("failed to parse download url (%v): %v", platform.SpacefoxURL, err)
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

	fmt.Printf("extracting %v into %v...\n", archive, cacheDir)
	err = extract(archive, cacheDir, false, verbose)
	if err != nil {
		return fmt.Errorf("failed to extract archive %v: %v", archive, err)
	}

	prefix := path.Join(cacheDir, platform.SpacefoxDir)
	cp := newCopyWrapper()
	cp.Copy(path.Join(prefix, "data/twbt_art"), path.Join(prefix, "data/art"), verbose)
	cp.Copy(path.Join(prefix, "data/twbt_init"), path.Join(prefix, "data/init"), verbose)
	cp.Copy(path.Join(prefix, "raw/twbt_graphics"), path.Join(prefix, "raw/graphics"), verbose)
	cp.Copy(path.Join(prefix, "raw/twbt_objects"), path.Join(prefix, "raw/objects"), verbose)
	cp.Copy(path.Join(prefix, "data"), path.Join(dfdir, "data"), verbose)
	cp.Copy(path.Join(prefix, "raw"), path.Join(dfdir, "raw"), verbose)

	return nil
}

func hasMultilevel(confPath string) (bool, error) {
	conf, err := os.Open(confPath)
	if err != nil {
		return false, fmt.Errorf("failed to open %v: %v", confPath, err)
	}
	defer conf.Close()

	scanner := bufio.NewScanner(conf)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "multilevel") {
			return true, nil
		}
	}

	return false, nil
}

func enableMultilevel(dfdir string, verbose bool) error {
	confPath := path.Join(dfdir, "dfhack-config/default/init/dfhack.init")
	exists, err := hasMultilevel(confPath)
	if err != nil {
		return fmt.Errorf("failed to check for multilevel feature: %v", err)
	}
	if exists {
		if verbose {
			fmt.Printf("multilevel already enabled in %v\n", confPath)
		}
		return nil
	}

	fmt.Println("enabling multilevel rendering")

	conf, err := os.OpenFile(confPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open dfhack conf: %v", err)
	}
	defer conf.Close()

	if _, err := conf.WriteString("multilevel 5"); err != nil {
		return fmt.Errorf("failed to enable multilevel rendering: %v", err)
	}

    if verbose {
        fmt.Println("multilevel rendering enabled!")
    }

	return nil
}

func ApplyPatches(dfdir string, verbose bool) error {
	return enableMultilevel(dfdir, verbose)
}
