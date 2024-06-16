package taskcontroller

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gulldan/lct2024_copyright/bff/internal/model"
	"github.com/gulldan/lct2024_copyright/bff/internal/pkg/ffmpeg"
	"github.com/gulldan/lct2024_copyright/bff/internal/pkg/videocopy"
	"github.com/gulldan/lct2024_copyright/bff/internal/pkg/wav2vec"
	"github.com/gulldan/lct2024_copyright/bff/internal/repository/minio"
	"github.com/gulldan/lct2024_copyright/bff/pkg/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/xid"
	"github.com/rs/zerolog"

	pgsql "github.com/gulldan/lct2024_copyright/bff/internal/repository/postgres"
)

var (
	once       sync.Once
	controller *TaskController
)

type TaskController struct {
	ffmpegExec    *ffmpeg.FfmpegExecutor
	minioClient   *minio.MinioClient
	log           *zerolog.Logger
	pgConn        *pgsql.Queries
	wav2vecClient *wav2vec.Client
	videocopy     *videocopy.Client
}

func New(cfg *config.Config, log *zerolog.Logger) *TaskController {
	once.Do(func() {
		httpCl := http.DefaultClient
		httpCl.Timeout = time.Hour

		pg, err := pgx.Connect(context.Background(), cfg.Postgres.Addr)
		if err != nil {
			panic(fmt.Errorf("failed to connect to postgres: %w", err))
		}

		m, err := minio.NewMinioClient(&cfg.Minio)
		if err != nil {
			panic(fmt.Errorf("failed to create minio client: %w", err))
		}

		w2vecCl, err := wav2vec.NewClient(cfg.Wav2VecAddr, wav2vec.WithHTTPClient(httpCl))
		if err != nil {
			panic(fmt.Errorf("failed to connect to wav2vec: %w", err))
		}

		videocopyCl, err := videocopy.NewClient(cfg.VideocopyAddr, videocopy.WithHTTPClient(httpCl))
		if err != nil {
			panic(fmt.Errorf("failed to connect to wav2vec: %w", err))
		}

		controller = &TaskController{
			minioClient:   m,
			log:           log,
			pgConn:        pgsql.New(pg),
			ffmpegExec:    ffmpeg.New(log),
			wav2vecClient: w2vecCl,
			videocopy:     videocopyCl,
		}
	})

	return controller
}

