package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

var (
	ErrRetrievingHomeDirectory    = errors.New("error retrieving home directory")
	ErrRetrievingWorkingDirectory = errors.New("error retrieving current directory")
	ErrDirectoryDoesNotExist      = errors.New("error, directory does not exist")
	ErrServeFileError             = errors.New("encountered an error when attempting to serve the file system")
)

func getHomeDir() (home string, err error) {
	home, err = os.UserHomeDir()
	if err != nil {
		return "", ErrRetrievingHomeDirectory
	}
	return home, nil
}

// getWorkingDir returns a rooted path name to the current director.
func getWorkingDir() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		return "", ErrRetrievingWorkingDirectory
	}
	return dir, nil
}

func checkDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return ErrDirectoryDoesNotExist
	}
	return nil
}

func parseArgs(dir, home string) (path string) {
	var arg string

	if len(os.Args) == 1 {
		return home
	}

	if len(os.Args) > 1 && len(os.Args) == 2 {
		arg = os.Args[1]
	}

	if len(os.Args) > 2 {
		fmt.Println("You must supply only one rootFileSystem directory for the file system\ne.g. fs ~, fs /home/users/will/")
		os.Exit(1)
	}

	switch arg {
	case "~":
		return home
	case "~/":
		if err := checkDirectoryExists(arg); err != nil {
			switch err {
			case ErrDirectoryDoesNotExist:
				log.Printf("the specified directory %q does not exist", arg)
				return home
			default:
				log.Println("error with checking directory for serving", err)
				return home
			}
		}
	}
	return arg // home = arg
}

// GetHumanReadableSize ...
func getHumanReadableSize(f os.FileInfo) string {
	if f.IsDir() {
		return "--"
	}
	bytes := f.Size()
	mb := float32(bytes) / (1024.0 * 1024.0)
	return fmt.Sprintf("%.2f MB", mb)
}

// generatePath returns a path to the rootFileSystem filesystem.
func generatePath() (rootPath string) {
	dir, err := getWorkingDir()
	if err != nil {
		log.Println(err)
	}

	home, err := getHomeDir()
	if err != nil {
		log.Println(err)
	}
	rootPath = parseArgs(dir, home)
	return rootPath
}

func respondWithError(w http.ResponseWriter, message string, status int, err error) {
	http.Error(w, message, status)
	log.Println("error while serving file system", err.Error())
	return
}

func open(fs http.FileSystem, name string) (file http.File, err error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, ErrServeFileError
	}
	defer func(f http.File) {
		err := f.Close()
		if err != nil {
			log.Println("failed to close file in timely fashion")
		}
	}(f)
	return f, nil
}

func stat(f http.File) (file fs.FileInfo, err error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, ErrServeFileError
	}
	return fileInfo, nil
}
