package main

import (
	"errors"
	"fmt"
	"github.com/brcodingdev/chat-app/service/pkg/app"
	"github.com/brcodingdev/chat-app/service/pkg/broker"
	"github.com/brcodingdev/chat-app/service/port/api"
	"github.com/brcodingdev/chat-app/service/port/http/middleware"
	"github.com/brcodingdev/chat-app/service/port/http/route"
	"github.com/brcodingdev/chat-app/service/port/repository"
	"github.com/brcodingdev/chat-app/service/server"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	stdlog "log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run() error {
	// configure logs
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)

	level.Info(logger).Log("service starting...")

	// setup database and migrate model
	db, err := repository.Connect()

	if err != nil {
		return err
	}
	// migrate DB with table models
	err = repository.MigrateDB()

	if err != nil {
		return err
	}

	level.Info(logger).Log("Database", "migrated")

	// connect RabbitMQ
	rmqHost := os.Getenv("RABBIT_HOST")
	rmqUserName := os.Getenv("RABBIT_USERNAME")
	rmqPassword := os.Getenv("RABBIT_PASSWORD")
	rmqPort := os.Getenv("RABBIT_PORT")

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		rmqUserName,
		rmqPassword,
		rmqHost,
		rmqPort,
	)
	fmt.Println("DSN ", dsn)

	receiverQueue := os.Getenv("RECEIVER_QUEUE")
	publisherQueue := os.Getenv("PUBLISHER_QUEUE")

	fmt.Println(
		"RECEIVER_QUEUE ", receiverQueue,
		"PUBLISHER_QUEUE", publisherQueue,
	)

	if receiverQueue == "" ||
		publisherQueue == "" {
		stdlog.Println("required RECEIVER_QUEUE, PUBLISHER_QUEUE env vars set")
		return errors.New("could not init the server ")
	}
	// register RabbitMQ
	rabbit, err := broker.NewRabbitMQ(
		dsn,
		receiverQueue,
		publisherQueue)

	if err != nil {
		stdlog.Println("could not initialize RabbitMQ")
		return err
	}

	defer rabbit.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("JWT Secret not set")
	}
	// register repositories
	chatRoomRepository := repository.NewChatRoomDB(db)
	chatRepository := repository.NewChatDB(db)
	chatApp := app.NewChatApp(
		chatRoomRepository,
		chatRepository,
	)

	r := route.RegisterRoute()

	userRepository := repository.NewUserDB(db)
	authApp := app.NewAuthApp(userRepository)
	// register authorization routes
	authAPI := api.NewAuthAPI(authApp)
	authProvider := route.NewAuthProvider(r, authAPI)
	authProvider.RegisterRoutes()
	// register chat routes
	chatAPI := api.NewChatAPI(chatApp)
	chatProvider := route.NewChatProvider(r, chatAPI)
	chatProvider.RegisterRoutes()
	// register ws routes
	wsProvider := route.NewWebSocketProvide(r, chatApp, rabbit)
	wsProvider.RegisterRoutes()
	// route with logging and cors middleware
	loggingMiddleware := middleware.RegisterLoggingMiddleware(logger)
	loggedRoutes := loggingMiddleware(r)
	handler := middleware.Cors(loggedRoutes)
	// start api appServer
	port := os.Getenv("PORT")
	appServer := server.NewServer(port, handler)

	// handle graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("received interrupt signal. shutting down gracefully...")
		// Perform any cleanup or shutdown logic here

		// close the RabbitMQ connection
		if err := rabbit.Close(); err != nil {
			fmt.Println("could not close RabbitMQ:", err)
		}

		// close the database connection
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}

		fmt.Println("graceful shutdown completed")
		os.Exit(0)
	}()

	// start server
	level.Info(logger).Log("Server", "starting", "port", port)
	defer appServer.Shutdown()
	err = appServer.Serve()
	return err
}
