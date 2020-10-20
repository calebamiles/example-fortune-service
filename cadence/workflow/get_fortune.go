package workflow

import (
	"time"

	"github.com/calebamiles/example-fortune-service/cadence/activity"
	"go.uber.org/cadence/workflow"
	"go.uber.org/zap"
)

func init() {
	workflow.Register(GetFortune)
}

const (
	// TaskList is the queue to use for GetFortune execution
	TaskList = "getFortuneTaskList"
)

// GetFortune returns a new fortune, or the default fortune
func GetFortune(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		TaskList:               TaskList,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	future := workflow.ExecuteActivity(ctx, activity.GetFortune)
	var result string
	if err := future.Get(ctx, &result); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Done", zap.String("result", result))
	return nil
}
