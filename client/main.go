// Package main démarre l'application cliente (front web) et configure le serveur HTTP.
package main

import (
	"log"
	"net/http"
	"os"

	"projet-forum/client/app"
	"projet-forum/client/config"
)

// main initialise l'application cliente et démarre le serveur HTTP.
func main() {
	// Le Procfile lance "<binaire> migrate" après chaque déploiement :
	// seule l'API a des migrations, le client n'a rien à faire.
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Client : aucune migration")
		return
	}

	app := app.InitApp()

	// Scalingo fournit le port dans $PORT ; 3000 par défaut en local.
	port := config.GetEnvWithDefault("PORT", "3000")
	log.Printf("Serveur Frontend (Client) lancé sur le port %s", port)
	serveErr := http.ListenAndServe(":"+port, app.Router)
	if serveErr != nil {
		log.Fatalf("Erreur lancement serveur client - %s", serveErr.Error())
	}
}
