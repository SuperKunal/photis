package http

import (
	"encoding/json"
	"errors"
	"github.com/savsgio/atreugo/v11"
	"net/http"
	"photis/domain"
	"strconv"
	"strings"
)

type ImageHandler struct {
	AlbumUsecase domain.AlbumUsecase
}

func NewImageHandler(router *atreugo.Router, auc domain.AlbumUsecase) {
	handler := &ImageHandler{AlbumUsecase: auc}
	router.POST("/album", handler.CreateAlbum)
	router.GET("/album", handler.FindAlbums)
	router.DELETE("/album/{id}", handler.DeleteAlbum)

	router.POST("/album/{id}/image", handler.Upload)
	router.GET("/album/{id}/image", handler.FindImages)
	router.GET("/album/{id}/image/{imgId}", handler.GetImageById)
	router.DELETE("/album/{id}/image/{imgId}", handler.DeleteImage)
}

func (handler *ImageHandler) CreateAlbum(ctx *atreugo.RequestCtx) error {
	var payload domain.CreateAlbumRequest
	if err := json.Unmarshal(ctx.PostBody(), &payload); err != nil || len(payload.Name) == 0 {
		return ctx.JSONResponse(map[string]string{"error": "name is required"}, http.StatusBadRequest)
	}

	album, err := handler.AlbumUsecase.CreateAlbum(payload.Name)
	if err != nil {
		return ctx.JSONResponse(map[string]string{"error": err.Error()}, http.StatusInternalServerError)
	}

	return ctx.JSONResponse(album, http.StatusOK)
}

func (handler *ImageHandler) FindAlbums(ctx *atreugo.RequestCtx) error {
	cursor, _ := strconv.Atoi(string(ctx.QueryArgs().Peek("cursor")))

	albums, err := handler.AlbumUsecase.FindAlbums(cursor)
	if err != nil {
		return ctx.JSONResponse(map[string]string{"error": err.Error()}, http.StatusInternalServerError)
	}

	return ctx.JSONResponse(albums, http.StatusOK)
}

func (handler *ImageHandler) DeleteAlbum(ctx *atreugo.RequestCtx) error {
	id := ctx.UserValue("id").(string)
	if err := handler.AlbumUsecase.RemoveAlbum(id); err != nil{
		return ctx.JSONResponse(map[string]string{"error": err.Error()}, http.StatusInternalServerError)
	}

	return ctx.JSONResponse(map[string]string{"message": "success"}, http.StatusOK)
}

func (handler *ImageHandler) Upload(context *atreugo.RequestCtx) error {
	id := context.UserValue("id").(string)
	if len(id) == 0 {
		return context.JSONResponse(map[string]string{"error": "invalid album id"}, http.StatusBadRequest)
	}

	contentType := string(context.Request.Header.ContentType())

	if strings.HasPrefix(contentType, "multipart/form-data") {
		imageData, err, errCode := parseMultipartImage(context)
		if err != nil {
			return context.ErrorResponse(err, errCode)
		}

		img, err := handler.AlbumUsecase.AddImage(id, imageData.Data, imageData.FileName)
		if err != nil {
			return context.JSONResponse(map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		}

		return context.JSONResponse(img, http.StatusOK)
	}

	return context.JSONResponse(map[string]string{"error": "bad request"}, http.StatusBadRequest)
}

func (handler *ImageHandler) GetImageById(context *atreugo.RequestCtx) error {
	id := context.UserValue("id").(string)
	imgId := context.UserValue("imgId").(string)
	doc, err := handler.AlbumUsecase.FindImageById(id, imgId)

	if err != nil {
		return context.JSONResponse(map[string]string{"error": err.Error()}, http.StatusBadRequest)
	}

	return context.JSONResponse(doc, http.StatusOK)
}

func (handler *ImageHandler) FindImages(context *atreugo.RequestCtx) error {
	albumId := context.UserValue("id").(string)
	if len(albumId) == 0 {
		return context.JSONResponse(map[string]string{"error": "invalid album id"}, http.StatusBadRequest)
	}

	cursor, _ := strconv.Atoi(string(context.QueryArgs().Peek("cursor")))
	docs, err := handler.AlbumUsecase.FindImages(albumId, cursor)

	if err != nil {
		return context.JSONResponse(map[string]string{"error": err.Error()}, http.StatusBadRequest)
	}

	return context.JSONResponse(docs, http.StatusOK)
}

func (handler *ImageHandler) DeleteImage(context *atreugo.RequestCtx) error {
	id := context.UserValue("id").(string)
	imgId := context.UserValue("imgId").(string)

	err := handler.AlbumUsecase.RemoveImageById(id, imgId)

	if err != nil {
		return context.JSONResponse(map[string]string{"error": err.Error()}, http.StatusBadRequest)
	}

	return context.JSONResponse(map[string]string{"message": "success"}, http.StatusOK)
}

func parseMultipartImage(ctx *atreugo.RequestCtx) (*domain.ImageSubmission, error, int) {
	form, err := ctx.FormFile("data")
	if err != nil {
		return nil, errors.New("invalid file"), 400
	}

	contentLength := form.Size
	if contentLength == 0 {
		return nil, errors.New("invalid content-length"), 400
	}

	open, err := form.Open()
	if err != nil {
		return nil, errors.New("invalid content"), 400
	}

	buffer := make([]byte, contentLength)
	_, _ = open.Read(buffer)

	filename := string(ctx.FormValue("filename"))

	return &domain.ImageSubmission{
		Data: buffer,
		FileName: filename,
	}, nil, 200
}