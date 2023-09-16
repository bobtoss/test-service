package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"test-service/internal/cache"
	"test-service/internal/config"
	"test-service/internal/handler"
	"test-service/internal/repository"
	"test-service/internal/service"
	"test-service/pkg/log"
	"test-service/pkg/server"
	"time"

	"go.uber.org/zap"
)

// Run initializes whole application.
func Run() {
	logger := log.LoggerFromContext(context.Background())

	configs, err := config.New()
	if err != nil {
		logger.Error("ERR_INIT_CONFIGS", zap.Error(err))
		return
	}

	repositories, err := repository.New(
		repository.WithMongoStore(configs.MONGO.DSN, "test"))
	if err != nil {
		logger.Error("ERR_INIT_REPOSITORIES", zap.Error(err))
		return
	}
	defer repositories.Close()

	caches, err := cache.New(
		cache.Dependencies{
			UserRepository: repositories.User,
		},
		cache.WithRedisStore("redis://red-ck29ag021fec73d6mg0g:6379"))
	if err != nil {
		logger.Error("ERR_INIT_CACHES", zap.Error(err))
		return
	}
	defer caches.Close()

	userService, err := service.New(
		service.WithUserRepository(repositories.User),
		service.WithUserCache(caches.User))

	if err != nil {
		logger.Error("ERR_INIT_USER_SERVICE", zap.Error(err))
		return
	}

	handlers, err := handler.New(
		handler.Dependencies{
			Configs:     configs,
			UserService: userService,
		},
		handler.WithHTTPHandler())
	if err != nil {
		logger.Error("ERR_INIT_HANDLERS", zap.Error(err))
		return
	}

	servers, err := server.New(
		server.WithHTTPServer(handlers.HTTP, configs.APP.Port))
	if err != nil {
		logger.Error("ERR_INIT_SERVERS", zap.Error(err))
		return
	}

	// Run our server in a goroutine so that it doesn't block.
	if err = servers.Run(logger); err != nil {
		logger.Error("ERR_RUN_SERVERS", zap.Error(err))
		return
	}
	logger.Info("http server started on http://localhost:" + configs.APP.Port + "/swagger/index.html")

	// Graceful Shutdown
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the httpServer gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	quit := make(chan os.Signal, 1) // create channel to signify a signal being sent

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel
	<-quit                                             // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")

	// create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err = servers.Stop(ctx); err != nil {
		panic(err) // failure/timeout shutting down the httpServer gracefully
	}

	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here

	fmt.Println("Server was successful shutdown.")
}
