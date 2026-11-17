# Projet Forum — Plateforme communautaire « Forum Gaming »

Plateforme de discussion communautaire sur le thème du **jeu vidéo**, inspirée de JeuxVideos.com. Les utilisateurs peuvent s'inscrire, ouvrir des fils de discussion, publier des
messages, réagir (like / dislike), filtrer par catégorie, rechercher des sujets, et les
administrateurs disposent d'un tableau de bord de modération.

---


## Membres de l'équipe

| Nom | Rôle |
|-----|------|
| **Clément Belmondo** | Développement |
| **Samuel Berard** | Développement |

---

## Technologies utilisées

| Domaine | Technologie |
|---------|-------------|
| Langage backend | **Go**|
| Routeur HTTP | `gorilla/mux` |
| Authentification | JWT — `golang-jwt/jwt/v5` |
| Base de données | **PostgreSQL** (driver `pgx`) |
| Migrations | **goose** (fichiers SQL versionnés) |
| Hachage des mots de passe | **SHA-512** |
| Hébergement | **Scalingo** (2 apps + addon PostgreSQL) |

---

## Architecture

Le projet est un **client serveur organisé en couches** , regroupé dans un
seul dépôt et un seul module Go. Il est découpé en deux serveurs HTTP :

- **API** (`api/`, port 8080 en local) — renvoie du JSON sous `/api/...` ; c'est la seule à parler à la base.
- **Client** (`client/`, port 3000 en local) — génère les pages HTML et appelle l'API côté serveur (variable `API_BASE_URL`).

Chaque serveur respecte une séparation claire des responsabilités :

- **Router** — définition des routes.
- **Controller** — réception des requêtes, validation des entrées, appel de la logique métier.
- **Service** — logique métier.
- **Repository** — accès aux données et communication avec PostgreSQL.
- **Model / DTO** — représentation des données manipulées.
- **View** — affichage des pages HTML (côté client).

---



## Prérequis

Avant d'installer le projet, assurez-vous d'avoir :

- **Go** 1.26 ou supérieur
- **Docker** (Docker Desktop) — pour lancer PostgreSQL en local, comme en production
- **Git** — pour cloner le dépôt

Vérification rapide :

```bash
go version
docker --version
```

---

## Installation

### 1. Cloner le dépôt

```bash
git clone https://github.com/Zyonnn13/Projet-Forum.git
cd Projet-Forum
```

### 2. Lancer PostgreSQL (Docker)

```bash
docker compose up -d
```

### 3. Configurer les variables d'environnement

```bash
cp .env.example .env
```

Puis remplacez `JWT_SECRET=a-remplacer` dans `.env` par une valeur longue et aléatoire
(par exemple le résultat de `openssl rand -base64 32`). Le fichier `.env` n'est jamais commité.

### 4. Créer les tables et insérer les données de test

```bash
go run ./api migrate
go run ./api seed
```

`migrate` applique les migrations de `api/migrations/` (goose) ; `seed` insère les catégories,
utilisateurs, fils et messages de démonstration (`api/seed/seed.sql`). Les deux sont relançables sans risque.

---

## Lancement de l'application

L'application nécessite **deux serveurs** lancés en parallèle (dans deux terminaux).

**Terminal 1 — API (backend, port 8080) :**

```bash
go run ./api
```

**Terminal 2 — Client (frontend, port 3000) :**

```bash
go run ./client
```

Ouvrez ensuite votre navigateur sur :

> **http://localhost:3000**

L'API expose aussi `GET http://localhost:8080/health`, qui répond `{"status":"ok"}` si l'API et la base répondent.

---

## Comptes de démonstration

Les données de test (`api/seed/seed.sql`) créent les comptes suivants.
Mot de passe commun : **`Password!123`**

| Identifiant | Email | Rôle |
|-------------|-------|------|
| `admin` | admin@example.com | **administrateur** |
| `alex_gaming` | alex@example.com | utilisateur |
| `sophie_dev` | sophie@example.com | utilisateur |
| `ludo_retro` | ludo@example.com | utilisateur |
| `chloe_fps` | chloe@example.com | utilisateur |
| `max_hardware` | max@example.com | utilisateur |

---

## Structure du projet

