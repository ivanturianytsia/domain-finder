package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go"
)

var client pusher.Client

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using evironmental variables.")
	}

	client = pusher.Client{
		AppId:   os.Getenv("PUSHER_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: "eu",
		Secure:  true,
	}

	mux := http.NewServeMux()

	// t := template.Must(template.ParseFiles("app/views/index.html"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("app/views/index.html"))
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
		key := strconv.Itoa(rand.Int())

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(key))

		go handleWord(words[0], key)
	})
	mux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./app/public"))))

	http.ListenAndServe(":"+os.Getenv("PORT"), mux)
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
