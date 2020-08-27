package main

import (
	"context"
	"github.com/savsgio/atreugo/v11"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	_albumRepository "photis/album/repository/mongo"
	"photis/config"
	_imgRepository "photis/image/repository/mongo"
	"photis/middleware"
	"photis/services"

	_handler "photis/album/delivery/http"
	_albumUsecase "photis/album/usecase"
	_imgUsecase "photis/image/usecase"
)

func main() {
	_config := config.NewConfig()
	_middleware := middleware.InitMiddleware()

	databaseClient := InitiateMongoClient(_config.ConnectionString).Database(_config.Database)
	rabbitMqClient := services.NewRabbitMqClient(_config.RabbitMqAddr, _config.QueueName)

	imgRepo := _imgRepository.NewMongoImageRepository(databaseClient)
	albumRepo := _albumRepository.NewMongoAlbumRepository(databaseClient)

	imgUsecase := _imgUsecase.NewImageUsecase(imgRepo)
	albumUsecase := _albumUsecase.NewAlbumUsecase(albumRepo, imgUsecase, rabbitMqClient)

	server := atreugo.New(atreugo.Config{Addr: "0.0.0.0:7000", PanicView: _middleware.PanicHandler})
	api := server.NewGroupPath("/api")

	_handler.NewImageHandler(api, albumUsecase)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func InitiateMongoClient(connectionString string) *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}

	return client
}