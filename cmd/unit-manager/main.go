/*
Copyright 2026 WspDev-Go

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0
*/

package main

import (
	"log"
	"net/http"
	"os"
	"runtime/trace"
	"unit-service/internal/app"

	_ "net/http/pprof"
)

const configPath = "config.yml"

const isProfUsage = false

func main() {

	if isProfUsage {
		go func() {
			// Pprof default registered on /debug/pprof/
			_ = http.ListenAndServe("localhost:6060", nil)
		}()

		f, err := os.Create("Unit-trace.out")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		trace.Start(f)
		defer trace.Stop()
	}

	application, err := app.NewApp(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.RunQueuePipeline(); err != nil {
		log.Fatal(err)
	}
}
