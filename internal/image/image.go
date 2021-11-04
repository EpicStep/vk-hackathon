package image

import (
	imageDB "github.com/EpicStep/vk-hackathon/internal/image/database"
	"github.com/EpicStep/vk-hackathon/pkg/database"
)

type Service struct {
	db *imageDB.ImageDB
}

func New(db *database.DB) *Service {
	return &Service{db: imageDB.New(db)}
}