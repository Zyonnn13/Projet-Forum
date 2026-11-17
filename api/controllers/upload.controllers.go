package controllers

import (
	"net/http"

	"projet-forum/api/helper"
	"projet-forum/api/services"

	"github.com/gorilla/mux"
)

type UploadControllers struct {
	service *services.UploadService
}

func InitUploadController(service *services.UploadService) *UploadControllers {
	return &UploadControllers{service: service}
}

// Upload reçoit un fichier en multipart/form-data (champ "file"),
// l'enregistre en base de données et renvoie son URL publique.
func (c *UploadControllers) Upload(w http.ResponseWriter, r *http.Request) {
	// Taille maximale en mémoire : 10 Mo.
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Fichier trop volumineux ou requête invalide")
		return
	}

	// Récupération du fichier envoyé sous le champ "file".
	file, header, err := r.FormFile("file")
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, "Aucun fichier reçu (champ 'file' attendu)")
		return
	}
	defer file.Close()

	chemin, err := c.service.SaveImage(file, header)
	if err != nil {
		helper.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// URL absolue servie par l'API, utilisable directement dans une balise <img>.
	// Sur Scalingo, le HTTPS est géré en amont : il nous l'indique dans X-Forwarded-Proto.
	scheme := "http"
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	urlPublique := scheme + "://" + r.Host + chemin

	helper.WriteJSON(w, http.StatusCreated, map[string]string{"url": urlPublique})
}

// ServeImage renvoie une image enregistrée en base (ex : GET /uploads/upload-123.png).
func (c *UploadControllers) ServeImage(w http.ResponseWriter, r *http.Request) {
	nom := mux.Vars(r)["nom"]

	contentType, data, err := c.service.GetImage(nom)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
