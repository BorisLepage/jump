# Skills Test - Backend Engineer

## Introduction

Merci par avance pour le temps que vous allez consacrer à cet exercice.

Via cet exercice nous voulons avant tout évaluer vos méthodes de travail et votre raisonnement.

Prenez tout le temps qu'il vous semblera nécessaire pour effectuer ce projet. L'objectif est de s'approcher le plus possible d'un résultat **"Production Ready"**. N'hésitez pas à ajouter par écrit vos idées si vous n'avez pas le temps de les implémenter.

Quand vous aurez fini, vous pouvez soit nous envoyer votre travail dans une archive compressée, soit nous envoyer un lien vers un dépôt de code public. Dans les deux cas, envoyez les informations nécessaires à l'adresse email `engineering+backend_interview@join-jump.com`.

## Contexte

Chez Jump, nous essayons d'automatiser au maximum les tâches administratives pour faciliter la vie de nos clients.

Tous nos clients facturent leurs prestations et le suivi de ces factures est une tâche longue et fastidieuse, mais nécessaire car c'est le règlement de ces factures qui leur permet par la suite l'avoir l'argent suffisant pour se verser un salaire.

Pour commencer vous allez devoir implémenter un outil de facturation capable de créer une facture et de modifier son statut.

## Instructions

### Généralités

- Votre service devra être écrit avec le langage **Golang**
- Pour la base de données, vous avez le choix (voir la section [Supports](#Supports)), chez Jump nous utilisons **PostgreSQL**
- Vous êtes libre d'utiliser les bibliothèques de votre choix
- Vous pouvez également fournir un document texte avec votre exercice (README, docx, Google Docs, …) pour décrire :
  - Vos choix d'implémentation
  - Votre approche
  - Les problèmes que vous avez rencontrés
  - Les améliorations que vous auriez aimé faire avec plus de temps
  - Si vous auriez fait les choses différemment de la façon dont elles sont demandées
- N'hésitez pas non plus à fournir des instructions sur comment lancer votre projet

> ⚠️ **Important** : La qualité du code doit être d'un niveau qui vous semble suffisant pour mettre votre projet en production, à vous de juger ce que vous considérez comme indispensable dans un projet pour être considéré comme propre à déployer en production.

### Implémentation

Votre service devra exposer au minimum les **3 routes suivantes** :

#### 1. Route GET - Récupérer un utilisateur

- **URL** : `/user/:id`
- **Méthode** : `GET`
- **Réponse** : JSON au format suivant

```jsonc
{
  "user_id": ...,
  "first_name": ...,
  "last_name": ...,
  "balance": ...
}
```

#### 2. Route POST - Créer une facture

- **URL** : `/invoices`
- **Méthode** : `POST`
- **Données d'entrée** : JSON au format suivant

```jsonc
{
  "user_id": ...,
  "amount": ...,
  "label": ...
}
```

#### 3. Route - Marquer une facture comme payée

- Une route pour mettre le statut `paid` à une facture
- La route doit renvoyer une erreur si la facture est déjà `paid`

> 📝 **Note** : Veillez à l'utilisation de codes HTTP appropriés.

## Supports

### Base de données

Un schéma de base de données est disponible dans le dossier `/database` de ce repository.

### Docker

Un Dockerfile est également disponible pour vous aider à lancer la base de données et tester votre projet.

Pour lancer la base de données :

```bash
git clone https://github.com/Freelance-launchpad/backend-interviews.git
cd backend-interviews/database
docker build . -t jump-database
docker run -p 5432:5432 jump-database
```

### Alternative

**Vous pouvez également utiliser votre propre schéma de base de données si vous le souhaitez**. Si c'est le cas, pensez à nous fournir un dockerfile ou les instructions nécessaires pour que nous puissions lancer la base de données de notre côté.

---

## Good luck ! 🚀