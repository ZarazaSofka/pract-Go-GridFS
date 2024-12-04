package main

import (
	"net/http"
	"os"
	"pr9/pkg/handlers"
	"pr9/pkg/repositories"
	"pr9/pkg/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	zapCfg := zap.NewProductionConfig()
	zapCfg.OutputPaths = []string{
		"/var/log/pr9/pr9.log",
		"stdout",
	}
	zapLogger, err := zapCfg.Build()
	if err != nil {
		panic("Logger creating error")
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	uploadedFile, ok := os.LookupEnv("FILE_NAME")
	if !ok {
		uploadedFile = "uploaded_file"
	}

	fh := handlers.FileHandler{
		FileService: &services.FileService{
			FileRepo: repositories.GetNewFileRepo(),
		},
		Logger: logger,
		UploadedFile: uploadedFile,
	}
	r := mux.NewRouter()
	r.HandleFunc("/files", fh.UploadFile).Methods(http.MethodPost)
	r.HandleFunc("/files", fh.DownloadFiles).Methods(http.MethodGet)
	r.HandleFunc("/files/{FILE_ID}", fh.DownloadFile).Methods(http.MethodGet)
	r.HandleFunc("/files/{FILE_ID}/info", fh.DownloadFileInfo).Methods(http.MethodGet)
	r.HandleFunc("/files/{FILE_ID}", fh.RenameFile).Methods(http.MethodPatch)
	r.HandleFunc("/files/{FILE_ID}", fh.UpdateFile).Methods(http.MethodPut)
	r.HandleFunc("/files/{FILE_ID}", fh.DeleteFile).Methods(http.MethodDelete)


	addr, ok := os.LookupEnv("PORT")
	if !ok {
		addr = ":8080"
	}

	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		panic(err)
	}
}
