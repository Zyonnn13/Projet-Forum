package app

import (
	"database/sql"

	"projet-forum/api/config"
	"projet-forum/api/controllers"
	"projet-forum/api/middleware"
	"projet-forum/api/repositories"
	"projet-forum/api/routers"
	"projet-forum/api/services"

	"github.com/gorilla/mux"
)

type App struct {
	Db     *sql.DB
	Router *mux.Router
}

func InitApp() *App {
	// Chargement des variables d'environnement
	config.LoadEnv()

	// Initialisation de la connexion à la base de données
	db := config.InitDB()

	// Initialisation de la connexion DB pour le middleware
	middleware.Init(db)

	// Initialisation des repositories
	utilisateurRepository := repositories.InitUtilisateurRepository(db)
	filRepository := repositories.InitFilRepository(db)
	messageRepository := repositories.InitMessageRepository(db)
	imageRepository := repositories.InitImageRepository(db)

	// Initialisation des services
	utilisateurService := services.InitUtilisateurService(utilisateurRepository)
	filService := services.InitFilService(filRepository)
	messageService := services.InitMessageService(messageRepository, filRepository)
	actusService := services.InitActusService()
	uploadService := services.InitUploadService(imageRepository)

	// Initialisation des controllers
	utilisateurController := controllers.InitUtilisateurController(utilisateurService)
	filController := controllers.InitFilController(filService)
	messageController := controllers.InitMessageController(messageService)
	actusController := controllers.InitActusController(actusService)
	uploadController := controllers.InitUploadController(uploadService)
	healthController := controllers.InitHealthController(db)

	// Routeur racine : santé (/health), images (/uploads/...) et sous-routeur /api.
	rootRouter := mux.NewRouter()
	rootRouter.HandleFunc("/health", healthController.Health).Methods("GET")
	rootRouter.HandleFunc("/uploads/{nom}", uploadController.ServeImage).Methods("GET")

	// Enregistrement des routes de l'API (avec ajout du préfixe "/api/...")
	router := rootRouter.PathPrefix("/api").Subrouter()

	routers.RegisterUtilisateurRoutes(router, utilisateurController)
	routers.RegisterFilRoutes(router, filController)
	routers.RegisterMessageRoutes(router, messageController)
	routers.RegisterActusRoutes(router, actusController)
	routers.RegisterUploadRoutes(router, uploadController)

	return &App{
		Db:     db,
		Router: rootRouter,
	}
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
}
