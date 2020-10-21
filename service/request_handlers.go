package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"net/http"

	"github.com/calebamiles/example-fortune-service/cadence/workflow"
	"github.com/calebamiles/example-fortune-service/provider"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"

	"go.uber.org/cadence/client"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultFortune = `
Software engineering is what happens to programming when you add time and other programmers.

-–Russ Cox
`

const (
	// HostPort is the location of the cadence frontend
	HostPort = "127.0.0.1:7933"

	// Domain is the namespace to use
	Domain = "hcp"

	// ClientName is the name of the worker
	ClientName = "get-fortune-http-handler"

	// ClientService is the Cadence service to connect to
	ClientService = "cadence-frontend"
)

// HandleGetFortuneCadence returns a new fortune through executation of a Cadence workflow
func HandleGetFortuneCadence(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cadenceOpts := &client.Options{
		Identity: ClientName,
	}

	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		logger.Fatal("Failed to setup logger", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		logger.Error("Failed to setup tchannel", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			ClientService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})

	if err := dispatcher.Start(); err != nil {
		logger.Error("Failed to start dispatcher", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	thriftService := workflowserviceclient.New(dispatcher.ClientConfig(ClientService))
	cadence := client.NewClient(thriftService, "hcp", cadenceOpts)

	startOpts := client.StartWorkflowOptions{
		TaskList:                        workflow.TaskList,
		ExecutionStartToCloseTimeout:    60 * time.Second,
		DecisionTaskStartToCloseTimeout: 20 * time.Second,
		WorkflowIDReusePolicy:           client.WorkflowIDReusePolicyAllowDuplicateFailedOnly,
		Memo:                            map[string]interface{}{"workflow-type": "local-development"},
	}

	future, err := cadence.ExecuteWorkflow(ctx, startOpts, workflow.GetFortune)
	if err != nil {
		logger.Error("error executing GetFortune workflow: %s", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var fortune string
	err = future.Get(ctx, &fortune)
	if err != nil {
		logger.Error("error getting GetFortune workflow result: %s", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write([]byte(fortune))
	if err != nil {
		logger.Error("error: writing HandleGetFortune response: %s", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n != len(fortune) {
		logger.Error(fmt.Sprintf("error: expected to write %d bytes, but only wrote %d", len(fortune), n))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// HandleGetFortuneDirect returns a new fortune through direct use of a provider.Fortune
func HandleGetFortuneDirect(w http.ResponseWriter, req *http.Request) {
	defaultFortune := []byte(defaultFortune)

	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		logger.Error("Failed to setup logger", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fortune := provider.NewFortune(defaultFortune)
	rawTxt, err := fortune.Get()
	if err != nil {
		logger.Error("Failed to get fortune", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(rawTxt)
	if err != nil {
		logger.Error("error: writing HandleGetFortuneDirect response: %s", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n != len(rawTxt) {
		logger.Error(fmt.Sprintf("error: expected to write %d bytes, but only wrote %d", len(rawTxt), n))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// HandleGetHealthz returns server health as ok
func HandleGetHealthz(w http.ResponseWriter, req *http.Request) {
	config := zap.NewDevelopmentConfig()
	config.Level.SetLevel(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		logger.Error("Failed to setup logger", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response struct {
		Status string `json:"status"`
	}

	response.Status = "ok"
	responseJSON, err := json.Marshal(response)
	if err != nil {
		logger.Error("encoding status to JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n, err := w.Write(responseJSON)
	if err != nil {
		logger.Error("writing response: %s", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if n != len(responseJSON) {
		logger.Error(fmt.Sprintf("expected to write %d bytes, but only wrote %d", len(responseJSON), n))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
