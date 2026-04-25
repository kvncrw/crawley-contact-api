package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Company string `json:"company"`
	Message string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	pushoverToken := os.Getenv("PUSHOVER_TOKEN")
	pushoverUser := os.Getenv("PUSHOVER_USER")
	discordWebhook := os.Getenv("DISCORD_WEBHOOK_URL")

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/contact", func(w http.ResponseWriter, r *http.Request) {
		var req ContactRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		if req.Name == "" || req.Email == "" || req.Message == "" {
			http.Error(w, "name, email, and message are required", http.StatusBadRequest)
			return
		}

		log.Printf("contact: name=%q email=%q company=%q",
			req.Name, req.Email, req.Company)

		msg := formatMessage(req)

		if discordWebhook != "" {
			if err := sendDiscord(discordWebhook, msg); err != nil {
				log.Printf("discord send failed: %v", err)
			}
		}

		if pushoverToken != "" && pushoverUser != "" {
			if err := sendPushover(pushoverToken, pushoverUser, req.Name, msg); err != nil {
				log.Printf("pushover send failed: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"https://crawley.systems"},
		AllowedMethods: []string{"POST"},
		AllowedHeaders: []string{"Content-Type"},
	}).Handler(mux)

	log.Printf("contact-api listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func formatMessage(req ContactRequest) string {
	company := req.Company
	if company == "" {
		company = "(not provided)"
	}
	return fmt.Sprintf("**New Contact — crawley.systems**\n\n**Name:** %s\n**Email:** %s\n**Company:** %s\n\n**Message:**\n%s",
		req.Name, req.Email, company, req.Message)
}

func sendDiscord(webhookURL, message string) error {
	body, _ := json.Marshal(map[string]string{"content": message})
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord returned %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func sendPushover(token, user, title, message string) error {
	body, _ := json.Marshal(map[string]string{
		"token":   token,
		"user":    user,
		"title":   title,
		"message": message,
	})
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post("https://api.pushover.net/1/messages.json", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pushover returned %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
