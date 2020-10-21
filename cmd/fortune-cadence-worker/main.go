package main

import (
	"net/http"

	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"

	"github.com/calebamiles/example-fortune-service/cadence/activity"
	"github.com/calebamiles/example-fortune-service/cadence/workflow"
	"github.com/calebamiles/example-fortune-service/service"
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
	Domain = "hcp"

	// ClientName is the name of the worker
	ClientName = "fortune-worker"
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

	thriftService := workflowserviceclient.New(dispatcher.ClientConfig(cadenceService))

	// TaskListName identifies set of client workflows, activities, and workers.
	// It could be your group or client or application name.
	workerOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(workflow.TaskList, map[string]string{}),
	}

	worker := worker.New(
		thriftService,
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

	logger.Info("Started Cadence worker.", zap.String("worker", workflow.TaskList))

	http.HandleFunc("/fortune", service.HandleGetFortuneCadence)
	http.HandleFunc("/healthz", service.HandleGetHealthz)

	logger.Info("Starting HTTP server.", zap.String("service", "fortune"))
	http.ListenAndServe(":8080", nil)
}
