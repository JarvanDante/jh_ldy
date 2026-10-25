package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/JarvanDante/jh_ldy/internal/conf"
	"github.com/JarvanDante/jh_ldy/internal/httpserver"
)

func main() {
	root := os.Getenv("JH_LDY_ROOT")
	if root == "" {
		if wd, err := os.Getwd(); err == nil {
			root = wd
		} else {
			root = "."
		}
	}
	cfgPath := os.Getenv("JH_LDY_CONFIG")
	if cfgPath == "" {
		cfgPath = filepath.Join(root, "config.json")
	}
	cfg, err := conf.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	srv := httpserver.New(cfg, root)
	log.Printf("jh_ldy listen %s", cfg.Listen)
	if err = http.ListenAndServe(cfg.Listen, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
