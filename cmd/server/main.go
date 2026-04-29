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
	var err error
	//if err = logger.Initialize("info"); err != nil {
	if err = logger.Initialize("debug"); err != nil {
		panic(err)
	}

	var storage repository.Storage
	if storeCfg.DB != "" {
		storage, err = repository.NewDBStorage(storeCfg.DB)
		if err != nil {
			logger.Log.Warn("Can't connect to DB", zap.String("db", err.Error()))
			if storeCfg.FilePath != "" {
				storage = repository.NewFileStorage(storeCfg.FilePath)
			} else {
				storage = repository.NewMemStorage()
			}
		}
	} else if storeCfg.FilePath != "" {
		storage = repository.NewFileStorage(storeCfg.FilePath)
	} else {
		storage = repository.NewMemStorage()
	}

	//if storeCfg.Restore && storeCfg.FilePath != "" {
	//	if err := storage.LoadFromFile(storeCfg.FilePath); err == nil {
	//		logger.Log.Info("Restored metrics from file", zap.String("path", storeCfg.FilePath))
	//	}
	//}

	mux := server.Router(storage, storeCfg.DB)

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
			//for range ticker.C {
			//	if err := storage.SaveToFile(storeCfg.FilePath); err == nil {
			//		logger.Log.Info("Metrics saved to file", zap.String("path", storeCfg.FilePath))
			//	}
			//}
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

	//if storeCfg.FilePath != "" {
	//	if err := storage.SaveToFile(storeCfg.FilePath); err == nil {
	//		logger.Log.Info("Metrics saved on shutdown", zap.String("path", storeCfg.FilePath))
	//	}
	//}

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
	dbString := flag.String("d", "", "DB connection string")

	flag.Parse()

	var envCfg EnvConfig

	if err := env.Parse(&envCfg); err != nil {
		panic(err)
	}

	if envCfg.Address != nil {
		if err := addr.Set(*envCfg.Address); err != nil {
			panic(err)
		}
	}

	storeCfg := StoreConfig{
		Interval: time.Duration(DefaultStoreInterval) * time.Second,
		FilePath: DefaultFileStoragePath,
		Restore:  DefaultRestore,
		DB:       "",
	}

	if envCfg.StoreInterval != nil {
		storeCfg.Interval = time.Duration(*envCfg.StoreInterval) * time.Second
	} else {
		storeCfg.Interval = time.Duration(*storeInterval) * time.Second
	}

	if envCfg.FileStorePath != nil {
		storeCfg.FilePath = *envCfg.FileStorePath
	} else {
		storeCfg.FilePath = *fileStoragePath
	}

	if envCfg.Restore != nil {
		storeCfg.Restore = *envCfg.Restore
	} else {
		storeCfg.Restore = *restore
	}

	if storeCfg.FilePath == "" {
		storeCfg.Interval = 0
	}

	if envCfg.DB != nil {
		storeCfg.DB = *envCfg.DB
	} else {
		storeCfg.DB = *dbString
	}

	if err := run(&addr, storeCfg); err != nil {
		panic(err)
	}
}
