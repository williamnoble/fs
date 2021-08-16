package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	ErrUnableToRetrieveHostname = errors.New("unable to retrieve hostname")
	ErrUnableToRetrieveLocalIP  = errors.New("unable to retrieve local ip address")
	ErrGeneralisedHostError     = errors.New("error in obtaining local ip address")
)

const (
	Port = ":3000"
)

type App struct {
	IP       string
	Port     string
	Server   FileServer
	RootPath string
}

func NewApp(rootDir string) *App {
	return &App{
		IP:       localIPAddres(),
		Port:     Port,
		RootPath: rootDir,
		Server: FileServer{http.Dir(rootDir),
			IndexTemplate{
				Name:           indexTemplate,
				FileName:       indexTemplateName,
				DirectoryFiles: DirectoryContent{},
			},
		},
	}
}

func (a *App) health() {
	fmt.Printf("Fs Serving %q via %s%s\n", a.RootPath, a.IP, a.Port)
}

func (a *App) listen() {
	err := http.ListenAndServe(a.Port, a.logger(a.Server))
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	rootDir := generatePath()
	app := NewApp(rootDir)
	app.health()
	app.listen()

}
