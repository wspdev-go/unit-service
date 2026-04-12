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
			http.ListenAndServe("localhost:6060", nil)
		}()

		f, err := os.Create("CDR-trace.out")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		trace.Start(f)
		defer trace.Stop()
	}
	app.InitApp(configPath)
}
