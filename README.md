# raytracer

## A simple 2d ray tracing app in Go and p5.js

### To start:

`cd cmd`  
`go run main.go`

Go to `localhost:8008` in the browser and move the light source.

### Takes 2 optional flags:

`port` - http port, defaults to `8008`  
`config` - path to the config file, defaults to `config.txt`

### Config hot reload:

The app listens for signal `SIGHUP` to reload its configuration and for `SIGTERM` to exit.

![alt text](https://i.ibb.co/LCCDxM4/scene.jpg)
