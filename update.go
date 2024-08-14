package goself

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// isUpdateAvailable - true
func (options Options) CheckUpdate() (bool, Release) {
	resp, err := http.Get("https://api.github.com/repos/" + options.Author + "/" + options.Repo + "/releases/latest")
	if err != nil {
		fmt.Println(err)
		return false, Release{}
	}
	defer resp.Body.Close()

	jsonDoc, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false, Release{}
	}

	var release Release
	if err := json.Unmarshal(jsonDoc, &release); err != nil {
		fmt.Println(err)
		return false, Release{}
	}

	if release.Version != options.CurrentVersion {
		return true, release
	} else {
		return false, Release{}
	}

}

func (options Options) findSource(release Release, tagEnd string) (int, Source) {
	var source Source
	count := 0
	for _, asset := range release.Assets {
		if count > 0 {
			source.Name += " "
			source.Download_Url += " "
		}
		if strings.HasSuffix(asset.Name, tagEnd) {
			count++
			source.Name += asset.Name
			source.Download_Url += asset.Download_Url
		}
	}

	return count, source
}

func (options Options) DownloadUpdate(release Release) error {
	if options.TagEnd == "" {
		return fmt.Errorf("TagEnd require!")
	}

	count, source := options.findSource(release, options.TagEnd)

	if count > 1 {
		return fmt.Errorf("multiple source found! please change 'TagEnd': %v", source)
	}

	if options.TagEnd2 == "" {
		return fmt.Errorf("not found any source! TagEnd2 not found, can't search for another source!")
	}

	count, source = options.findSource(release, options.TagEnd2)

	if count > 1 {
		return fmt.Errorf("multiple source found! please change 'TagEnd': %v", source)
	} else if count < 1 {
		return fmt.Errorf("not found any source! plese change 'TagEnd'")
	}

	file, err := Download_Update_File(source.Download_Url)
	if err != nil {
		return fmt.Errorf("file cannot download: %v", err)
	}

	if strings.HasSuffix(source.Name, ".zip") {
		options.ZipExtractor(file)
	} else if strings.HasSuffix(source.Name, ".tar.gz") {
		options.Targz_extractor(file)
	}

	return nil
}

func (options Options) ApplyUpdate() error {
	if err := options.StartUpdate(); err != nil {
		return err
	}

	if err := options.EndUpdate(); err != nil {
		return err
	}
	return nil
}

// Change old executable with new executable
func (options Options) StartUpdate() error {

	tmpFolderName := options.TmpFolderName
	if tmpFolderName == "" {
		options.TmpFolderName = ".update-tmp"
	}

	var err error
	newBinary := options.TmpFolderName + "/" + options.AppName

	err = os.Chmod(newBinary, 0755)
	if err != nil {
		return fmt.Errorf("exec permission cannot get: %v", err)
	}

	cmd := exec.Command(newBinary, "updated")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("new app cannot start: %v", err)
	}

	return nil
}

// Close old executable file
func (options Options) EndUpdate() error {
	tmpFolderName := options.TmpFolderName
	if tmpFolderName == "" {
		options.TmpFolderName = ".update-tmp"
	}

	newBinary := options.TmpFolderName + "/" + options.AppName
	oldBinary := "./" + options.AppName

	bakPath := oldBinary + ".bak"
	err := os.Rename(oldBinary, bakPath)
	if err != nil {
		return fmt.Errorf("old file cannot backup: %v", err)
	}

	err = os.Rename(newBinary, oldBinary)
	if err != nil {
		os.Rename(bakPath, oldBinary)
		return fmt.Errorf("file cannot changed: %v", err)
	}

	os.Remove(bakPath)
	os.Remove(options.TmpFolderName)
	os.Exit(0)

	return nil
}
