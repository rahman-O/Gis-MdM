package application

import (
	"context"
	"log/slog"
	"time"

	"github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/domain"
	pluginpostgres "github.com/gis-mdm/server-backend-go/internal/modules/plugins/push/adapter/persistence/postgres"
	notifpostgres "github.com/gis-mdm/server-backend-go/internal/modules/notifications/adapter/persistence/postgres"
)

// ScheduleRunner processes due plugin_push_schedule rows.
type ScheduleRunner struct {
	schedule *pluginpostgres.ScheduleRepository
	queue    *notifpostgres.QueueRepository
	repo     *pluginpostgres.MessageRepository
	log      *slog.Logger
}

func NewScheduleRunner(
	schedule *pluginpostgres.ScheduleRepository,
	queue *notifpostgres.QueueRepository,
	repo *pluginpostgres.MessageRepository,
	log *slog.Logger,
) *ScheduleRunner {
	return &ScheduleRunner{schedule: schedule, queue: queue, repo: repo, log: log}
}

func (r *ScheduleRunner) RunOnce(ctx context.Context) {
	tasks, err := r.schedule.ListAll(ctx)
	if err != nil {
		if r.log != nil {
			r.log.Warn("schedule list failed", slog.String("error", err.Error()))
		}
		return
	}
	now := time.Now()
	for _, task := range tasks {
		if !MatchesSchedule(task, now) {
			continue
		}
		r.sendTask(ctx, task)
	}
}

func (r *ScheduleRunner) sendTask(ctx context.Context, task domain.PluginPushSchedule) {
	deviceIDs, err := r.schedule.ResolveDeviceIDs(ctx, task)
	if err != nil {
		if r.log != nil {
			r.log.Warn("schedule resolve devices", slog.Int64("task", task.ID), slog.String("error", err.Error()))
		}
		return
	}
	payload := task.Payload
	for _, deviceID := range deviceIDs {
		if deviceID <= 0 {
			continue
		}
		if err := r.queue.Enqueue(ctx, deviceID, task.MessageType, payload); err != nil && r.log != nil {
			r.log.Warn("schedule enqueue", slog.Int64("device", deviceID), slog.String("error", err.Error()))
			continue
		}
		_ = r.repo.InsertHistory(ctx, task.CustomerID, deviceID, task.MessageType, payload)
	}
}
