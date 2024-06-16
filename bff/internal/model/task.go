package model

type TaskStatus uint

const (
	TaskStatusDone TaskStatus = iota
	TaskStatusInProgress
	TaskStatusFailed
)

type Copyright struct {
	CopyrightStart int    `json:"copyright_start,omitempty"`
	CopyrightEnd   int    `json:"copyright_end,omitempty"`
	OrigStart      int    `json:"orig_start,omitempty"`
	OrigEnd        int    `json:"orig_end,omitempty"`
	OrigID         string `json:"orig_id,omitempty"`
}

type Task struct {
	TaskID       int64
	VideoName    string
	VideoIDUrl   string
	PreviewIDUrl string
	Status       TaskStatus
	Copyright    []Copyright
}
