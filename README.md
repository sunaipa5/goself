# Goself

self-update library for go 

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
