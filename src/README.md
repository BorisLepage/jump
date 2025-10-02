# Projet de Plateforme de Facturation

## Description

Projet de plateforme de facturation développée en Go avec une architecture microservices.
Le projet comprend une API de gestion des factures et une base de données PostgreSQL pour le stockage des données.

## Architecture
```
src/
├── facture/ # API de gestion des factures
│ ├── bin/facture/ # Point d'entrée de l'API
│ ├── internal/ # Logique repository
│ ├── Dockerfile # Image Docker pour l'API
│ └── docker-compose.yml # Orchestration des services
├── receipts/ # Schémas de base de données
└── helper_api/ # Utilitaires partagés
```
## Prérequis

- Docker et Docker Compose
- Go 1.24+ (pour le développement)
- Make

## Installation

1. **Cloner le projet**
   ```bash
   git clone <repository-url>
   cd jump
   ```

2. **Initialiser les dépendances Go**
   ```bash
   cd src
   go mod tidy
   ```

## Utilisation

### Démarrage de l'ensemble du projet

```bash
# Démarrer l'API facture avec sa base de données
make run pkg=facture
```

Cette commande va :
- Démarrer une instance PostgreSQL (receipts-db)
- Charger automatiquement le schéma de base de données
- Démarrer l'API facture sur le port 8080

### Commandes disponibles

```bash
# Afficher l'aide
make help

# Lancer les tests
make test pkg=facture

# Compiler le binaire
make build pkg=facture

# Nettoyer les artefacts
make clean pkg=facture

# Build de l'image Docker
make docker-build pkg=facture
```

### Accès aux services

- **API Facture** : http://localhost:8080
- **Base de données** : localhost:5432
  - Base : `receipts`
  - Utilisateur : `jump`
  - Mot de passe : `password`

## Endpoints disponibles

### GET /user/{id}
Récupérer les informations d'un utilisateur
```bash
curl http://localhost:8080/user/1
```

### POST /invoices
Créer une nouvelle facture
```bash
curl -X POST http://localhost:8080/invoices \
  -H "Content-Type: application/json" \
  -d '
  {
    "user_id": 1,
    "amount": 50000,
    "label": "Service de consultation"
  }'
```

### PUT /invoices/{id}/paid
Marquer une facture comme payée
```bash
curl -X PUT http://localhost:8080/invoices/1/_paid
```

## Développement

### Structure du projet

- **Monorepo** : Support de plusieurs APIs
- **Docker** : Chaque API a son propre Dockerfile
- **Makefile** : Orchestration centralisée
- **Tests** : Tests unitaires avec mocks

### Ajouter une nouvelle API

1. Créer le dossier `src/nouvelle-api/`
2. Ajouter un `Dockerfile` et `docker-compose.yml`
3. Utiliser les commandes make avec `pkg=nouvelle-api`

## Maintenance

### Logs
```bash
# Voir les logs de l'API
docker logs facture-api

# Voir les logs de la base de données
docker logs receipts-db
```

### Arrêt des services
```bash
make clean pkg=facture
```
