package repositories

import (
	"database/sql"
	"errors"
	"time"

	"projet-forum/api/models"
	"projet-forum/api/utils"
)

type FilRepository struct {
	db *sql.DB
}

func InitFilRepository(db *sql.DB) *FilRepository {
	return &FilRepository{db: db}
}

func (r *FilRepository) Create(f *models.FilDeDiscussion, categoryIDs []int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	query := `INSERT INTO fils_de_discussion (titre, etat, auteur_id) VALUES (?, ?, ?) RETURNING id`
	var id int64
	err = tx.QueryRow(utils.Rebind(query), f.Titre, f.Etat, f.AuteurID).Scan(&id)
	if err != nil {
		tx.Rollback()
		return err
	}
	f.ID = int(id)

	if len(categoryIDs) > 0 {
		catQuery := `INSERT INTO fils_categories (fil_id, categorie_id) VALUES (?, ?)`
		for _, catID := range categoryIDs {
			_, err = tx.Exec(utils.Rebind(catQuery), f.ID, catID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *FilRepository) FindAll(page, limit int, search string, searchType string, categoryID int) ([]models.FilDeDiscussion, error) {
	offset := (page - 1) * limit
	args := []interface{}{}

	query := `
		SELECT DISTINCT f.id, f.titre, f.etat, f.auteur_id, f.created_at, f.updated_at, u.username, u.email, u.role, u.banned
		FROM fils_de_discussion f
		LEFT JOIN utilisateurs u ON f.auteur_id = u.id
	`

	if categoryID > 0 {
		query += ` INNER JOIN fils_categories fc ON f.id = fc.fil_id `
	}

	query += ` WHERE f.etat != 'archive' `

	if search != "" {
		if searchType == "authors" {
			query += ` AND u.username ILIKE ? `
		} else if searchType == "messages" {
			query += ` AND f.id IN (SELECT DISTINCT fil_id FROM messages WHERE contenu ILIKE ?) `
		} else {
			query += ` AND f.titre ILIKE ? `
		}
		args = append(args, "%"+search+"%")
	}

	if categoryID > 0 {
		query += ` AND fc.categorie_id = ? `
		args = append(args, categoryID)
	}

	query += ` ORDER BY f.created_at DESC LIMIT ? OFFSET ? `
	args = append(args, limit, offset)

	rows, err := r.db.Query(utils.Rebind(query), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fils []models.FilDeDiscussion
	for rows.Next() {
		var f models.FilDeDiscussion
		var createdAt, updatedAt []uint8
		var authorUsername, authorEmail, authorRole string
		var authorBanned bool
		if err := rows.Scan(&f.ID, &f.Titre, &f.Etat, &f.AuteurID, &createdAt, &updatedAt, &authorUsername, &authorEmail, &authorRole, &authorBanned); err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, string(createdAt)); err == nil {
			f.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339Nano, string(updatedAt)); err == nil {
			f.UpdatedAt = t
		}
		f.Auteur = &models.Utilisateur{
			ID:       f.AuteurID,
			Username: authorUsername,
			Email:    authorEmail,
			Role:     authorRole,
			Banned:   authorBanned,
		}
		fils = append(fils, f)
	}

	return fils, nil
}

func (r *FilRepository) CountAll(search string, searchType string, categoryID int) (int, error) {
	args := []interface{}{}
	query := `
		SELECT COUNT(DISTINCT f.id) 
		FROM fils_de_discussion f
		LEFT JOIN utilisateurs u ON f.auteur_id = u.id
	`

	if categoryID > 0 {
		query += ` INNER JOIN fils_categories fc ON f.id = fc.fil_id `
	}

	query += ` WHERE f.etat != 'archive' `

	if search != "" {
		if searchType == "authors" {
			query += ` AND u.username ILIKE ? `
		} else if searchType == "messages" {
			query += ` AND f.id IN (SELECT DISTINCT fil_id FROM messages WHERE contenu ILIKE ?) `
		} else {
			query += ` AND f.titre ILIKE ? `
		}
		args = append(args, "%"+search+"%")
	}

	if categoryID > 0 {
		query += ` AND fc.categorie_id = ? `
		args = append(args, categoryID)
	}

	var count int
	err := r.db.QueryRow(utils.Rebind(query), args...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindAllAdmin retourne tous les fils, quel que soit leur etat (y compris archives),
// pour le tableau de bord d'administration.
func (r *FilRepository) FindAllAdmin() ([]models.FilDeDiscussion, error) {
	query := `
		SELECT f.id, f.titre, f.etat, f.auteur_id, f.created_at, f.updated_at, u.username, u.email, u.role, u.banned
		FROM fils_de_discussion f
		LEFT JOIN utilisateurs u ON f.auteur_id = u.id
		ORDER BY f.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fils []models.FilDeDiscussion
	for rows.Next() {
		var f models.FilDeDiscussion
		var createdAt, updatedAt []uint8
		var authorUsername, authorEmail, authorRole string
		var authorBanned bool
		if err := rows.Scan(&f.ID, &f.Titre, &f.Etat, &f.AuteurID, &createdAt, &updatedAt, &authorUsername, &authorEmail, &authorRole, &authorBanned); err != nil {
			return nil, err
		}
		if t, err := time.Parse(time.RFC3339Nano, string(createdAt)); err == nil {
			f.CreatedAt = t
		}
		if t, err := time.Parse(time.RFC3339Nano, string(updatedAt)); err == nil {
			f.UpdatedAt = t
		}
		f.Auteur = &models.Utilisateur{
			ID:       f.AuteurID,
			Username: authorUsername,
			Email:    authorEmail,
			Role:     authorRole,
			Banned:   authorBanned,
		}
		fils = append(fils, f)
	}
	return fils, nil
}

func (r *FilRepository) FindByID(id int) (*models.FilDeDiscussion, error) {
	query := `
		SELECT f.id, f.titre, f.etat, f.auteur_id, f.created_at, f.updated_at, u.username, u.email, u.role, u.banned
		FROM fils_de_discussion f
		LEFT JOIN utilisateurs u ON f.auteur_id = u.id
		WHERE f.id = ?
	`
	row := r.db.QueryRow(utils.Rebind(query), id)

	f := &models.FilDeDiscussion{}
	var createdAt, updatedAt []uint8
	var authorUsername, authorEmail, authorRole string
	var authorBanned bool
	err := row.Scan(&f.ID, &f.Titre, &f.Etat, &f.AuteurID, &createdAt, &updatedAt, &authorUsername, &authorEmail, &authorRole, &authorBanned)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("fil introuvable")
		}
		return nil, err
	}
	if t, err := time.Parse(time.RFC3339Nano, string(createdAt)); err == nil {
		f.CreatedAt = t
	}
	if t, err := time.Parse(time.RFC3339Nano, string(updatedAt)); err == nil {
		f.UpdatedAt = t
	}
	f.Auteur = &models.Utilisateur{
		ID:       f.AuteurID,
		Username: authorUsername,
		Email:    authorEmail,
		Role:     authorRole,
		Banned:   authorBanned,
	}

	return f, nil
}

func (r *FilRepository) Update(f *models.FilDeDiscussion) error {
	query := `UPDATE fils_de_discussion SET titre = ?, etat = ? WHERE id = ?`
	_, err := r.db.Exec(utils.Rebind(query), f.Titre, f.Etat, f.ID)
	return err
}

func (r *FilRepository) Delete(id int) error {
	query := `DELETE FROM fils_de_discussion WHERE id = ?`
	_, err := r.db.Exec(utils.Rebind(query), id)
	return err
}
