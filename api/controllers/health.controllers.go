package controllers

import (
	"database/sql"
	"net/http"

	"projet-forum/api/helper"
)

type HealthControllers struct {
	db *sql.DB
}

func InitHealthController(db *sql.DB) *HealthControllers {
	return &HealthControllers{db: db}
}

// Health répond 200 si l'API tourne et que la base de données répond.
func (c *HealthControllers) Health(w http.ResponseWriter, r *http.Request) {
	if err := c.db.Ping(); err != nil {
		helper.WriteError(w, http.StatusServiceUnavailable, "Base de données indisponible")
		return
	}
	helper.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
