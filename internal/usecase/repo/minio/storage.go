package minio

//
//import (
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/SETTER2000/prove/internal/clients"
//	"github.com/SETTER2000/prove/internal/entity"
//	"github.com/SETTER2000/prove/pkg/er"
//	"github.com/SETTER2000/prove/pkg/log/logger"
//	"io"
//)
//
//type minioStorage struct {
//	client *clients.ClientMinio
//	logger logger.Logger
//}
//
//// NewStorage
////TODO ошибка аргумента возвращаемого NewStorage - file.Storage это должен быть интерфейс,
//// который ссылается на repo. Вроде как! 🤨
////func NewStorage(endpoint, accessKeyID, secretAccessKey string, logger logger.Logger) (file.Storage, error) {
////	client, err := clients.NewClientMinio(endpoint, accessKeyID, secretAccessKey, logger)
////	if err != nil {
////		return nil, fmt.Errorf("failed to create minio clients. err: %w", err)
////	}
////	return &minioStorage{
////		client: client,
////	}, nil
////}
//
//func (m *minioStorage) GetFile(ctx context.Context, bucketName, fileID string) (*entity.File, error) {
//	obj, err := m.client.GetFile(ctx, bucketName, fileID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get file. err: %w", err)
//	}
//	defer obj.Close()
//	objectInfo, err := obj.Stat()
//	if err != nil {
//		return nil, fmt.Errorf("failed to get file. err: %w", err)
//	}
//	buffer := make([]byte, objectInfo.Size)
//	_, err = obj.Read(buffer)
//	if err != nil && err != io.EOF {
//		return nil, fmt.Errorf("failed to get objects. err: %w", err)
//	}
//	f := entity.File{
//		ID:    objectInfo.Key,
//		Name:  objectInfo.UserMetadata["Name"],
//		Size:  objectInfo.Size,
//		Bytes: buffer,
//	}
//	return &f, nil
//}
//
//func (m *minioStorage) GetFilesByNoteUUID(ctx context.Context, noteUUID string) ([]*entity.File, error) {
//	objects, err := m.client.GetBucketFiles(ctx, noteUUID)
//	if err != nil {
//		return nil, fmt.Errorf("failed to get objects. err: %w", err)
//	}
//	if len(objects) == 0 {
//		return nil, er.ErrNotFound
//	}
//
//	var files []*entity.File
//	for _, obj := range objects {
//		stat, err := obj.Stat()
//		if err != nil {
//			m.logger.Error("failed to get objects. err: %v", err)
//			continue
//		}
//		buffer := make([]byte, stat.Size)
//		_, err = obj.Read(buffer)
//		if err != nil && err != io.EOF {
//			m.logger.Error("failed to get objects. err: %v", err)
//			continue
//		}
//		f := entity.File{
//			ID:    stat.Key,
//			Name:  stat.UserMetadata["Name"],
//			Size:  stat.Size,
//			Bytes: buffer,
//		}
//		files = append(files, &f)
//		obj.Close()
//	}
//
//	return files, nil
//}
//
//func (m *minioStorage) CreateFile(ctx context.Context, noteUUID string, file *entity.File) error {
//	err := m.client.UploadFile(ctx, file.ID, file.Name, noteUUID, file.Size, bytes.NewBuffer(file.Bytes))
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (m *minioStorage) DeleteFile(ctx context.Context, noteUUID, fileId string) error {
//	err := m.client.DeleteFile(ctx, noteUUID, fileId)
//	if err != nil {
//		return err
//	}
//	return nil
//}
