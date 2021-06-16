package main

import (
	"errors"
	"log"
	"runtime"
	"test_task/internal/config"
	"test_task/internal/server"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() (err error) {
	defer func() {
		if e := recover(); e != nil {
			switch ee := e.(type) {
			case error:
				err = ee
			case string:
				err = errors.New(ee)
			default:
				err = errors.New("undefined error")
			}
		}
	}()
	conf, e := config.New("./config/config.json")
	if e != nil {
		return e
	}
	sr := server.New(conf)
	if err := sr.Run(); err != nil {
		return err
	}

	return nil
}
