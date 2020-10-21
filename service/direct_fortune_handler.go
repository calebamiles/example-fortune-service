package service

import (
	"fmt"
	"net/http"

	"github.com/calebamiles/example-fortune-service/provider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultFortune = `
Software engineering is what happens to programming when you add time and other programmers.

-–Russ Cox
`

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
