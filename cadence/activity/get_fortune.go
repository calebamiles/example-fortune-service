package activity

import (
	"context"

	"github.com/calebamiles/example-fortune-service/fortune"
	"go.uber.org/cadence/activity"
	"go.uber.org/zap"
)

const defaultFortune = `
Software engineering is what happens to programming when you add time and other programmers.

–Russ Cox
`

// GetFortune returns a new fortune, or the default
func GetFortune(ctx context.Context) (string, error) {
	defaultFortune := []byte(defaultFortune)

	f := fortune.NewProvider(defaultFortune)
	txt, err := f.Get()
	if err != nil {
		return "", err
	}

	activity.GetLogger(ctx).Info("GetFortune called.", zap.String("Fortune", string(txt)))
	return string(txt), nil
}
