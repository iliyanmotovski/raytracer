package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/iliyanmotovski/raytracer/backend"
	"github.com/iliyanmotovski/raytracer/backend/persistent"
	"github.com/iliyanmotovski/raytracer/backend/server/http/api"
)

var (
	httpPort   = flag.String("port", "8008", "http listen address")
	configPath = flag.String("config", "config.txt", "path to config file")
)

func main() {
	flag.Parse()

	sceneRepo := persistent.NewInMemorySceneRepository()
	configRepo := persistent.NewInMemoryConfigRepository()

	ctx := context.Background()
	configurator := backend.NewTextFileConfigurator(*configPath)

	c, err := configurator.Parse(ctx, configRepo)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	cc := make(chan *backend.ConfigChan)
	initialSrrc := make(chan *backend.SceneReloadResponse)
	hotReloadSrrc := make(chan *backend.SceneReloadResponse)
	createConfigHandlerSrrc := make(chan *backend.SceneReloadResponse)

	srrcFactory := backend.SceneReloadResponseChanFactory{
		backend.Initial:             initialSrrc,
		backend.HotReload:           hotReloadSrrc,
		backend.CreateConfigHandler: createConfigHandlerSrrc,
	}

	sceneReloadDaemon := backend.NewSceneReloadDaemon(sceneRepo, cc, srrcFactory)
	sceneReloadDaemon.Start(1)

	cc <- &backend.ConfigChan{Ctx: ctx, Config: c, ResponseChan: backend.Initial}
	resp := <-srrcFactory[backend.Initial]
	if resp.Err != nil {
		log.Println(resp.Err)
	}

	apiRoot := mux.NewRouter().PathPrefix("/api/v1").Subrouter()
	apiRoot.Handle("/scene", api.GetScene(sceneRepo)).Methods("GET")
	apiRoot.Handle("/scene/config", api.CreateConfiguration(cc, srrcFactory)).Methods("POST")

	fileServer := http.FileServer(http.Dir("../frontend"))

	http.HandleFunc("/{path:.+\\.[a-z0-9]+$}", handler)
	http.Handle("/", http.StripPrefix("/", fileServer))
	http.Handle("/api/v1/", apiRoot)

	errChan := make(chan error)
	go func() {
		l, err := net.Listen("tcp", ":"+*httpPort)
		if err != nil {
			errChan <- err
		}

		if err := http.Serve(l, nil); err != nil {
			errChan <- err
		}
	}()

	// configuration hot reload
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGTERM)
	for {
		select {
		case s := <-sigchan:
			switch s {
			case syscall.SIGHUP:
				log.Println("Reloading configuration")
				c, err = configurator.Parse(ctx, configRepo)
				if err != nil {
					log.Println(err)
					continue
				}

				cc <- &backend.ConfigChan{Ctx: ctx, Config: c, ResponseChan: backend.HotReload}
				resp := <-srrcFactory[backend.HotReload]
				if resp.Err != nil {
					log.Println(resp.Err)
				}

				continue
			case syscall.SIGTERM:
				os.Exit(0)
			}
		case httpErr := <-errChan:
			log.Println(httpErr)
			os.Exit(1)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("../frontend/index.html")
	t.Execute(w, "")
}
