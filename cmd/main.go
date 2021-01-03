package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sergeikus/go-rest-template/pkg/conf"
	"github.com/sergeikus/go-rest-template/pkg/handler"
	"github.com/sergeikus/go-rest-template/pkg/storage"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], `[--config <path>] server`)
		flag.PrintDefaults()
	}
	configuration := flag.String("config", "../configs/config.yaml", "Configuration file in YAML format")
	flag.Parse()

	c, err := conf.ReadConf(*configuration)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	if err := c.Validate(); err != nil {
		log.Fatalf("configuration validation failed: %v", err)
	}

	// Change working directory to specify files relativly to the configuration file location
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	if err := os.Chdir(filepath.Dir(filepath.Join(currentDir, *configuration))); err != nil {
		log.Fatalf("failed to change working directory: %v", err)
	}

	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%s", strconv.Itoa(c.Port)),
	}

	log.Printf("Initializing storage...")
	var api handler.API
	switch c.DatabaseType {
	case storage.DatabaseTypeInMemory:
		api = handler.API{DB: &storage.InMemoryStorage{}}
	default:
		log.Fatalf("unsupported database type: '%s'", c.DatabaseType)
	}

	log.Printf("Performing connection to database...")
	if err := api.DB.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Printf("Successfully connected to database")

	// Set handlers
	http.HandleFunc("/api/data/get", api.GetKey)
	http.HandleFunc("/api/data/store", api.AddKey)

	if c.TLS {
		log.Printf("Starting HTTPS server")
		if err := httpServer.ListenAndServeTLS(c.TLSCertPath, c.TLSKeyPath); err != nil {
			log.Fatalf("failed to initialize TLS server: %v", err)
		}
	} else {
		log.Printf("Starting HTTP server")
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed to initialize server: %v", err)
		}
	}

}