```
Projet-Forum/
├── api/                  # Serveur API (backend, :8080)
│   ├── main.go
│   ├── app/              # Assemblage de l'application (injection de dépendances)
│   ├── config/           # Chargement .env et connexion PostgreSQL
│   ├── migrations/       # Migrations SQL versionnées (goose), intégrées au binaire
│   ├── seed/             # Données de démonstration (go run ./api seed)
│   ├── routers/          # Définition des routes /api/...
│   ├── controllers/      # Réception des requêtes, validation
│   ├── services/         # Logique métier
│   ├── repositories/     # Accès aux données (SQL)
│   ├── models/           # Entités
│   ├── dto/              # Objets de transfert de données
│   ├── middleware/       # Authentification JWT, contrôle des bannis
│   ├── auth/             # Génération et validation des JWT
│   ├── helper/ utils/    # Fonctions utilitaires (dont hachage SHA-512)
│   └── ...
├── client/               # Serveur web (frontend, :3000)
│   ├── main.go
│   ├── app/              # Assemblage de l'application cliente
│   ├── routers/          # Routes des pages + assets
│   ├── controllers/      # Rendu des pages
│   ├── services/         # Appels métier
│   ├── api/              # Client HTTP vers l'API backend
│   ├── templates/        # Pages HTML (index, forum, thread, login, admin...)
│   └── config/
├── static/               # CSS, JS et images servis sur /static
│   ├── css/
│   ├── js/
│   └── img/
├── .github/workflows/    # CI GitHub Actions (api.yml, client.yml)
├── Procfile              # Commandes lancées par Scalingo (web + postdeploy)
├── docker-compose.yml    # PostgreSQL local
├── .env.example          # Modèle des variables d'environnement (sans secret)
├── go.mod / go.sum
└── README.md
```

### Schéma de la base de données

`utilisateurs`, `categories`, `fils_de_discussion`, `fils_categories` (liaison N–N),
`messages`, `reactions`, `images`. Les suppressions sont propagées via `ON DELETE CASCADE`
(supprimer un fil supprime ses messages et leurs réactions). Un trigger met à jour
`updated_at` à chaque modification d'un fil ou d'un message.

Les images envoyées (avatars) sont stockées **dans la table `images`** et servies par
`GET /uploads/{nom}` : le disque d'un conteneur Scalingo est effacé à chaque redéploiement.

---

## Mise en production (Scalingo)

### Architecture

```
Navigateur ──HTTPS──▶ App client (bin/client) ──HTTPS──▶ App API (bin/api) ──SQL──▶ PostgreSQL
                                                              ▲                    (addon Starter 512)
Navigateur ──<img src="https://<app-api>/uploads/…">──────────┘
```

- Les deux apps déploient **le même dépôt**. `go.mod` compile les deux binaires
  (`// +scalingo install ./api ./client`) et le `Procfile` lance celui indiqué par `APP_BIN`.
- Pas de `PROJECT_DIR` : le client lit `./static` et `./client/templates` depuis la racine du dépôt.
- Pas de CORS : le navigateur n'appelle jamais l'API en JavaScript. Il ne charge depuis l'API que
  des images via `<img>`, ce qui ne nécessite pas CORS.
- Après chaque déploiement, le process `postdeploy` lance `bin/api migrate` (côté client, `migrate` ne fait rien).

### Variables d'environnement

| Variable | App | Rôle |
|----------|-----|------|
| `APP_BIN` | API et client | Binaire lancé par le Procfile : `api` ou `client` |
| `JWT_SECRET` | API | Clé de signature des JWT |
| `SCALINGO_POSTGRESQL_URL` | API | URL de la base, fournie par l'addon PostgreSQL |
| `DATABASE_URL` | API (local) | URL de la base locale (Docker) |
| `API_BASE_URL` | client | URL de l'API, ex. `https://<app-api>.osc-fr1.scalingo.io/api` |
| `PORT` | API et client | Port d'écoute, fourni par Scalingo |

### Procédure de déploiement

API et base :

```bash
scalingo create --remote scalingo-api <app-api>
scalingo --app <app-api> addons-add postgresql postgresql-starter-512
scalingo --app <app-api> env-set APP_BIN=api JWT_SECRET=<clé générée avec openssl rand -base64 32>
git push scalingo-api main
scalingo --app <app-api> scale web:1:S
scalingo --app <app-api> run bin/api seed     # une seule fois
```

Client :

```bash
scalingo create --remote scalingo-client <app-client>
scalingo --app <app-client> env-set APP_BIN=client API_BASE_URL=https://<app-api>.osc-fr1.scalingo.io/api
git push scalingo-client main
scalingo --app <app-client> scale web:1:S
```

Contrôle : `https://<app-api>.osc-fr1.scalingo.io/health` répond `{"status":"ok"}`.

### Déploiement continu

Chaque app est liée au dépôt GitHub (onglet **Deploy** du dashboard) avec **Auto Deploy** sur `main`.
Les workflows `.github/workflows/api.yml` et `client.yml` lancent format (`gofmt`), lint (`go vet`),
migrations et tests ; Scalingo attend qu'ils soient verts avant de déployer.

### Coûts

| Ressource | Estimation mensuelle | Réel |
|-----------|----------------------|------|
| 2 conteneurs S | _à compléter (simulateur Scalingo)_ | _à compléter_ |
| PostgreSQL Starter 512 | _à compléter_ | _à compléter_ |
