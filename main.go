package main

import (
	"embed"
	"encoding/gob"
	"io/fs"
	"log"
	"mime"
	"net/http"

	"github.com/hellosam123/pompeii_golang/internal/handlers"
	"github.com/hellosam123/pompeii_golang/internal/middleware"
	"github.com/hellosam123/pompeii_golang/internal/models"
)

//go:embed static
var staticFiles embed.FS

func main() {
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	gob.Register([]models.AnsweredVocabID{})

	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")

	pompeii := http.NewServeMux()

	mux := http.NewServeMux()
	// api routes
	mux.HandleFunc("/get_vocab", handlers.GetVocabHandler)
	mux.HandleFunc("/get_random_vocab", handlers.GetRandomVocabHandler)
	mux.HandleFunc("/check_answer", handlers.CheckAnswerHandler)
	mux.HandleFunc("/get_score", handlers.GetScoreHandler)
	mux.HandleFunc("/get_game_score", handlers.GetGameScoreHandler)

	// user endpoints
	mux.HandleFunc("/", handlers.IndexHandler)
	mux.HandleFunc("/vocabulary", handlers.VocabularyHandler)
	mux.HandleFunc("/game_settings", handlers.GameSettingsHandler)
	mux.HandleFunc("/load_game", handlers.LoadGameHandler)
	mux.HandleFunc("/classic", handlers.ClassicGameModeHandler)
	mux.HandleFunc("/game", handlers.NormalGameModeHandler)
	mux.HandleFunc("/game_over", handlers.GameOverHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	pompeii.Handle("/pompeii/", http.StripPrefix("/pompeii", mux))

	server := http.Server{
		Addr:    ":5030",
		Handler: middleware.Logging(pompeii),
	}

	log.Println("Server starting on http://localhost:5030/pompeii/")
	log.Fatal(server.ListenAndServe())
}