func (ctl *TaskController) CreateTask(ctx context.Context, file io.Reader, filename string) (int64, error) {
	fmt.Println(file, filename)
	videoID, previewID, videoLen, err := ctl.makePreviewUploadVideo(ctx, file)
	if err != nil {
		return 0, fmt.Errorf("failed to upload video: %w", err)
	}

	fmt.Println(videoID, previewID, videoLen, err)

	hash, err := ctl.getHashFromVideo(ctx, videoID, ctl.minioClient.GetVideoBucketName())
	if err != nil {
		return 0, fmt.Errorf("failed to calculate hash for video: %w", err)
	}

	videos, err := ctl.pgConn.GetOrigVideosByHash(ctx, pgtype.Text{
		String: hash,
		Valid:  true,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to compare hash with original videos: %w", err)
	}

	if len(videos) != 0 {
		//nolint:shadow
		task, errC := ctl.pgConn.CreateTask(ctx, pgsql.CreateTaskParams{
			VideoID: pgtype.Text{
				String: videoID,
				Valid:  true,
			},
			PreviewID: pgtype.Text{
				String: previewID,
				Valid:  true,
			},
			Status: pgsql.NullTaskStatus{
				TaskStatus: pgsql.TaskStatusDone,
				Valid:      true,
			},
			VideoName: pgtype.Text{
				String: filename,
				Valid:  true,
			},
		})
		if errC != nil {
			return 0, fmt.Errorf("create task failed: %w", err)
		}

		c := model.Copyright{
			CopyrightStart: 0,
			CopyrightEnd:   videoLen,
			OrigStart:      0,
			OrigEnd:        videoLen,
			OrigID:         videos[0].VideoID.String,
		}

		copyright, errC := json.Marshal(c)
		if errC != nil {
			return 0, fmt.Errorf("failed to marshal copyright to json: %w", err)
		}

		if errC = ctl.pgConn.UpdateTaskCopyright(ctx, pgsql.UpdateTaskCopyrightParams{
			TaskID:    task.TaskID,
			Copyright: copyright,
		}); errC != nil {
			return 0, fmt.Errorf("failed to update task copyright: %w", err)
		}

		return task.TaskID, nil
	}

	task, err := ctl.pgConn.CreateTask(ctx, pgsql.CreateTaskParams{
		VideoID: pgtype.Text{
			String: videoID,
			Valid:  true,
		},
		PreviewID: pgtype.Text{
			String: previewID,
			Valid:  true,
		},
		Status: pgsql.NullTaskStatus{
			TaskStatus: pgsql.TaskStatusInProgress,
			Valid:      true,
		},
		VideoName: pgtype.Text{
			String: filename,
			Valid:  true,
		},
	})
	if err != nil {
		return 0, fmt.Errorf("create task failed: %w", err)
	}

	go func() {
		if err := ctl.checkForCopyright(context.Background(), task); err != nil {
			ctl.log.Error().Err(err).Any("task", task).Msg("check for copyright failed")
		}
	}()

	return task.TaskID, nil
}

func (ctl *TaskController) checkForCopyright(ctx context.Context, task pgsql.Task) error {
	if err := ctl.pgConn.UpdateTaskStatus(ctx, pgsql.UpdateTaskStatusParams{
		TaskID: task.TaskID,
		Status: pgsql.NullTaskStatus{
			TaskStatus: pgsql.TaskStatusInProgress,
			Valid:      true,
		},
	}); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	audio, video, errR := ctl.runVideoAudioCopyrightCheck(ctx, task)
	if errR != nil {
		if errK := ctl.pgConn.UpdateTaskStatus(ctx, pgsql.UpdateTaskStatusParams{
			TaskID: task.TaskID,
			Status: pgsql.NullTaskStatus{
				TaskStatus: pgsql.TaskStatusFail,
				Valid:      true,
			},
		}); errK != nil {
			ctl.log.Error().Err(errK).Msg("failed to update task status")
		}

		return fmt.Errorf("failed to run video audio copyright: %w", errR)
	}

	var copyright model.Copyright
	if video.PiracyStartFrame == -1 && audio.IDLicenseWav2vec != "-" {
		segOrig := strings.Split(audio.SegmentsWav2vec, "-")
		origStart, err := strconv.Atoi(segOrig[0])
		if err != nil {
			return fmt.Errorf("segOrig to int failed: %s %w", audio.SegmentsWav2vec, err)
		}

		origEnd, err := strconv.Atoi(segOrig[1])
		if err != nil {
			return fmt.Errorf("segEnd to int failed: %s %w", audio.SegmentsWav2vec, err)
		}

		segCopy := strings.Split(audio.SegmentWav2vec, "-")
		copyStart, err := strconv.Atoi(segCopy[0])
		if err != nil {
			return fmt.Errorf("segCopy to int failed: %s %w", audio.SegmentWav2vec, err)
		}

		copyEnd, err := strconv.Atoi(segCopy[1])
		if err != nil {
			return fmt.Errorf("segEnd to int failed: %s %w", audio.SegmentWav2vec, err)
		}

		copyright.OrigID = audio.IDLicenseWav2vec
		copyright.OrigStart = origStart
		copyright.OrigEnd = origEnd
		copyright.CopyrightStart = copyStart
		copyright.CopyrightEnd = copyEnd
	} else {
		copyright.OrigID = video.LicenceName
		copyright.OrigStart = video.LicenseStartFrame
		copyright.OrigEnd = video.LicenseEndFrame
		copyright.CopyrightStart = video.PiracyStartFrame
		copyright.CopyrightEnd = video.PiracyEndFrame
	}

	if copyright.OrigID != "" {
		copyright, errC := json.Marshal(copyright)
		if errC != nil {
			return fmt.Errorf("failed to marshal copyright to json: %w", errC)
		}

		if errC = ctl.pgConn.UpdateTaskCopyright(ctx, pgsql.UpdateTaskCopyrightParams{
			TaskID:    task.TaskID,
			Copyright: copyright,
		}); errC != nil {
			return fmt.Errorf("failed to update task copyright: %w", errC)
		}
	}

	if err := ctl.pgConn.UpdateTaskStatus(ctx, pgsql.UpdateTaskStatusParams{
		TaskID: task.TaskID,
		Status: pgsql.NullTaskStatus{
			TaskStatus: pgsql.TaskStatusDone,
			Valid:      true,
		},
	}); err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

func (ctl *TaskController) runVideoAudioCopyrightCheck(ctx context.Context, task pgsql.Task) (wav2vec.CopyrightAnswer, videocopy.SearchResponse, error) {
	errCh := make(chan error, 2)
	defer close(errCh)
	var (
		wg             sync.WaitGroup
		audioCopyright wav2vec.CopyrightAnswer
		videoCopyright videocopy.SearchResponse
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		audioCopyright, err = ctl.getAudioCopyright(ctx, task)
		if err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		videoCopyright, err = ctl.getVideoCopyright(ctx, task)
		if err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	var retErr error
	for err := range errCh {
		ctl.log.Error().Err(err).Msg("run copyright failed")
		retErr = err
	}

	if retErr != nil {
		return wav2vec.CopyrightAnswer{}, videocopy.SearchResponse{}, retErr
	}

	return audioCopyright, videoCopyright, nil
}

func (ctl *TaskController) getVideoCopyright(ctx context.Context, task pgsql.Task) (videocopy.SearchResponse, error) {
	videoReader, err := ctl.minioClient.GetFileReader(ctx, task.VideoID.String, ctl.minioClient.GetVideoBucketName())
	if err != nil {
		return videocopy.SearchResponse{}, err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	if fw, err = w.CreateFormFile("video_file", task.VideoName.String); err != nil {
		return videocopy.SearchResponse{}, err
	}

	if _, err = io.Copy(fw, videoReader); err != nil {
		return videocopy.SearchResponse{}, err
	}

	resp, err := ctl.videocopy.SearchFindVideoPostWithBody(ctx, w.FormDataContentType(), &b)
	if err != nil {
		return videocopy.SearchResponse{}, fmt.Errorf("failed to find video infringement: %w", err)
	}
	defer resp.Body.Close()

	var copyright videocopy.SearchResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return videocopy.SearchResponse{}, fmt.Errorf("read all failed: %w", err)
	}

	ctl.log.Info().Str("resp_body", string(body)).Msg("video infringement")

	if err = json.Unmarshal(body, &copyright); err != nil {
		return videocopy.SearchResponse{}, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return copyright, nil
}

func (ctl *TaskController) getAudioCopyright(ctx context.Context, task pgsql.Task) (wav2vec.CopyrightAnswer, error) {
	videoReader, err := ctl.minioClient.GetFileReader(ctx, task.VideoID.String, ctl.minioClient.GetVideoBucketName())
	if err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}

	tmpfile, err := os.CreateTemp("", "*.wav")
	if err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	if _, err = io.Copy(tmpfile, videoReader); err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}

	audioFileName, err := ctl.ffmpegExec.GetAudioFromVideo(tmpfile.Name())
	if err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}

	audioFile, err := os.Open(audioFileName)
	if err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}
	defer os.Remove(audioFileName)
	defer audioFile.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	if fw, err = w.CreateFormFile("audio_file", task.VideoName.String); err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}

	if _, err = io.Copy(fw, audioFile); err != nil {
		return wav2vec.CopyrightAnswer{}, err
	}

	resp, err := ctl.wav2vecClient.UpdateDataBaseFindCopyrightInfringementPostWithBody(ctx, w.FormDataContentType(), &b)
	if err != nil {
		return wav2vec.CopyrightAnswer{}, fmt.Errorf("failed to find audio infringement: %w", err)
	}
	defer resp.Body.Close()

	var copyright wav2vec.CopyrightAnswer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return wav2vec.CopyrightAnswer{}, fmt.Errorf("read all failed: %w", err)
	}

	ctl.log.Info().Str("resp_body", string(body)).Msg("audio infringement")

	if err = json.Unmarshal(body, &copyright); err != nil {
		return wav2vec.CopyrightAnswer{}, fmt.Errorf("json unmarshal failed: %w", err)
	}

	return copyright, nil
}

