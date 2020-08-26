package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Album struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id"`
	Name      string               `bson:"name" json:"name"`
	Images    []primitive.ObjectID `json:"images" bson:"images"`
	CreatedAt time.Time            `json:"created_at" bson:"created_at"`
}

type AlbumUsecase interface {
	CreateAlbum(name string) (*Album, error)
	RemoveAlbum(id string) error
	FindAlbums(cursor int) (*[]Album, error)
	AddImage(albumId string, file []byte, fileName string) (*Image, error)
	FindImages(albumId string, cursor int) (*[]Image, error)
	FindImageById(albumId, imageId string) (*Image, error)
	RemoveImageById(albumId, imageId string) error
}

type AlbumRepository interface {
	Store(album *Album) error
	GetByID(id primitive.ObjectID) (*Album, error)
	Find(take, skip int) (*[]Album, error)
	Exist(id primitive.ObjectID) bool
	Remove(id primitive.ObjectID) error

	AddImageToAlbum(albumId, imageId primitive.ObjectID) (error)
	RemoveImageFromAlbum(albumId, imageId primitive.ObjectID) (error)
}

// Models

type CreateAlbumRequest struct {
	Name string `json:"name"`
}