package repository

import (
	"database/sql"

	"github.com/alekslesik/neuro-news/internal/app/model"
	"github.com/alekslesik/neuro-news/pkg/logger"
)

type MySQLImageRepository struct {
	db *sql.DB
	l  *logger.Logger
}

type ImageQueries struct {
	insert                 string
	selectArticleRandLimit string
}

var imageQueries = ImageQueries{
	insert: `INSERT INTO image
	(image_path, image_size, image_name, image_alt)
	VALUES(?, ?, ?, ?);`,

	selectArticleRandLimit: `SELECT image_path, image_alt
	FROM image
	ORDER BY RAND()
	LIMIT ?;`,
}

// InsertImage insert image to DB
func (r *MySQLImageRepository) InsertImage(image *model.Image) error {
	const op = "repository.InsertImage()"

	result, err := r.db.Exec(imageQueries.insert, image.ImagePath, image.Size, image.Name, image.Alt)
	if err != nil {
		r.l.Warn().Msgf("%s: query exec insert image error > %s", op, err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.l.Warn().Msgf("%s: query exec insert image row affected error > %s", op, err)
		return err
	}

	if rows != 1 {
		r.l.Warn().Msgf("%s: query exec insert image number affected rows is > %d", op, rows)
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.l.Warn().Msgf("%s: query exec insert image get id error > %d", op, rows)
		return err
	}

	image.ImageID = id

	return nil
}

// SelectGalleryImage return random images
func (r *MySQLImageRepository) SelectGalleryImage(limit int) ([]model.Image, error) {
	const op = "repository.SelectGalleryImage()"
	var is = make([]model.Image, 0, limit)

	rows, err := r.db.Query(imageQueries.selectArticleRandLimit, limit)
	if err != nil {
		r.l.Warn().Msgf("%s: query select images for gallery error > %s", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var i model.Image
		err = rows.Scan(&i.ImagePath, &i.Alt)
		if err != nil {
			r.l.Error().Msgf("%s: query scan images for gallery error > %s", op, err)
			return nil, err
		}

		is = append(is, i)
	}

	return is, nil
}
