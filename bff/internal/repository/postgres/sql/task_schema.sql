CREATE TYPE task_status AS ENUM ('in_progress', 'fail', 'done');

CREATE TABLE task (
  task_id BIGSERIAL PRIMARY KEY,
  video_name TEXT,
  video_id TEXT,
  preview_id TEXT,
  status task_status,
  copyright JSONB
);

CREATE TABLE origvideo (
  video_id TEXT,
  video_hash TEXT,
  video_minio_id TEXT
);
