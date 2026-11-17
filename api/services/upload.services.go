package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"projet-forum/api/repositories"
)

// Extensions d'images autorisées à l'upload.
var allowedImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadService enregistre les images en base de données
// (le disque d'un conteneur Scalingo est effacé à chaque redéploiement).
type UploadService struct {
	imageRepository *repositories.ImageRepository
}

func InitUploadService(imageRepository *repositories.ImageRepository) *UploadService {
	return &UploadService{imageRepository: imageRepository}
}

// SaveImage valide l'image, l'enregistre en base sous un nom unique
// et retourne son chemin public (ex : /uploads/upload-123.png).
func (s *UploadService) SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExtensions[ext] {
		return "", fmt.Errorf("format non autorisé (jpg, jpeg, png, gif, webp uniquement)")
	}

	// Limite de 5 Mo par image.
	if header.Size > 5<<20 {
		return "", fmt.Errorf("image trop lourde (5 Mo maximum)")
	}

	// Lecture du contenu du fichier.
	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la lecture du fichier : %s", err.Error())
	}

	// Nom unique basé sur l'heure actuelle en nanosecondes.
	nom := fmt.Sprintf("upload-%d%s", time.Now().UnixNano(), ext)

	// DetectContentType devine le type MIME (image/png, image/jpeg...) à partir du contenu.
	err = s.imageRepository.Create(nom, http.DetectContentType(data), data)
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'enregistrement de l'image : %s", err.Error())
	}

	return "/uploads/" + nom, nil
}

// GetImage retourne le type MIME et le contenu d'une image enregistrée.
func (s *UploadService) GetImage(nom string) (string, []byte, error) {
	return s.imageRepository.FindByNom(nom)
}
