# Goself

[![Go Reference](https://pkg.go.dev/badge/github.com/sunaipa5/reqtor.svg)](https://pkg.go.dev/github.com/sunaipa5/goself)
self-update library for go

## How to work

It checks the version in the public github repo and updates the application, it only works on single executable (binaries),
you can also manually download and update the file if you want, `tar.gz` and `.zip` archives are supported, file extraction is automatic, tested on `cli` and `wails` apps.

## Example

```go
func CheckUpdate() string {
	updaterOptions := goself.Options{
		Author:         "yourGithubName",
		Repo:           "yourGithubRepoName",
		CurrentVersion: "0.0.1",
		TagEnd:         "linux_amd64.tar.gz",
                TagEnd2:        "linux_amd64.zip", //not mandatory, option 2 when tagEnd is not found
		AppName:        "yourExecutableName",
	}

	isUpdateAvailable, release := updaterOptions.CheckUpdate()

	if !isUpdateAvailable {
		return "App up to date"
	}

	go func() {
		if err := updaterOptions.DownloadUpdate(release); err != nil {
			fmt.Println(err)
			return
		}

		if err := updaterOptions.ApplyUpdate(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	return "New version installing..."
}
```

## Wails Example

```go
func (a *App) CheckUpdate() string {
	updaterOptions := goself.Options{
		Author:         "yourGithubName",
		Repo:           "yourGithubRepoName",
		CurrentVersion: "0.0.1",
		TagEnd:         "linux_amd64.tar.gz",
                TagEnd2:        "linux_amd64.zip", //not mandatory, option 2 when tagEnd is not found
		AppName:        "yourExecutableName",
	}

	isUpdateAvailable, release := updaterOptions.CheckUpdate()

	if !isUpdateAvailable {
		return "App up to date"
	}

	go func() {
		if err := updaterOptions.DownloadUpdate(release); err != nil {
			fmt.Println(err)
			return
		}

		if err := updaterOptions.ApplyUpdate(); err != nil {
			fmt.Println(err)
			return
		}
	}()

	return "New version installing..."
}
```
