CREATE DATABASE bazadannih;
CREATE TYPE task_status AS ENUM ('in_progress', 'fail', 'done');
CREATE TABLE IF NOT EXISTS task (
  task_id BIGSERIAL PRIMARY KEY,
  video_name TEXT,
  video_id TEXT,
  preview_id TEXT,
  status task_status,
  copyright JSONB
);
CREATE TABLE IF NOT EXISTS origvideo (
  video_id TEXT,
  video_hash TEXT UNIQUE,
  video_minio_id TEXT
)
