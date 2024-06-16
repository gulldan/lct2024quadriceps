package taskcontroller

import (
	"encoding/json"

	"github.com/gulldan/lct2024_copyright/bff/internal/model"

	pgsql "github.com/gulldan/lct2024_copyright/bff/internal/repository/postgres"
)

func statusToModel(s pgsql.TaskStatus) model.TaskStatus {
	switch s {
	case pgsql.TaskStatusDone:
		return model.TaskStatusDone
	case pgsql.TaskStatusFail:
		return model.TaskStatusFailed
	case pgsql.TaskStatusInProgress:
		return model.TaskStatusInProgress
	default:
		return model.TaskStatusFailed
	}
}

func taskSliceToModel(t []pgsql.Task) ([]model.Task, error) {
	m := make([]model.Task, len(t))
	var err error

	for i := range t {
		m[i], err = taskToModel(t[i])
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func taskToModel(t pgsql.Task) (model.Task, error) {
	var cArr []model.Copyright
	var c model.Copyright
	if err := json.Unmarshal(t.Copyright, &c); err != nil {
		c = model.Copyright{}
	}

	if c.OrigID != "" {
		cArr = append(cArr, c)
	}

	return model.Task{
		VideoName:    t.VideoName.String,
		TaskID:       t.TaskID,
		VideoIDUrl:   t.VideoID.String,
		PreviewIDUrl: t.PreviewID.String,
		Status:       statusToModel(t.Status.TaskStatus),
		Copyright:    cArr,
	}, nil
}
