package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/sergeikus/go-rest-template/pkg/conf"
	"github.com/sergeikus/go-rest-template/pkg/handler"
	"github.com/sergeikus/go-rest-template/pkg/storage"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], `[--config <path>]
		`)
		flag.PrintDefaults()
	}
	configuration := flag.String("config", "../configs/config.yaml", "Path to server configuration (supported format is YAML)")
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
	// This environmental variable will override configuration
	inMemoryOverride := os.Getenv("DB_TYPE_INMEMORY")
	if strings.ToLower(inMemoryOverride) == "true" {
		c.Database.Type = "in-memory"
	}

	var api handler.API
	log.Printf("Database type is: %s", c.Database.Type)
	switch c.Database.Type {
	case storage.DatabaseTypeInMemory:
		api = handler.API{DB: &storage.InMemoryStorage{}}
	case storage.DatabaseTypePostgre:
		api = handler.API{
			DB: storage.DefinePostgresStorage(
				c.Database.Username, c.Database.Password, c.Database.Name, c.Database.Host, c.Database.Port,
			),
		}
	default:
		log.Fatalf("unsupported database type: '%s'", c.Database.Type)
	}

	log.Printf("Performing connection to database...")
	if err := api.DB.Connect(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Printf("Successfully connected to database")
	// Close connection when application shuts down
	defer api.DB.Close()

	// Set handlers
	http.HandleFunc("/api/data/get", api.GetData)
	http.HandleFunc("/api/data/get/all", api.GetAllData)
	http.HandleFunc("/api/data/store", api.Store)

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
