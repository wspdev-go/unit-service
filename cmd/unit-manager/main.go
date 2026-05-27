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

		f, err := os.Create("CDR-trace.out")
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
