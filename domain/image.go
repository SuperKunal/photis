package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Image struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Binary   []byte             `bson:"binary"`
	FIleName string             `json:"file_name" bson:"file_name"`
	UploadAt time.Time          `json:"upload_at" bson:"upload_at"`
}

type ImageUsecase interface {
	UploadImage(file []byte, filename string) (*Image, error)
	GetImageByID(id string) (*Image, error)
	FindImages(ids []primitive.ObjectID, cursor int) (*[]Image, error)
	RemoveImage(id string) (*primitive.ObjectID, error)
	RemoveImages(ids []primitive.ObjectID) error
}

type ImageRepository interface {
	Store(s *Image) error
	GetByID(id primitive.ObjectID) (*Image, error)
	Find(take, skip int) (*[]Image, error)
	Remove(id primitive.ObjectID) error
	RemoveMany(ids []primitive.ObjectID) error
	FindByIds(ids []primitive.ObjectID, take, skip int) (*[]Image, error)
}

type ImageSubmission struct {
	Data     []byte
	FileName string
}