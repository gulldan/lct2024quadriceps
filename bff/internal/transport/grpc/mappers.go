package grpc

import (
	"github.com/gulldan/lct2024_copyright/bff/internal/model"

	bffv1 "github.com/gulldan/lct2024_copyright/bff/proto/bff/v1"
)

func statusToProto(st model.TaskStatus) bffv1.TaskStatus {
	switch st {
	case model.TaskStatusDone:
		return bffv1.TaskStatus_TASK_STATUS_DONE
	case model.TaskStatusFailed:
		return bffv1.TaskStatus_TASK_STATUS_FAIL
	case model.TaskStatusInProgress:
		return bffv1.TaskStatus_TASK_STATUS_IN_PROGRESS
	default:
		return bffv1.TaskStatus_TASK_STATUS_UNSPECIFIED
	}
}

func tasksPreviewToProto(m []model.Task) []*bffv1.TaskPreview {
	p := make([]*bffv1.TaskPreview, len(m))
	for i := range m {
		p[i] = &bffv1.TaskPreview{
			Id:         m[i].TaskID,
			Name:       m[i].VideoName,
			PreviewUrl: m[i].PreviewIDUrl,
			Status:     statusToProto(m[i].Status),
		}
	}

	return p
}

func taskToProto(m model.Task) *bffv1.GetTaskResponse {
	return &bffv1.GetTaskResponse{
		Id:        m.TaskID,
		Name:      m.VideoName,
		VideoUrl:  m.VideoIDUrl,
		Status:    statusToProto(m.Status),
		Copyright: copyrightToProto(m.Copyright),
	}
}

func copyrightToProto(c []model.Copyright) []*bffv1.CopyrightTimestamp {
	p := make([]*bffv1.CopyrightTimestamp, len(c))
	for i := range p {
		p[i] = &bffv1.CopyrightTimestamp{
			CopyrightStart: uint64(c[i].CopyrightStart),
			CopyrightEnd:   uint64(c[i].CopyrightEnd),
			OrigStart:      uint64(c[i].OrigStart),
			OrigEnd:        uint64(c[i].OrigEnd),
			OrigId:         c[i].OrigID,
			OrigUrl:        c[i].OrigURL,
		}
	}

	return p
}
