package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"photis/domain"
)

const ImageCollectionName string = "images"

type mongoImageRepository struct {
	collection *mongo.Collection
}

func (m *mongoImageRepository) Store(image *domain.Image) error {
	_, err := m.collection.InsertOne(context.Background(), image)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoImageRepository) GetByID(albumId, imgId primitive.ObjectID) (*domain.Image, error) {
	var image domain.Image
	result := m.collection.FindOne(context.Background(), bson.M{"_id" : imgId, "album_id": albumId})
	if err := result.Decode(&image); err != nil {
		return nil, nil
	}
	return &image, nil
}

func (m *mongoImageRepository) Find(albumId primitive.ObjectID, take, skip int) (*[]domain.Image, error) {
	var images []domain.Image

	findOptions := options.FindOptions{}
	limit := int64(take)
	s := int64(skip)
	findOptions.Limit = &limit
	findOptions.Skip = &s
	findOptions.SetSort(bson.D{{"_id", -1}})
	cur, err := m.collection.Find(context.TODO(), bson.M{"album_id": albumId}, &findOptions)

	if err != nil {
		return nil, nil
	}

	for cur.Next(context.TODO()) {
		var elem domain.Image
		err := cur.Decode(&elem)
		if err != nil {
			break
		}
		images = append(images, elem)
	}

	return &images, nil
}

func (m *mongoImageRepository) Remove(albumId, id primitive.ObjectID) error {
	filter := bson.M{"_id": id, "album_id": albumId}
	if _, err := m.collection.DeleteOne(context.Background(), filter); err != nil{
		return err
	}
	return nil
}

func (m *mongoImageRepository) RemoveManyByAlbumId(albumId primitive.ObjectID) error {
	filter := bson.M{"album_id": albumId}
	_, err := m.collection.DeleteMany(context.Background(), filter)
	return err
}

func (m *mongoImageRepository) FindByIds(ids []primitive.ObjectID, take, skip int) (*[]domain.Image, error) {
	var images []domain.Image

	findOptions := options.FindOptions{}
	limit := int64(take)
	s := int64(skip)
	findOptions.Limit = &limit
	findOptions.Skip = &s
	findOptions.SetSort(bson.D{{"_id", -1}})
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cur, err := m.collection.Find(context.TODO(), filter, &findOptions)

	if err != nil {
		return nil, nil
	}

	for cur.Next(context.TODO()) {
		var elem domain.Image
		err := cur.Decode(&elem)
		if err != nil {
			break
		}
		images = append(images, elem)
	}

	return &images, nil
}

func NewMongoImageRepository(Conn *mongo.Database) domain.ImageRepository {
	return &mongoImageRepository{Conn.Collection(ImageCollectionName)}
}
