package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Commit représente un commit GitHub dans un événement de type PushEvent.
type Commit struct {
	Sha    string `json:"sha"`    // Identifiant du commit (SHA)
	Author struct {
		Name string `json:"name"` // Nom de l'auteur du commit
	} `json:"author"`
	Message string `json:"message"` // Message du commit
}

// Actor représente l'acteur qui a effectué l'événement sur GitHub (utilisateur).
type Actor struct {
	Login string `json:"login"` // Nom d'utilisateur de l'acteur
}

// Repo représente un dépôt GitHub sur lequel l'événement a eu lieu.
type Repo struct {
	Name string `json:"name"` // Nom du dépôt
}

// Payload représente la charge utile d'un événement GitHub.
// Cela peut inclure l'action effectuée ainsi que les commits pour un PushEvent.
type Payload struct {
	Action  string   `json:"action"`  // Action effectuée (ex. "added", "removed")
	Commits []Commit `json:"commits"` // Liste des commits dans un PushEvent
}

// Event représente un événement GitHub générique.
type Event struct {
	Type      string  `json:"type"`      // Type d'événement (ex. "PushEvent")
	Actor     Actor   `json:"actor"`     // Acteur ayant réalisé l'événement
	Repo      Repo    `json:"repo"`      // Dépôt concerné par l'événement
	Payload   Payload `json:"payload"`   // Charge utile de l'événement
	CreatedAt string  `json:"created_at"` // Date de création de l'événement
}

// LineEventOutput représente l'événement sous une forme plus simple à afficher.
type LineEventOutput struct {
	Type      string    `json:"type"`      // Type d'événement
	CreatedAt time.Time `json:"created_at"` // Date de création formatée de l'événement
	RepoName  string    `json:"repo_name"`  // Nom du dépôt
	Action    string    `json:"action"`     // Action effectuée dans l'événement
	Commits   []Commit  `json:"commits"`    // Liste des commits pour un PushEvent
}

// fetchGitHubEvents récupère les événements d'un utilisateur GitHub à partir de l'API.
func fetchGitHubEvents(username string) ([]Event, error) {
	// Construction de l'URL de l'API GitHub pour récupérer les événements de l'utilisateur
	apiURL := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	// Envoi de la requête HTTP GET à GitHub
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub events: %w", err)
	}
	defer resp.Body.Close() // Assurez-vous de fermer la réponse à la fin

	// Vérification du code de statut HTTP
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from GitHub API: %s", resp.Status)
	}

	// Décodage de la réponse JSON dans une slice d'événements
	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return events, nil
}

// parseGitHubEvents transforme les événements récupérés en une forme plus simple pour l'affichage.
func parseGitHubEvents(events []Event) []LineEventOutput {
	var parsedEvents []LineEventOutput

	// Pour chaque événement GitHub
	for _, event := range events {
		// Conversion de la date en format time.Time
		createdAt, err := time.Parse(time.RFC3339, event.CreatedAt)
		if err != nil {
			// Si la date est mal formatée, ignorer cet événement
			fmt.Printf("Error parsing time for event (%s): %v\n", event.Type, err)
			continue
		}
		// Ajouter l'événement transformé dans la liste des événements parsés
		parsedEvents = append(parsedEvents, LineEventOutput{
			Type:      event.Type,
			CreatedAt: createdAt,
			RepoName:  event.Repo.Name,
			Action:    event.Payload.Action,
			Commits:   event.Payload.Commits, // Ajouter les commits si c'est un PushEvent
		})
	}

	return parsedEvents
}

// displayEvents affiche les événements GitHub sous une forme lisible.
func displayEvents(events []LineEventOutput) {
	// Titre pour la sortie des événements
	fmt.Println("GitHub Events Summary:")
	fmt.Println("======================")
	// Parcours de chaque événement et affichage de ses détails
	for _, event := range events {
		// Affichage des informations de base de l'événement
		fmt.Printf("- %s | %s | Repo: %s | Action: %s\n",
			event.Type,
			event.CreatedAt.Format("2006-01-02 15:04:05"), // Format de date lisible
			event.RepoName,
			event.Action,
		)

		// Si des commits sont associés à cet événement (ex. PushEvent), on les affiche
		if len(event.Commits) > 0 {
			fmt.Println("  Commits:")
			for _, commit := range event.Commits {
				// Affichage des détails de chaque commit
				fmt.Printf("    - %s: %s (%s)\n",
					commit.Sha[:7],          // Afficher seulement les 7 premiers caractères du SHA
					commit.Message,
					commit.Author.Name,      // Auteur du commit
				)
			}
		}
	}
}

// main est le point d'entrée du programme. Il récupère les événements GitHub pour un utilisateur.
func main() {
	// Vérifie que l'utilisateur a fourni un nom d'utilisateur via la ligne de commande
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-cli <username>")
		os.Exit(1)
	}

	// Récupère le nom d'utilisateur depuis les arguments de la ligne de commande
	username := os.Args[1]

	// Récupère les événements GitHub de l'utilisateur
	events, err := fetchGitHubEvents(username)
	if err != nil {
		// Si une erreur survient, affiche un message et quitte le programme
		fmt.Printf("Error fetching events: %v\n", err)
		os.Exit(1)
	}

	// Parse les événements récupérés pour les transformer en une forme simple à afficher
	parsedEvents := parseGitHubEvents(events)

	// Affiche les événements parsés
	displayEvents(parsedEvents)
}
