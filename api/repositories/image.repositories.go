package repositories

import "database/sql"

type ImageRepository struct {
	db *sql.DB
}

func InitImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

// Create enregistre une image en base (nom unique, type MIME et contenu binaire).
func (r *ImageRepository) Create(nom string, contentType string, data []byte) error {
	query := `INSERT INTO images (nom, content_type, data) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, nom, contentType, data)
	return err
}

// FindByNom retrouve une image par son nom (ex : upload-123.png).
func (r *ImageRepository) FindByNom(nom string) (string, []byte, error) {
	query := `SELECT content_type, data FROM images WHERE nom = $1`
	var contentType string
	var data []byte
	err := r.db.QueryRow(query, nom).Scan(&contentType, &data)
	return contentType, data, err
}
