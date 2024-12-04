package services

import (
	"io"
	"pr9/pkg/file"
	"pr9/pkg/repositories"
)

type FileService struct{
	FileRepo *repositories.FileRepo
}

func (s *FileService) UploadFile(file io.Reader, fileName string) (string, error) {
	return s.FileRepo.UploadFile(file, fileName)
}

func (s *FileService) DownloadFiles() ([]file.File, error) {
	return s.FileRepo.DownloadFiles()
}

func (s *FileService) DownloadFile(id string) ([]byte, string, error) {
	return s.FileRepo.DownloadFile(id)
}

func (s *FileService) DownloadFileInfo(id string) (file.File, error) {
	return s.FileRepo.DownloadFileInfo(id)
}

func (s *FileService) UpdateFile(file io.Reader, id, fileName string) (error) {
	return s.FileRepo.UpdateFile(file, id, fileName)
}

func (s *FileService) RenameFile(id, fileName string) (error) {
	return s.FileRepo.RenameFile(id, fileName)
}

func (s *FileService) DeleteFile(id string) (error) {
	return s.FileRepo.DeleteFile(id)
}