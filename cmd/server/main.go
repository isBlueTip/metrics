package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/isBlueTip/metrics/internal/logger"
	"github.com/isBlueTip/metrics/internal/repository"
	"github.com/isBlueTip/metrics/internal/server"
	"go.uber.org/zap"
)

const (
	DefaultStoreInterval   = 300
	DefaultFileStoragePath = "/tmp/metrics-db.json"
	DefaultRestore         = true
)

func run(address *ServerAddress, storeCfg StoreConfig) error {
	if err := logger.Initialize("info"); err != nil {
		panic(err)
	}

	storage := repository.NewStorage()

	if storeCfg.Restore && storeCfg.FilePath != "" {
		if err := storage.LoadFromFile(storeCfg.FilePath); err == nil {
			logger.Log.Info("Restored metrics from file", zap.String("path", storeCfg.FilePath))
		}
	}

	mux := server.Router(storage)

	addrString := net.JoinHostPort(address.Host, address.Port)

	logger.Log.Info("Running server", zap.String("address", address.String()))

	srv := &http.Server{
		Addr:    addrString,
		Handler: mux,
	}

	go func() {
		if storeCfg.FilePath != "" && storeCfg.Interval > 0 {
			ticker := time.NewTicker(storeCfg.Interval)
			defer ticker.Stop()
			for range ticker.C {
				if err := storage.SaveToFile(storeCfg.FilePath); err == nil {
					logger.Log.Info("Metrics saved to file", zap.String("path", storeCfg.FilePath))
				}
			}
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if storeCfg.FilePath != "" {
		if err := storage.SaveToFile(storeCfg.FilePath); err == nil {
			logger.Log.Info("Metrics saved on shutdown", zap.String("path", storeCfg.FilePath))
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func main() {
	addr := ServerAddress{Host: "", Port: "8080"}
	flag.Var(&addr, "a", "Parsing net address")

	storeInterval := flag.Int("i", DefaultStoreInterval, "Store interval in seconds")
	fileStoragePath := flag.String("f", DefaultFileStoragePath, "File storage path")
	restore := flag.Bool("r", DefaultRestore, "Restore metrics on startup")

	flag.Parse()

	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Address != nil {
		if err := addr.Set(*cfg.Address); err != nil {
			panic(err)
		}
	}

	storeCfg := StoreConfig{
		Interval:  time.Duration(DefaultStoreInterval) * time.Second,
		FilePath: DefaultFileStoragePath,
		Restore: DefaultRestore,
	}

	if cfg.StoreInterval != nil {
		storeCfg.Interval = time.Duration(*cfg.StoreInterval) * time.Second
	} else {
		storeCfg.Interval = time.Duration(*storeInterval) * time.Second
	}

	if cfg.FileStorePath != nil {
		storeCfg.FilePath = *cfg.FileStorePath
	} else {
		storeCfg.FilePath = *fileStoragePath
	}

	if cfg.Restore != nil {
		storeCfg.Restore = *cfg.Restore
	} else {
		storeCfg.Restore = *restore
	}

	if storeCfg.FilePath == "" {
		storeCfg.Interval = 0
	}

	if err := run(&addr, storeCfg); err != nil {
		panic(err)
	}
}
