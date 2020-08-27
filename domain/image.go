package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Image struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	AlbumID  primitive.ObjectID `json:"-" bson:"album_id"`
	Binary   []byte             `bson:"binary"`
	FIleName string             `json:"file_name" bson:"file_name"`
	UploadAt time.Time          `json:"upload_at" bson:"upload_at"`
}

type ImageUsecase interface {
	UploadImage(file []byte, filename string, albumId primitive.ObjectID) (*Image, error)
	GetImageByID(albumId primitive.ObjectID, id string) (*Image, error)
	FindImages(albumId primitive.ObjectID, cursor int) (*[]Image, error)
	RemoveImage(albumId primitive.ObjectID, id string) error
	RemoveImagesByAlbumID(albumId primitive.ObjectID) error
}

type ImageRepository interface {
	Store(s *Image) error
	GetByID(albumId, imgId primitive.ObjectID) (*Image, error)
	Find(albumId primitive.ObjectID, take, skip int) (*[]Image, error)
	Remove(albumId, id primitive.ObjectID) error
	RemoveManyByAlbumId(id primitive.ObjectID) error
	FindByIds(ids []primitive.ObjectID, take, skip int) (*[]Image, error)
}

type ImageSubmission struct {
	Data     []byte
	FileName string
}