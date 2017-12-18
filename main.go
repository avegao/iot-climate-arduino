package main

import (
	"github.com/avegao/gocondi"
	"os"
	"os/signal"
	"syscall"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/heroku/rollrus"
	"net"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	pb "github.com/avegao/iot-climate-arduino/proto"
	"github.com/avegao/iot-climate-arduino/controller"
)

const version = "1.0.0"

var (
	debug      = flag.Bool("debug", false, "Print debug logs")
	grpcPort   = flag.Int("port", 50000, "gRPC port. Default 50000")
	buildDate  string
	commitHash string
	container  *gocondi.Container
	parameters map[string]interface{}
	server     *grpc.Server
)

func initContainer() {
	flag.Parse()

	parameters = map[string]interface{}{
		"build_date":    buildDate,
		"commit_hash":   commitHash,
		"debug":         *debug,
		"rollbar_token": "",
		"version":       version,
	}

	logger := initLogger()
	gocondi.Initialize(logger)
	container = gocondi.GetContainer()

	for name, value := range parameters {
		container.SetParameter(name, value)
	}
}

func initLogger() *logrus.Logger {
	logLevel := logrus.InfoLevel
	environment := "release"

	if *debug {
		logLevel = logrus.DebugLevel
		environment = "debug"
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetFormatter(&logrus.TextFormatter{})

	hook := rollrus.NewHook(fmt.Sprintf("%v", parameters["rollbar_token"]), environment)

	log := logrus.New()
	log.SetLevel(logLevel)
	log.Hooks.Add(hook)

	return log
}

func initGrpc() {
	logger := container.GetLogger()
	logger.Debugf("initGrpc() - START")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))

	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	logger.Debugf("gRPC listening in %d port", *grpcPort)

	server = grpc.NewServer()
	pb.RegisterArduinoServer(server, controller.Controller{})
	reflection.Register(server)

	if err := server.Serve(listen); err != nil {
		logger.Fatalf("failed to server: %v", err)
	}

	logger.Debugf("initGrpc() - END")
}

func powerOff(exitCode int) {
	if nil != gocondi.GetContainer() {
		gocondi.GetContainer().GetLogger().Infof("Shutting down...")
	}

	os.Exit(exitCode)
}

func handleInterrupt() {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		<-gracefulStop
		powerOff(0)
	}()
}

func main() {
	initContainer()
	handleInterrupt()
	defer powerOff(0)

	logger := container.GetLogger()
	logger.Infof(
		"Dashboard Mobile Quality Query Exporter for Spain started v%s (commit %s, build date %s)",
		container.GetStringParameter("version"),
		container.GetStringParameter("commit_hash"),
		container.GetStringParameter("build_date"),
	)

	// this lock the execution
	initGrpc()
}
