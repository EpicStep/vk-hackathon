package image

import (
	"bytes"
	"github.com/EpicStep/vk-hackathon/internal/image/model"
	"github.com/EpicStep/vk-hackathon/internal/jsonutil"
	v1 "github.com/EpicStep/vk-hackathon/pkg/api/v1"
	"github.com/corona10/goimagehash"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"net/http"
)

const (
	MaxSize = 10 << 20
)

func (s *Service) Upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := r.ParseMultipartForm(MaxSize)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(3, "Cannot parse file"))
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(3, "Cannot parse file form"))
		return
	}

	defer file.Close()

	buffer := make([]byte, header.Size)
	_, err = file.Read(buffer)
	if err != nil {
		return
	}

	imgj, err := jpeg.Decode(bytes.NewReader(buffer))
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "Image not in JPEG"))
		return
	}

	hash, err := goimagehash.DifferenceHash(imgj)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(7, "Failed to calculate hash"))
		return
	}

	images, err := s.db.GetAllImagesHash(ctx)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(6, "Failed to get all images"))
		return
	}

	imgcfg, _, err := image.DecodeConfig(bytes.NewReader(buffer))
	if err != nil {
		return
	}

	for _, v := range images {
		h := goimagehash.NewImageHash(v.Hash, goimagehash.DHash)
		distance, err := hash.Distance(h)
		if err != nil {
			return
		}

		if distance < 10 {
			if imgcfg.Width > v.Width && imgcfg.Height > v.Height {
				img := model.Image{
					ID:     v.ID,
					Image:  buffer,
					Hash:   hash.GetHash(),
					Height: imgcfg.Height,
					Width: imgcfg.Width,
				}

				err := s.db.UpdateByID(ctx, &img)
				if err != nil {
					jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(6, "Failed to update image"))
					return
				}
			}

			jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(v1.UploadResponse{ID: v.ID}))
			return
		}
	}

	img := model.Image{
		ID:    uuid.New().String(),
		Image: buffer,
		Hash: hash.GetHash(),
		Height: imgcfg.Height,
		Width: imgcfg.Width,
	}

	err = s.db.Create(ctx, &img)
	if err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(4, "Failed to add image to database"))
		return
	}

	jsonutil.MarshalResponse(w, http.StatusOK, jsonutil.NewSuccessfulResponse(v1.UploadResponse{ID: img.ID}))
}
