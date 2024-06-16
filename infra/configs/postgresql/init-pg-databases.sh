#!/bin/bash

psql -v ON_ERROR_STOP=1 --username postgres <<-EOSQL
  CREATE DATABASE bazadannih;
EOSQL

psql -U postgres -d bazadannih -c "CREATE TYPE task_status AS ENUM ('in_progress', 'fail', 'done');"
psql -U postgres -d bazadannih -c "CREATE TABLE task ( \
  task_id BIGSERIAL PRIMARY KEY, \
  video_name TEXT, \
  video_id TEXT, \
  preview_id TEXT, \
  status task_status, \
  copyright JSONB \
);"
psql -U postgres -d bazadannih -c "CREATE TABLE IF NOT EXISTS origvideo ( \
  video_id TEXT \
);"
