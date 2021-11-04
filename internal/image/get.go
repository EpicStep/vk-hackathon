package image

import (
	"bytes"
	"database/sql"
	"errors"
	"github.com/EpicStep/vk-hackathon/internal/jsonutil"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"net/http"
	"strconv"
)

func (s *Service) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := r.URL.Query()

	id := params.Get("id")
	scale, err := strconv.ParseFloat(params.Get("scale"), 64)
	if err != nil || scale <= 0 {
		scale = 1
	}

	img, err := s.db.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(2, "Image not found"))
		} else {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(5, "Failed to get image"))
		}

		return
	}

	var resized image.Image
	if scale != 1 {
		imgj, err := jpeg.Decode(bytes.NewReader(img.Image))
		if err != nil {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(3, "Image not in JPEG"))
			return
		}

		resized = resize.Resize(uint(float64(img.Width)*scale), uint(float64(img.Height)*scale), imgj, resize.Lanczos3)
	}

	w.Header().Set("Content-Type", "image/jpeg")

	if scale != 1 {
		buf := new(bytes.Buffer)
		_ = jpeg.Encode(buf, resized, nil)

		w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
		if _, err := w.Write(buf.Bytes()); err != nil {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(6, "Failed to return image"))
			return
		}

		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(img.Image)))
	if _, err := w.Write(img.Image); err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(6, "Failed to return image"))
		return
	}
}
