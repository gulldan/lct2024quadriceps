package grpc

import (
	"context"
	"fmt"

	"github.com/gulldan/lct2024_copyright/bff/pkg/config"
	"github.com/rs/zerolog"

	taskcontroller "github.com/gulldan/lct2024_copyright/bff/internal/controller/task_controller"
	bffv1 "github.com/gulldan/lct2024_copyright/bff/proto/bff/v1"
)

type Handler struct {
	bffv1.UnimplementedScanTasksServiceServer

	taskCtl *taskcontroller.TaskController
}

func New(cfg *config.Config, log *zerolog.Logger) *Handler {
	ctl := taskcontroller.New(cfg, log)

	return &Handler{
		taskCtl: ctl,
	}
}

func (h *Handler) CreateTaskFromFile(ctx context.Context, req *bffv1.CreateTaskFromFileRequest) (*bffv1.CreateTaskFromFileResponse, error) {
	return nil, nil
}

func (h *Handler) UploadOriginalVideo(ctx context.Context, req *bffv1.UploadOriginalVideoRequest) (*bffv1.UploadOriginalVideoResponse, error) {
	return nil, nil
}

func (h *Handler) GetTasksPreview(ctx context.Context, req *bffv1.GetTasksPreviewRequest) (*bffv1.GetTasksPreviewResponse, error) {
	tasks, total, err := h.taskCtl.GetTasks(ctx, req.Limit, req.Page*req.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return &bffv1.GetTasksPreviewResponse{
		TasksPreview: tasksPreviewToProto(tasks),
		Total:        uint64(total),
	}, nil
}

func (h *Handler) GetTask(ctx context.Context, req *bffv1.GetTaskRequest) (*bffv1.GetTaskResponse, error) {
	task, err := h.taskCtl.GetTask(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("get task failed: %w", err)
	}

	return taskToProto(task), nil
}
