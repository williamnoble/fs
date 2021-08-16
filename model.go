package main

import "net/url"

type DirectoryContent struct {
	DirName string
	Files   []FileContent
	IPAddr  string
}

func newDirectoryContent(dirName string, files []FileContent, ipAddr string) *DirectoryContent {
	return &DirectoryContent{
		DirName: dirName,
		Files:   files,
		IPAddr:  ipAddr,
	}
}

type FileContent struct {
	Name      string
	Size      string
	URL       url.URL
	Extension string
}

func newFileContent(name, size string, url url.URL, extension string) FileContent {
	return FileContent{
		Name:      name,
		Size:      size,
		URL:       url,
		Extension: extension,
	}
}