func (ctl *TaskController) UploadOriginalVideo(ctx context.Context, file io.Reader, filename string) (string, error) {
	origVideoID, err := ctl.uploadOriginalVideoToMinio(ctx, file)
	if err != nil {
		return "", err
	}

	hash, err := ctl.getHashFromVideo(ctx, origVideoID, ctl.minioClient.GetOrigVideoBucket())
	if err != nil {
		return "", fmt.Errorf("failed to get hash from orig video: %w", err)
	}

	if err = ctl.updateOriginalEmbeddings(ctx, origVideoID, filename); err != nil {
		return "", fmt.Errorf("failed to update original embeddings: %w", err)
	}

	if _, err = ctl.pgConn.CreateOrigVideo(ctx, pgsql.CreateOrigVideoParams{
		VideoID: pgtype.Text{
			String: filename,
			Valid:  true,
		},
		VideoHash: pgtype.Text{
			String: hash,
			Valid:  true,
		},
		VideoMinioID: pgtype.Text{
			String: origVideoID,
			Valid:  true,
		},
	}); err != nil {
		return "", nil
	}

	return origVideoID, nil
}

func (ctl *TaskController) updateOriginalEmbeddings(ctx context.Context, origVideoID, filename string) error {
	errCh := make(chan error, 2)
	defer close(errCh)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ctl.updateOriginalVideoEmbeddings(ctx, origVideoID, filename); err != nil {
			errCh <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ctl.updateOriginalAudioEmbeddings(ctx, origVideoID, filename); err != nil {
			errCh <- err
		}
	}()

	wg.Wait()
	var retErr error
	for err := range errCh {
		ctl.log.Error().Err(err).Msg("update original embeddings failed")
		retErr = err
	}

	return retErr
}

