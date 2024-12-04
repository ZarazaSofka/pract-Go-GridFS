package repositories

import (
	"context"
	"fmt"
	"io"
	"os"
	"pr9/pkg/file"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FileRepo struct {
	Files *gridfs.Bucket
}

func GetNewFileRepo() *FileRepo {
	ctx := context.Background()
	sess, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db"))
	if err != nil {
		panic(err)
	}

	dbName, ok := os.LookupEnv("DB_NAME")
	if !ok {
		dbName = "myDB"
	}

	db := sess.Database(dbName)
	files, err := gridfs.NewBucket(db)
	if err != nil {
	   panic(err)
	}

	return &FileRepo{
		Files: files,
	}
}

func (repo *FileRepo) UploadFile(file io.Reader, fileName string) (string, error) {
	objectID, err := repo.Files.UploadFromStream(fileName, file)
	if err != nil {
		return "", nil
	}
	return objectID.Hex(), nil
}

func (repo *FileRepo) UpdateFile(file io.Reader, id, fileName string) error {
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = repo.Files.Delete(fileID)
	if err != nil {
		return err
	}
	return repo.Files.UploadFromStreamWithID(fileID, fileName, file)
}

func (repo *FileRepo) DownloadFiles() ([]file.File, error) {
	cursor, err := repo.Files.Find(bson.D{})
	if err != nil {
		return nil, err
	}
	var files []file.File
	err = cursor.All(context.Background(), &files)
	return files, err
}

func (repo *FileRepo) DownloadFile(id string) ([]byte, string, error) {
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, "", err
	}
	file, err := repo.Files.OpenDownloadStream(fileID)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	fileName := file.GetFile().Name
	data := make([]byte, file.GetFile().Length)
	_, err = file.Read(data)
	return data, fileName, err
}

func (repo *FileRepo) DownloadFileInfo(id string) (file.File, error) {
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return file.File{}, err
	}
	cursor, err := repo.Files.Find(bson.D{{Key: "_id", Value: fileID}})
	if err != nil {
		return file.File{}, err
	}
	defer cursor.Close(context.Background())
	
	if cursor.Next(context.Background()) {
        var file file.File
        err = cursor.Decode(&file)
        return file, err
    }

    return file.File{}, fmt.Errorf("document not found with ID: %s", id)
}

func (repo *FileRepo) RenameFile(id, fileName string) error {
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return repo.Files.Rename(fileID, fileName)
}

func (repo *FileRepo) DeleteFile(id string) error {
	fileID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return repo.Files.Delete(fileID)
}