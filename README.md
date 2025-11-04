# Mini-CRM CLI

Cette application est un mini-CRM (Customer Relationship Management) accessible en ligne de commande, permettant de gérer une liste de contacts.

## Fonctionnalités principales

- **Afficher les contacts** : Les contacts existants sont affichés à chaque tour du menu.
- **Ajouter un contact** : Saisissez le nom et l'email pour créer un nouveau contact.
- **Supprimer un contact** : Indiquez l'identifiant (ID) du contact à supprimer.
- **Mettre à jour un contact** : Modifiez le nom et/ou l'email d'un contact existant via son ID.
- **Quitter l'application** : Ferme le programme.

## Démarrage

1. Cloner le dépôt, puis lancez :

```bash
go run main.go
```

2. Utilisez le menu interactif pour accéder aux différentes fonctions (ajout, suppression, mise à jour, etc.).

## Ajout rapide via flags (optionnel)

Il est possible d’ajouter un contact directement via des options en ligne de commande :

```bash
go run main.go -add -name "Nom" -email "email@example.com"
```

## Concepts utilisés

- Boucle infinie pour le menu (`for {}`)
- Utilisation de `map` pour stocker les contacts
- Entrées utilisateur avec `bufio`, analyse et conversion avec `strconv`, gestion des erreurs avec `if err != nil`
- Sélection d’action avec `switch`
- "comma ok idiom" pour vérifier l’existence d’un contact

## Remarques

- Les contacts sont stockés en mémoire (pas de sauvegarde après l'arrêt).
- Les identifiants sont générés automatiquement.
