package handlers

import (
	"encoding/json"
	"net/http"
	"pr9/pkg/file"
	"pr9/pkg/helpers"
	"pr9/pkg/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type FileHandler struct {
	FileService *services.FileService
	Logger *zap.SugaredLogger
	UploadedFile string
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(200 * 1024 * 1024)
	file, header, err := r.FormFile(h.UploadedFile)
	if err != nil {
		helpers.JSONMessageSend(w, err.Error(), http.StatusBadRequest)
		h.Logger.Infof("File upload error: %w", err)
		return
	}
	defer file.Close()

	id, err := h.FileService.UploadFile(file, header.Filename)
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("File upload databases error: %w", err)
		return
	}

	helpers.JSONSend(w, map[string]string{"file ID": id}, http.StatusCreated)
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["FILE_ID"]
	if !ok {
		helpers.JSONMessageSend(w, "bad file id", http.StatusBadRequest)
		h.Logger.Info("File download error")
		return
	}

	fileData, fileName, err := h.FileService.DownloadFile(id)
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("Files download databases error: %w", err)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.WriteHeader(http.StatusOK)
	w.Write(fileData)
}

func (h *FileHandler) DownloadFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.FileService.DownloadFiles()
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("Files download databases error: %w", err)
		return
	}
	helpers.JSONSend(w, map[string][]file.File{"files": files}, http.StatusOK)
}

func (h *FileHandler) DownloadFileInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["FILE_ID"]
	if !ok {
		helpers.JSONMessageSend(w, "bad file id", http.StatusBadRequest)
		h.Logger.Info("File info download error")
		return
	}
	file, err := h.FileService.DownloadFileInfo(id)
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("File info download database error: %w", err)
		return
	}
	helpers.JSONSend(w, file, http.StatusOK)
}

func (h *FileHandler) RenameFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["FILE_ID"]
	if !ok {
		helpers.JSONMessageSend(w, "bad file id", http.StatusBadRequest)
		h.Logger.Info("File info download error")
		return
	}
	var newFileName struct {
		FileName string `json:"file_name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newFileName); err != nil {
		http.Error(w, "Error decoding new file name", http.StatusBadRequest)
		h.Logger.Info("Error decoding new file name")
		return
	}
	
	err := h.FileService.RenameFile(id, newFileName.FileName)
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("File rename error: %w", err)
		return
	}
	helpers.JSONMessageSend(w, "success", http.StatusOK)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["FILE_ID"]
	if !ok {
		helpers.JSONMessageSend(w, "bad file id", http.StatusBadRequest)
		h.Logger.Info("File info download error")
		return
	}
	err := h.FileService.DeleteFile(id)
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("File delete error: %w", err)
		return
	}
	helpers.JSONMessageSend(w, "success", http.StatusOK)
}

func (h *FileHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(200 * 1024 * 1024)
	file, header, err := r.FormFile(h.UploadedFile)
	if err != nil {
		helpers.JSONMessageSend(w, err.Error(), http.StatusBadRequest)
		h.Logger.Infof("File update error: %w", err)
		return
	}
	defer file.Close()

	vars := mux.Vars(r)
	id, ok := vars["FILE_ID"]
	if !ok {
		helpers.JSONMessageSend(w, "bad file id", http.StatusBadRequest)
		h.Logger.Info("File update id error")
		return
	}

	err = h.FileService.UpdateFile(file, id, header.Filename)
	if err != nil {
		helpers.DatabaseError(w, err, h.Logger)
		h.Logger.Errorf("File update database error: %w", err)
		return
	}

	helpers.JSONMessageSend(w, "success", http.StatusOK)
}