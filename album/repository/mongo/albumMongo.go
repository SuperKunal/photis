package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"photis/domain"
)

const AlbumCollectionName string = "albums"

type mongoAlbumRepository struct {
	collection *mongo.Collection
}

func (m *mongoAlbumRepository) Store(album *domain.Album) error {
	_, err := m.collection.InsertOne(context.Background(), album)
	if err != nil {
		return err
	}
	return nil
}

func (m *mongoAlbumRepository) GetByID(id primitive.ObjectID) (*domain.Album, error) {
	var album domain.Album
	result := m.collection.FindOne(context.Background(), bson.M{"_id" : id})
	if err := result.Decode(&album); err != nil {
		return nil, nil
	}
	return &album, nil
}

func (m *mongoAlbumRepository) Find(take, skip int) (*[]domain.Album, error) {
	var albums []domain.Album

	findOptions := options.FindOptions{}
	limit := int64(take)
	s := int64(skip)
	findOptions.Limit = &limit
	findOptions.Skip = &s
	findOptions.SetSort(bson.D{{"_id", -1}})
	findOptions.SetProjection(bson.D{{"images", 0}})
	cur, err := m.collection.Find(context.TODO(), bson.M{}, &findOptions)

	if err != nil {
		return nil, nil
	}

	for cur.Next(context.TODO()) {
		var elem domain.Album
		err := cur.Decode(&elem)
		if err != nil {
			break
		}
		albums = append(albums, elem)
	}

	return &albums, nil
}

func (m *mongoAlbumRepository) Exist(id primitive.ObjectID) bool {
	count, err := m.collection.CountDocuments(context.Background(), bson.M{"_id" : id})
	if err != nil {
		return false
	}
	return count != 0
}

func (m *mongoAlbumRepository) Remove(id primitive.ObjectID) error {
	if _, err := m.collection.DeleteOne(context.Background(), bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}

func (m *mongoAlbumRepository) AddImageToAlbum(albumId, imageId primitive.ObjectID) error {
	filter := bson.M{"_id": albumId}
	update := bson.M{"$push": bson.M{"images": imageId}}
	if err := m.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}
	return nil
}

func (m *mongoAlbumRepository) RemoveImageFromAlbum(albumId, imageId primitive.ObjectID) error {
	filter := bson.M{"_id": albumId}
	update := bson.M{"$pull": bson.M{"images": imageId}}
	if err := m.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}
	return nil
}

func NewMongoAlbumRepository(Conn *mongo.Database) domain.AlbumRepository {
	return &mongoAlbumRepository{Conn.Collection(AlbumCollectionName)}
}