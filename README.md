# README.md

# GitHub CLI

Ce projet est une interface de ligne de commande (CLI) qui récupère et affiche l'activité récente d'un utilisateur GitHub.

## Installation

Pour installer le projet, clonez le dépôt et exécutez la commande suivante :

```
go mod tidy
```

## Utilisation

Pour exécuter l'application, utilisez la commande suivante dans le terminal :

```
go run main.go <nom_utilisateur>
```

Remplacez `<nom_utilisateur>` par le nom d'utilisateur GitHub dont vous souhaitez voir l'activité.

## Exemples

```
go run main.go octocat
```

Cela affichera l'activité récente de l'utilisateur `octocat`.

## API GitHub

Ce projet interagit avec l'API GitHub pour récupérer les événements d'activité. Pour plus d'informations sur l'API, consultez la [documentation de l'API GitHub](https://docs.github.com/en/rest).
