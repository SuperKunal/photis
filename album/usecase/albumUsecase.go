package usecase

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"photis/domain"
	"photis/services"
	"time"
)

type albumUsecase struct {
	albumRepo      domain.AlbumRepository
	imageUsecase   domain.ImageUsecase
	rabbitMqClient *services.Client
}

func (uc *albumUsecase) CreateAlbum(name string) (*domain.Album, error) {
	album := &domain.Album{
		ID:        primitive.NewObjectID(),
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}

	err := uc.albumRepo.Store(album)
	if err != nil {
		return nil, errors.New("error occurred")
	}
	return album, nil
}

func (uc *albumUsecase) RemoveAlbum(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid album id")
	}

	album, err := uc.albumRepo.GetByID(objectId)
	if err != nil || album == nil {
		return errors.New("invalid album id")
	}

	if err = uc.imageUsecase.RemoveImagesByAlbumID(album.ID); err != nil{
		return err
	}

	if err := uc.albumRepo.Remove(objectId); err != nil {
		return errors.New("error occurred")
	}
	return nil
}

func (uc *albumUsecase) FindAlbums(cursor int) (*[]domain.Album, error) {
	const pageSize = 10

	take := pageSize
	skip := 0

	if cursor > 0 {
		skip = pageSize * cursor
	}

	docs, err := uc.albumRepo.Find(take, skip)

	if err != nil {
		return nil, errors.New("not found")
	}

	return docs, nil
}

func (uc *albumUsecase) AddImage(albumId string, file []byte, fileName string) (*domain.Image, error) {
	objectId, err := primitive.ObjectIDFromHex(albumId)
	if err != nil {
		return nil, errors.New("invalid album id")
	}

	if !uc.albumRepo.Exist(objectId) {
		return nil, errors.New("invalid album id")
	}

	image, err := uc.imageUsecase.UploadImage(file, fileName, objectId)
	if err != nil {
		return nil, errors.New("error occurred")
	}

	go uc.rabbitMqClient.Publish("Image: " + fileName + " received!")

	return image, nil
}

func (uc *albumUsecase) RemoveImageById(albumId, imageId string) error {
	objectId, err := primitive.ObjectIDFromHex(albumId)
	if err != nil {
		return errors.New("invalid album id")
	}

	if !uc.albumRepo.Exist(objectId) {
		return errors.New("invalid album id")
	}

	err = uc.imageUsecase.RemoveImage(objectId, imageId)
	if err != nil {
		return errors.New("error occurred")
	}

	go uc.rabbitMqClient.Publish("Image deleted!")

	return nil
}

func (uc *albumUsecase) FindImages(albumId string, cursor int) (*[]domain.Image, error) {
	objectId, err := primitive.ObjectIDFromHex(albumId)
	if err != nil {
		return nil, errors.New("invalid album id")
	}

	album, err := uc.albumRepo.GetByID(objectId)
	if err != nil || album == nil {
		return nil, errors.New("invalid album id")
	}

	images, err := uc.imageUsecase.FindImages(album.ID, cursor)
	if err != nil {
		return nil, errors.New("images not found")
	}

	return images, nil
}

func (uc *albumUsecase) FindImageById(albumId, imageId string) (*domain.Image, error) {
	albumObjectId, err := primitive.ObjectIDFromHex(albumId)
	if err != nil {
		return nil, errors.New("invalid album id")
	}

	if !uc.albumRepo.Exist(albumObjectId) {
		return nil, errors.New("invalid album id")
	}

	image, err := uc.imageUsecase.GetImageByID(albumObjectId, imageId)
	if err != nil {
		return nil, errors.New("image not found")
	}

	return image, nil
}

func NewAlbumUsecase(r domain.AlbumRepository, uc domain.ImageUsecase, client *services.Client) domain.AlbumUsecase {
	return &albumUsecase{albumRepo: r, imageUsecase: uc, rabbitMqClient: client}
}
