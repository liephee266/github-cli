package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Actor struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	DisplayLogin string `json:"display_login"`
	GravatarID   string `json:"gravatar_id"`
	URL          string `json:"url"`
	AvatarURL    string `json:"avatar_url"`
}

type Repo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Member struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Payload struct {
	Member Member `json:"member"`
	Action string `json:"action"`
}

type Event struct {
	ID        string  `json:"id"`
	Type      string  `json:"type"`
	Actor     Actor   `json:"actor"`
	Repo      Repo    `json:"repo"`
	Payload   Payload `json:"payload"`
	Public    bool    `json:"public"`
	CreatedAt string  `json:"created_at"`
}

type lineEventOutput struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Repo      struct {
		Name string `json:"name"`
	} `json:"repo"`
	Action string `json:"action"`
}

// main est le point d'entrée du programme
func main() {
	// Vérifie que l'utilisateur a fourni un nom d'utilisateur en argument de la ligne de commande
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-cli <username>")
		return
	}

	// Récupère le nom d'utilisateur depuis les arguments de la ligne de commande
	username := os.Args[1]
	// Construit l'URL de l'API GitHub pour récupérer les événements de l'utilisateur
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	// Fait une requête HTTP GET à l'URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	// Ferme le corps de la réponse à la fin de la fonction
	defer resp.Body.Close()

	// Vérifie que le code de statut de la réponse est 200 OK
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response code")
		return
	}

	// Décode la réponse JSON en une slice d'interfaces vides
	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	var eventOutput []lineEventOutput
	// Affiche chaque événement
	for _, event := range events {
		if event.Type == "PushEvent" {
			eventOutput = append(eventOutput, lineEventOutput{
				Type: event.Type,
				CreatedAt: func() time.Time {
					t, err := time.Parse(time.RFC3339, event.CreatedAt)
					if err != nil {
						fmt.Println("Error parsing time:", err)
					}
					return t
				}(),
				Repo: struct {
					Name string `json:"name"`
				}{
					Name: event.Repo.Name,
				},
				Action: event.Payload.Action,
			})
		}
	}

	// Affiche toutes les ligne d'événements trouvés
	for _, event := range eventOutput {
		fmt.Println(event.Type, event.CreatedAt, event.Repo.Name, event.Action)
	}
}
