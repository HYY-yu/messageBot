package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"messageBot/db"
	"messageBot/messenger"
	"messageBot/messenger/handlers"
	"messageBot/messenger/model"
	"net/http"
	"os"
	"os/signal"

	psh "github.com/platformsh/config-reader-go/v2"
)

func main() {
	prod := os.Getenv("IS_PROD")
	if prod == "true" {
		model.Prod = true
	}
	if model.Prod {
		model.VerifyToken = os.Getenv("VERIFY_TOKEN")
		model.AppSecret = os.Getenv("APP_SECRET")
		model.PageAccessToken = os.Getenv("PAGE_ACCESS_TOKEN")
		model.NLPToken = os.Getenv("NLP_TOKEN")
	}

	var srv http.Server
	srv.Addr = ":80"

	if model.Prod {
		// The Config Reader library provides Platform.sh environment information mapped to Go structs.
		config, err := psh.NewRuntimeConfig()
		if err != nil {
			panic("Not in a Platform.sh Environment.")
		}
		srv.Addr = config.Port()
	}

	// app_secret can read from config.
	webHooker := messenger.NewWebHooker(model.AppSecret)
	webHooker.AddMessageHandler(handlers.NewDatabaseMessageHandler(db.NewMessageRepository()))
	webHooker.AddMessageHandler(handlers.NewNLPHandler(model.NLPToken))
	webHooker.AddMessageHandler(handlers.NewRobotHandler(model.PageAccessToken, db.NewMessageTemplateRepository()))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, "Welcome to MessageBot!\n")
		if err != nil {
			log.Fatal(err)
		}
	}))

	http.Handle("/webhook", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// GET method is a webhook_verify request
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}

			mode := r.Form.Get("hub.mode")
			token := r.Form.Get("hub.verify_token")
			challenge := r.Form.Get("hub.challenge")

			if len(mode) != 0 && len(token) != 0 {
				if mode == "subscribe" && token == model.VerifyToken {
					log.Println("WEBHOOK_VERIFIED")
					w.WriteHeader(http.StatusOK)
					_, err := w.Write([]byte(challenge))
					if err != nil {
						log.Fatal(err)
					}
				} else {
					w.WriteHeader(http.StatusForbidden)
				}
			}
			return
		} else if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println("Error reading request body:", err)
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			// Verify message signature
			if !webHooker.VerifySignature(body, r.Header.Get("X-Hub-Signature")[5:]) {
				log.Println("invalid request signature")
				http.Error(w, "invalid request signature", http.StatusForbidden)
				return
			}

			var msg model.Message
			if err := json.Unmarshal(body, &msg); err != nil {
				log.Println("Error unmarshalling request body:", err)
				http.Error(w, "Error unmarshalling request body", http.StatusInternalServerError)
				return
			}

			// Process message
			go webHooker.HandleMessage(&msg)

			// return 200 OK as fast as possible
			// otherwise the fb server will send the message again.
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "Only support GET, POST methods", http.StatusMethodNotAllowed)
	}))

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		// The webHooker.Shutdown() must be called after the server.Shutdown().
		// Otherwise, It will lose the messages.
		webHooker.Shutdown()

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
