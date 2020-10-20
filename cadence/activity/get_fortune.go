package activity

import (
	"context"

	"github.com/calebamiles/example-fortune-service/provider"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

func init() {
	activity.Register(GetFortune)
}

const defaultFortune = `
Software engineering is what happens to programming when you add time and other programmers.

–Russ Cox
`

// GetFortune returns a new fortune, or the default
func GetFortune(ctx context.Context) (string, error) {
	defaultFortune := []byte(defaultFortune)

	fortune := provider.NewFortune(defaultFortune)
	txt, err := fortune.Get()
	if err != nil {
		return "", err
	}

	activity.GetLogger(ctx).Info("GetFortune called.", zap.String("Fortune", string(txt)))
	return string(txt), nil
}