func (ctl *TaskController) updateOriginalVideoEmbeddings(ctx context.Context, origVideoID, filename string) error {
	videoReader, err := ctl.minioClient.GetFileReader(ctx, origVideoID, ctl.minioClient.GetOrigVideoBucket())
	if err != nil {
		return err
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	if fw, err = w.CreateFormFile("audio_file", filename); err != nil {
		return err
	}

	if _, err = io.Copy(fw, videoReader); err != nil {
		return err
	}

	_, err = ctl.videocopy.UploadUploadVideoPostWithBody(ctx, w.FormDataContentType(), &b)
	if err != nil {
		return fmt.Errorf("update database failed: %w", err)
	}

	return nil
}

func (ctl *TaskController) updateOriginalAudioEmbeddings(ctx context.Context, origVideoID, filename string) error {
	videoReader, err := ctl.minioClient.GetFileReader(ctx, origVideoID, ctl.minioClient.GetOrigVideoBucket())
	if err != nil {
		return err
	}

	tmpfile, err := os.CreateTemp("", "*.mp4")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()

	if _, err = io.Copy(tmpfile, videoReader); err != nil {
		return err
	}

	audioFileName, err := ctl.ffmpegExec.GetAudioFromVideo(tmpfile.Name())
	if err != nil {
		return err
	}

	audioFile, err := os.Open(audioFileName)
	if err != nil {
		return err
	}
	defer os.Remove(audioFileName)
	defer audioFile.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	var fw io.Writer
	if fw, err = w.CreateFormFile("audio_file", filename); err != nil {
		return err
	}

	if _, err = io.Copy(fw, audioFile); err != nil {
		return err
	}

	_, err = ctl.wav2vecClient.UpdateDataBaseUpdateDatabasePostWithBody(ctx, w.FormDataContentType(), &b)
	if err != nil {
		return fmt.Errorf("update database failed: %w", err)
	}

	return nil
}

func (ctl *TaskController) GetTask(ctx context.Context, id int64) (model.Task, error) {
	pgtask, err := ctl.pgConn.GetTask(ctx, id)
	if err != nil {
		return model.Task{}, fmt.Errorf("get task failed: %w", err)
	}

	task, err := taskToModel(pgtask)
	if err != nil {
		return model.Task{}, err
	}

	task.VideoIDUrl, err = ctl.minioClient.GetFileURL(ctx, task.VideoIDUrl, ctl.minioClient.GetVideoBucketName())
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to get video url: %w", err)
	}

	task.PreviewIDUrl, err = ctl.minioClient.GetFileURL(ctx, task.PreviewIDUrl, ctl.minioClient.GetPreviewBucketName())
	if err != nil {
		return model.Task{}, fmt.Errorf("failed to get preview url: %w", err)
	}

	return task, nil
}

