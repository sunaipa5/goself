package goself

type Options struct {
	Author         string //Github username
	Repo           string //Github repo name
	CurrentVersion string //Your application current version
	PackageName    string //Installation package name
	TagEnd         string //Example: window_amd64.tar.gz
	TagEnd2        string //Option 2 if 'TagEnd' is not found
	AppName        string //Your executable app name
	TmpFolderName  string //Folder where the executable will be temporarily hosted Default: .update-tmp
}

type Release struct {
	Version string `json:"tag_name"`
	Assets  []Assets
}

type Assets struct {
	Name         string `json:"name"`
	Download_Url string `json:"browser_download_url"`
}

type Source struct {
	Name         string
	Download_Url string
}
