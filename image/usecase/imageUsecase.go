package usecase

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"photis/domain"
	"time"
)

type imageUsecase struct {
	imageRepo     domain.ImageRepository
}

func (i *imageUsecase) UploadImage(file []byte, filename string, albumId primitive.ObjectID) (*domain.Image, error) {

	image := &domain.Image{
		ID:       primitive.NewObjectID(),
		AlbumID:  albumId,
		Binary:   file,
		FIleName: filename,
		UploadAt: time.Now().UTC(),
	}

	err := i.imageRepo.Store(image)
	if err != nil {
		return nil, errors.New("error occured")
	}
	return image, nil
}

func (i *imageUsecase) GetImageByID(albumId primitive.ObjectID, id string) (*domain.Image, error) {
	imgId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid image id")
	}

	doc , err := i.imageRepo.GetByID(albumId, imgId)

	if err != nil{
		return nil, errors.New("not found")
	}

	return doc, nil
}

func (i *imageUsecase) FindImages(albumId primitive.ObjectID, cursor int) (*[]domain.Image, error) {
	const pageSize = 10

	take := pageSize
	skip := 0

	if cursor > 0 {
		skip = pageSize * cursor
	}

	docs, err := i.imageRepo.Find(albumId, take, skip)

	if err != nil {
		return nil, errors.New("not found")
	}

	return docs, nil
}

func (i *imageUsecase) RemoveImage(albumId primitive.ObjectID, imgId string) error {
	objectId, err := primitive.ObjectIDFromHex(imgId)
	if err != nil {
		return errors.New("invalid image id")
	}

	if err = i.imageRepo.Remove(albumId, objectId); err != nil {
		return err
	}
	return nil
}

func (i *imageUsecase) RemoveImagesByAlbumID(albumId primitive.ObjectID) error {
	return i.imageRepo.RemoveManyByAlbumId(albumId)
}

func NewImageUsecase(r domain.ImageRepository) domain.ImageUsecase {
	return &imageUsecase{imageRepo: r}
}