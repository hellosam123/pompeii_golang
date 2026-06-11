package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hellosam123/pompeii_golang/internal/handlers"
	"github.com/hellosam123/pompeii_golang/internal/middleware"
	"github.com/hellosam123/pompeii_golang/internal/models"
)

func main() {
	gob.Register([]models.AnsweredVocabID{})

	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "application/javascript")

	staticPath, err := getStaticPath()
	if err != nil {
		log.Fatal(err)
	}

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

	pompeii.Handle("/pompeii/", http.StripPrefix("/pompeii", mux))
	pompeii.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))

	server := http.Server{
		Addr:    ":5030",
		Handler: middleware.Logging(pompeii),
	}

	log.Println("Server starting on http://localhost:5030/pompeii/")
	log.Fatal(server.ListenAndServe())
}

func getStaticPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get current file path")
	}

	exeDir := filepath.Dir(exe)
	staticPath := filepath.Join(exeDir, "static")

	return staticPath, nil
}
