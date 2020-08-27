# photis
golang based album/image REST APIs that follows the clean architecture. 

#### To Run

- MongoDB Self hosted service should be running on local
- RabbitMQ service should be running on local
- go mod download
- go run main.go

### Features: 

- Images are being stored to Mongo as in binary format
- Code Structure Follows the Clean Architecture


### Improvements:

- JSON Error Codes
- S3 to store images 
- Unit Testing

### API Examples: 

#### Create Album
```http request
curl --location --request POST 'localhost:7000/api/album/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "album test"
}'
```

#### Get Albums
```http request
curl --location --request GET 'localhost:7000/api/album/'
```

#### Delete Album by ID
```http request
curl --location --request DELETE 'localhost:7000/api/album/5f478562757cbfec5f4b5c20'
```

#### Add an Image to Album
```http request
curl --location --request POST 'localhost:7000/api/album/5f47415d2f855088b50daeba/image' \
--form 'data=@/D:/path/to/file/13535985.jpg' \
--form 'filename=abc.jpg'
```

#### Get an Image by ID
```http request
curl --location --request GET 'localhost:7000/api/album/5f47415d2f855088b50daeba/image/5f4786a2757cbfec5f4b5c21'
```

#### Get all Images in Album
```http request
curl --location --request GET 'localhost:7000/api/album/5f47415d2f855088b50daeba/image'
```

#### Remove an Image from Album
```http request
curl --location --request DELETE 'localhost:7000/api/album/5f47415d2f855088b50daeba/image/5f4786a2757cbfec5f4b5c21'
```