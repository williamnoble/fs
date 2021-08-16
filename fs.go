package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"sort"
)

const (
	indexTemplateName = "index"
	indexTemplate     = "index.tmpl"
)

type FileServer struct {
	RootFileSystem http.FileSystem
	Index          IndexTemplate
}

type IndexTemplate struct {
	Name           string
	FileName       string
	DirectoryFiles DirectoryContent
}

func (fs *FileServer) SetDirectoryContent(files DirectoryContent) {
	fs.Index.DirectoryFiles = files
}

func (fs FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestPath := path.Clean(r.URL.Path)
	fs.ServeFile(w, r, fs.RootFileSystem, requestPath)
}

func (fs *FileServer) ServeFile(w http.ResponseWriter, r *http.Request, rootFileSystem http.FileSystem, name string) {
	f, err := rootFileSystem.Open(name)
	if err != nil {
		respondWithError(w, "encountered an error when attempting to serve the file system", http.StatusInternalServerError, err)
	}
	defer func(f http.File) {
		err := f.Close()
		if err != nil {
			log.Println("failed to close file in timely fashion")
		}
	}(f)
	fInfo, err := f.Stat()
	if err != nil {
		respondWithError(w, "encountered an error when attempting to serve the file system", http.StatusInternalServerError, err)
	}

	if fInfo.IsDir() {
		fs.ListDirectory(w, r, f, fs.Index)
		return
	}

	http.ServeContent(w, r, fInfo.Name(), fInfo.ModTime(), f)
}

func (fs *FileServer) RenderTemplate(w http.ResponseWriter) {
	var t *template.Template
	t = template.New("index")
	t, err := template.ParseFiles("index.tmpl")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, fs.Index.DirectoryFiles)
	if err != nil {
		fmt.Println("error when executing template")
	}
}

func (fs *FileServer) ListDirectory(w http.ResponseWriter, r *http.Request, f http.File, index IndexTemplate) {
	// Check rootDirectory is valid
	rootDirectory, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// Get a List of rootDirectory files (reading to end -1)
	rootFiles, err := f.Readdir(-1)
	if err != nil {
		http.Error(w, "error reading file list", http.StatusInternalServerError)
	}

	// Sort the File contents
	sort.Slice(rootFiles, func(i, j int) bool {
		return rootFiles[i].Name() < rootFiles[j].Name()
	})

	w.Header().Set("Content-Type", "text/html ;charset=utf-8")
	fc := fs.returnFileContent(rootFiles)
	d := newDirectoryContent(rootDirectory.Name(), fc, r.Host)
	fs.SetDirectoryContent(*d)
	fmt.Println("Rending Template having completed listing")
	fs.RenderTemplate(w)
}

func (fs *FileServer) returnFileContent(rootFiles []fs.FileInfo) (fc []FileContent) {

	for _, d := range rootFiles {
		name := d.Name()
		fileExtension := "page"
		if d.IsDir() {
			name += "/"
			fileExtension = "folder"
		} else if len(filepath.Ext(name)) > 1 {
			fileExtension = filepath.Ext(name)[1:]
		}

		u := url.URL{Path: name}
		fileContent := newFileContent(name, getHumanReadableSize(d), u, fileExtension)
		fc = append(fc, fileContent)
	}
	return fc
}
