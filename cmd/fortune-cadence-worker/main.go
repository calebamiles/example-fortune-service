package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"

	"github.com/calebamiles/example-fortune-service/cadence/activity"
	"github.com/calebamiles/example-fortune-service/cadence/workflow"
	"github.com/uber-go/tally"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// HostPort is the location of the cadence frontend
	HostPort = "127.0.0.1:7933"

	// Domain is the namespace to use
	Domain = "HCP"

	// ClientName is the name of the worker
	ClientName = "FortuneWorker"
)

func main() {
	cadenceService := "cadence-frontend"

	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		logger.Fatal("Failed to setup logger", zap.Error(err))
	}

	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		logger.Fatal("Failed to setup tchannel", zap.Error(err))
	}
	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			cadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})

	if err := dispatcher.Start(); err != nil {
		logger.Fatal("Failed to start dispatcher", zap.Error(err))
	}

	service := workflowserviceclient.New(dispatcher.ClientConfig(cadenceService))

	// TaskListName identifies set of client workflows, activities, and workers.
	// It could be your group or client or application name.
	workerOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(workflow.TaskList, map[string]string{}),
	}

	worker := worker.New(
		service,
		Domain,
		workflow.TaskList,
		workerOptions,
	)

	worker.RegisterActivity(activity.GetFortune)
	worker.RegisterWorkflow(workflow.GetFortune)

	err = worker.Start()
	if err != nil {
		logger.Error("Failed to start worker", zap.Error(err))
	}

	logger.Info("Started worker.", zap.String("worker", workflow.TaskList))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Cadence workers are expected to be long runnning; block until signaled
	<-sigs
	logger.Info("Shutting down worker.", zap.String("worker", workflow.TaskList))
	worker.Stop()
}