func (ctl *TaskController) GetTasks(ctx context.Context, limit, offset uint64) ([]model.Task, int64, error) {
	pgtasks, err := ctl.pgConn.GetTasks(ctx, pgsql.GetTasksParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("get tasks failed: %w", err)
	}

	tasks, err := taskSliceToModel(pgtasks)
	if err != nil {
		return nil, 0, err
	}

	for i := range tasks {
		tasks[i].VideoIDUrl, err = ctl.minioClient.GetFileURL(ctx, tasks[i].VideoIDUrl, ctl.minioClient.GetVideoBucketName())
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get video url: %w", err)
		}

		tasks[i].PreviewIDUrl, err = ctl.minioClient.GetFileURL(ctx, tasks[i].PreviewIDUrl, ctl.minioClient.GetPreviewBucketName())
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get preview url: %w", err)
		}
	}

	total, err := ctl.pgConn.GetTasksCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tasks count: %w", err)
	}

	return tasks, total, nil
}

func (ctl *TaskController) makePreviewUploadVideo(ctx context.Context, file io.Reader) (videoID string, previewID string, videoLen int, err error) {
	id := xid.New().String() + ".mp4"

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err = io.Copy(tmpFile, file); err != nil {
		return "", "", 0, fmt.Errorf("io.Copy failed: %w", err)
	}

	defer func() {
		if errDef := os.Remove(tmpFile.Name()); errDef != nil {
			ctl.log.Error().Err(errDef).Str("path", tmpFile.Name()).Msg("failed to remove tmp file")
		}
	}()

	stat, err := tmpFile.Stat()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get file metainfo: %w", err)
	}

	if _, err = tmpFile.Seek(0, 0); err != nil {
		return "", "", 0, fmt.Errorf("failed to reset reader tmpfile: %w", err)
	}

	if err = ctl.minioClient.UploadFile(ctx, tmpFile, stat.Size(), id, ctl.minioClient.GetVideoBucketName()); err != nil {
		return "", "", 0, fmt.Errorf("failed to upload video to minio: %w", err)
	}

	videoLenFloat, err := ctl.ffmpegExec.GetVideoLength(tmpFile.Name())
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get video length: %w", err)
	}

	previewPicture, err := ctl.ffmpegExec.GetScreenshotFromVideo(tmpFile.Name())
	if err != nil {
		return "", "", 0, fmt.Errorf("make preview failed: %w", err)
	}

	defer func() {
		if errDef := os.Remove(previewPicture); errDef != nil {
			ctl.log.Error().Err(errDef).Str("path", tmpFile.Name()).Msg("failed to remove preview picture file")
		}
	}()

	if err = ctl.minioClient.UploadFileFromOs(ctx, previewPicture, previewPicture, ctl.minioClient.GetPreviewBucketName()); err != nil {
		return "", "", 0, fmt.Errorf("failed to upload picture to minio: %w", err)
	}

	return id, previewPicture, int(videoLenFloat.Seconds()), nil
}

func (ctl *TaskController) uploadOriginalVideoToMinio(ctx context.Context, file io.Reader) (videoID string, err error) {
	id := xid.New().String() + ".mp4"

	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	if _, err = io.Copy(tmpFile, file); err != nil {
		return "", fmt.Errorf("io.Copy failed: %w", err)
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek failed: %w", err)
	}

	defer func() {
		if errK := os.Remove(tmpFile.Name()); errK != nil {
			ctl.log.Error().Err(errK).Str("path", tmpFile.Name()).Msg("failed to remove tmp file")
		}
	}()

	stat, err := tmpFile.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file metainfo: %w", err)
	}

	if err = ctl.minioClient.UploadFile(ctx, tmpFile, stat.Size(), id, ctl.minioClient.GetOrigVideoBucket()); err != nil {
		return "", fmt.Errorf("failed to upload video to minio: %w", err)
	}

	return id, nil
}

func (ctl *TaskController) getHashFromVideo(ctx context.Context, id, bucket string) (string, error) {
	rdr, err := ctl.minioClient.GetFileReader(ctx, id, bucket)
	if err != nil {
		return "", fmt.Errorf("failed to get reader from minio: %w", err)
	}

	h := md5.New()

	if _, err := io.Copy(h, rdr); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
