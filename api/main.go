package main

import (
	"log"
	"net/http"
	"os"

	"projet-forum/api/app"
	"projet-forum/api/config"
	"projet-forum/api/migrations"
	"projet-forum/api/seed"

	"github.com/pressly/goose/v3"
)

func main() {
	// Sous-commandes : "api migrate", "api rollback", "api seed"
	if len(os.Args) > 1 {
		runCommand(os.Args[1])
		return
	}

	app := app.InitApp()
	defer app.Close()

	// Scalingo fournit le port dans $PORT ; 8080 par défaut en local.
	port := config.GetEnvWithDefault("PORT", "8080")
	log.Printf("API lancée sur le port %s", port)
	if err := http.ListenAndServe(":"+port, app.Router); err != nil {
		log.Fatalf("Erreur lancement serveur - %s", err.Error())
	}
}

// runCommand exécute "migrate", "rollback" ou "seed" puis s'arrête.
func runCommand(cmd string) {
	config.LoadEnv()
	db := config.InitDB()
	defer db.Close()

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	var err error
	switch cmd {
	case "migrate":
		err = goose.Up(db, ".")
	case "rollback":
		err = goose.Down(db, ".")
	case "seed":
		_, err = db.Exec(seed.SQL)
	default:
		log.Fatalf("Commande inconnue : %s (migrate, rollback, seed)", cmd)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s : OK", cmd)
}
