package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/gulldan/lct2024_copyright/bff/pkg/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	taskcontroller "github.com/gulldan/lct2024_copyright/bff/internal/controller/task_controller"
	bffv1 "github.com/gulldan/lct2024_copyright/bff/proto/bff/v1"
)

func New(logger *zerolog.Logger, grpcEndpoint string) (http.Handler, error) {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := bffv1.RegisterScanTasksServiceHandlerFromEndpoint(context.Background(), mux, grpcEndpoint, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to register tasks handler: %w", err)
	}

	return mux, nil
}

func AddCors(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return true
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Location"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(h)

	return handler
}

func AddSwagger(f fs.FS, h http.Handler) (http.Handler, error) {
	mux, ok := h.(*runtime.ServeMux)
	if !ok {
		return nil, errors.New("got not a serveMux")
	}

	fsys, err := fs.Sub(f, "swagger-ui/docs")
	if err != nil {
		return nil, fmt.Errorf("fs sub failed: %w", err)
	}

	fileServer := http.FileServer(http.FS(fsys))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			mux.ServeHTTP(w, r)
			return
		}

		fileServer.ServeHTTP(w, r)
	}), nil
}

func HandleBinaryFileUpload(h http.Handler, cfg *config.Config, logger *zerolog.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v1/tasks/create/upload") {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse form: %s", err.Error()), http.StatusBadRequest)
				return
			}

			f, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to get file 'attachment': %s", err.Error()), http.StatusBadRequest)
				return
			}
			defer f.Close()

			taskCtl := taskcontroller.New(nil, nil)

			id, err := taskCtl.CreateTask(context.Background(), f, header.Filename)
			if err != nil {
				logger.Error().Err(err).Msg("create task failed")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			respBody, err := json.Marshal(bffv1.CreateTaskFromFileResponse{
				Id: id,
			})
			if err != nil {
				logger.Error().Err(err).Msg("json marshal failed")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(respBody)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

func HandleOriginalFileUpload(h http.Handler, cfg *config.Config, logger *zerolog.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/v1/original/upload") {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to parse form: %s", err.Error()), http.StatusBadRequest)
				return
			}

			f, header, err := r.FormFile("file")
			if err != nil {
				logger.Error().Err(err).Msg("orig formfile failed")
				http.Error(w, fmt.Sprintf("failed to get file 'attachment': %s", err.Error()), http.StatusBadRequest)
				return
			}
			defer f.Close()

			enableStr := r.FormValue("not_upload_embeddings")
			notUploadEmbeddings, err := strconv.ParseBool(enableStr)
			if err != nil {
				notUploadEmbeddings = false
			}

			taskCtl := taskcontroller.New(nil, nil)

			_, err = taskCtl.UploadOriginalVideo(context.Background(), f, header.Filename, !notUploadEmbeddings)
			if err != nil {
				logger.Error().Err(err).Msg("upload original video failed")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
