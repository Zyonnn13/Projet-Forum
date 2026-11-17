package config

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB() *sql.DB {
	// Sur Scalingo : SCALINGO_POSTGRESQL_URL (fournie par l'addon). En local : DATABASE_URL (.env).
	url := GetEnvWithDefault("SCALINGO_POSTGRESQL_URL", "")
	if url == "" {
		url = GetRequiredEnv("DATABASE_URL")
	}

	// Mise en place de la connexion.
	dbContext, dbContextErr := sql.Open("pgx", url)
	if dbContextErr != nil {
		log.Fatalf("Erreur connection base de donnees - Erreur : \n\t %s", dbContextErr.Error())
	}

	// Test de ping la base de donnees.
	pingErr := dbContext.Ping()
	if pingErr != nil {
		dbContext.Close()
		log.Fatalf("Erreur ping base de donnees - Erreur : \n\t %s", pingErr.Error())
	}

	log.Printf("BDD - Connexion reussie")
	return dbContext
}
