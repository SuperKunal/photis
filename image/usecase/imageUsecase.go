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

func (i *imageUsecase) UploadImage(file []byte, filename string) (*domain.Image, error) {

	image := &domain.Image{
		ID:       primitive.NewObjectID(),
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

func (i *imageUsecase) GetImageByID(id string) (*domain.Image, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid image id")
	}

	doc , err := i.imageRepo.GetByID(objectId)

	if err != nil{
		return nil, errors.New("not found")
	}

	return doc, nil
}

func (i *imageUsecase) FindImages(ids []primitive.ObjectID, cursor int) (*[]domain.Image, error) {
	const pageSize = 10

	take := pageSize
	skip := 0

	if cursor > 0 {
		skip = pageSize * cursor
	}

	docs, err := i.imageRepo.Find(take, skip)

	if err != nil {
		return nil, errors.New("not found")
	}

	return docs, nil
}

func (i *imageUsecase) RemoveImage(id string) (*primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid image id")
	}

	if err = i.imageRepo.Remove(objectId); err != nil {
		return nil, err
	}
	return &objectId, nil
}

func (i *imageUsecase) RemoveImages(ids []primitive.ObjectID) error {
	return i.imageRepo.RemoveMany(ids)
}

func NewImageUsecase(r domain.ImageRepository) domain.ImageUsecase {
	return &imageUsecase{imageRepo: r}
}