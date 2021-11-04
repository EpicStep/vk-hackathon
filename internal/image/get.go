package image

import (
	"database/sql"
	"errors"
	"github.com/EpicStep/vk-hackathon/internal/jsonutil"
	"net/http"
	"strconv"
)

func (s *Service) GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := r.URL.Query()

	id := params.Get("id")
	//scale := params.Get("scale")

	img, err := s.db.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			jsonutil.MarshalResponse(w, http.StatusBadRequest, jsonutil.NewError(2, "Image not found"))
		} else {
			jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(5, "Failed to get image"))
		}

		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	w.Header().Set("Content-Length", strconv.Itoa(len(img.Image)))
	if _, err := w.Write(img.Image); err != nil {
		jsonutil.MarshalResponse(w, http.StatusInternalServerError, jsonutil.NewError(6, "Failed to return image"))
		return
	}
}
