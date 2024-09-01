package main

import (
	"fmt"
	"log"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
}

type server struct {
	host string
	post string
	url  string
}

func main() {
	fmt.Println("Hello world!")
}
