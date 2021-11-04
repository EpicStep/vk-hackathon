package database

import (
	"context"
	"github.com/EpicStep/vk-hackathon/internal/image/model"
	"github.com/EpicStep/vk-hackathon/pkg/database"
)

type ImageDB struct {
	db *database.DB
}

func New(db *database.DB) *ImageDB {
	return &ImageDB{db: db}
}

func (db *ImageDB) Create(ctx context.Context, image *model.Image) error {
	_, err := db.db.DB.ExecContext(ctx, `
		INSERT INTO 
		    images
		    (id, image, hash, height, width)
		VALUES
			(?, ?, ?, ?, ?)
	`, image.ID, image.Image, image.Hash, image.Height, image.Width)

	if err != nil {
		return err
	}

	return nil
}

func (db *ImageDB) UpdateByID(ctx context.Context, image *model.Image) error {
	_, err := db.db.DB.ExecContext(ctx, `
		UPDATE 
		    images
		SET
		    image = ?,
			hash = ?,
			height = ?,
			width = ?
		WHERE
			id = ?
	`, image.Image, image.Hash, image.Height, image.Width, image.ID)

	if err != nil {
		return err
	}

	return nil
}


func (db *ImageDB) GetByID(ctx context.Context, id string) (*model.Image, error) {
	var i model.Image

	err := db.db.DB.QueryRowContext(ctx,`
		SELECT
			id, image
		FROM 
		    images
		WHERE
			id = ?
	`, id).Scan(&i.ID, &i.Image)

	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (db *ImageDB) GetAllImagesHash(ctx context.Context) ([]*model.Image, error) {
	var images []*model.Image

	rows, err := db.db.DB.QueryContext(ctx, `
		SELECT
			id, hash, height, width
		FROM 
		    images
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var m model.Image

		err := rows.Scan(&m.ID, &m.Hash, &m.Height, &m.Width)
		if err != nil {
			return nil, err
		}

		images = append(images, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}