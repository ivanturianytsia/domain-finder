package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go"
)

var client pusher.Client

func main() {
	client = pusher.Client{
		AppId:   os.Getenv("PUSHER_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: "eu",
		Secure:  true,
	}

	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using evironmental variables.")
	}

	mux := http.NewServeMux()

	// t := template.Must(template.ParseFiles("cmd/index.html"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("cmd/index.html"))
		t.ExecuteTemplate(w, "index", map[string]string{
			"PusherKey": os.Getenv("PUSHER_KEY"),
		})
	})
	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		words := r.URL.Query()["word"]
		if len(words) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Word missing"))
			return
		}
		keys := r.URL.Query()["key"]
		if len(keys) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Key missing"))
			return
		}
		word := words[0]
		key := keys[0]
		go handleWord(word, key)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	})

	http.ListenAndServe(":8000", mux)
}

func handleWord(word, key string) {
	domainsCh := make(chan domain)
	doneCh := make(chan struct{})

	go findDomains(word, domainsCh, doneCh)

	for {
		select {
		case dom := <-domainsCh:
			client.Trigger("Domains", "Result-"+key, map[string]interface{}{
				"name":      dom.name,
				"available": dom.available,
				"err":       dom.err,
				"from":      dom.from,
			})
		case <-doneCh:
			client.Trigger("Domains", "End-"+key, nil)
			return
		}
	}
}
