package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"

	"github.com/hellosam123/project_webapp_pompeii_go/internal/helpers"
	"github.com/hellosam123/project_webapp_pompeii_go/internal/models"
	"github.com/hellosam123/project_webapp_pompeii_go/internal/templates"
)

var (
	key   = []byte("25471bfb56548f8c2ce1ead21e3aa6c2a6114afe30fdc35eb39c68aaac71dca3")
	store = sessions.NewCookieStore(key)
)

func renderTemplate(w http.ResponseWriter, tmplName string, data any) {
	w.Header().Set("Content-Type", "text/html")

	tmpl, err := templates.GetTemplates(tmplName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("template error: %v", err)
		return
	}
}

func respondJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

func VocabularyHandler(w http.ResponseWriter, r *http.Request) {
	allVocab, err := helpers.GetAllVocab()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, "vocabulary.html", allVocab)
}

func GameSettingsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "game_settings.html", nil)
}

func LoadGameHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")

	var vocabGroup string
	var classicMode bool
	if r.Method == http.MethodPost {
		vocabGroup = r.FormValue("vocab-list")
		classicMode = r.FormValue("classic-mode") == "true"
	} else {
		classicMode = true
	}

	gameSession.Values["vocabGroup"] = vocabGroup
	gameSession.Values["classicMode"] = classicMode
	gameSession.Values["answeredVocabSlice"] = []models.AnsweredVocabID{}

	if classicMode {
		if r.Method == http.MethodPost {
			totalQuestions, err := strconv.Atoi(r.FormValue("num-questions"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			vocabList, err := helpers.GetVocabByGroup(vocabGroup)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var vocabIDStack []int

			for _, value := range vocabList {
				vocabIDStack = append(vocabIDStack, value.VocabID)
			}

			rand.Shuffle(len(vocabIDStack), func(i, j int) {
				vocabIDStack[i], vocabIDStack[j] = vocabIDStack[j], vocabIDStack[i]
			})

			if totalQuestions <= len(vocabIDStack) {
				vocabIDStack = vocabIDStack[:totalQuestions]
			}

			gameSession.Values["vocabIDStack"] = vocabIDStack
		}

		gameSession.Values["numQuestions"] = 0
		gameSession.Values["numCorrect"] = 0

	} else {
		difficultyMode := r.FormValue("difficulty-mode")
		gameSession.Values["difficultyMode"] = difficultyMode
		gameSession.Values["gameScore"] = 0
		var correctScore int
		switch difficultyMode {
		case "easy":
			correctScore = 50
		case "medium":
			correctScore = 100
		case "hard":
			correctScore = 200
		}

		gameSession.Values["correctScore"] = correctScore

	}

	err := gameSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if classicMode {
		http.Redirect(w, r, "/pompeii/classic", http.StatusFound)
	} else {
		http.Redirect(w, r, "/pompeii/game", http.StatusFound)
	}
}

func ClassicGameModeHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")

	if gameSession.Values["vocabIDStack"] != nil {
		renderTemplate(w, "classic.html", nil)
	} else {
		http.Redirect(w, r, "/pompeii/", http.StatusFound)
	}
}

func NormalGameModeHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")
	vocabGroup, ok := gameSession.Values["vocabGroup"].(string)
	if !ok {
		log.Printf("failed to get vocabGroup value from gameSession")
		return
	}
	difficultyMode, ok := gameSession.Values["difficultyMode"].(string)
	if !ok {
		log.Printf("failed to get difficultyMode value from gameSession")
		return
	}

	if vocabGroup != "" && difficultyMode != "" {
		data := models.GameData{
			VocabGroup:     vocabGroup,
			DifficultyMode: difficultyMode,
		}

		renderTemplate(w, "game.html", data)
	} else {
		http.Redirect(w, r, "/pompeii/", http.StatusFound)
	}
}

func GameOverHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")

	answeredVocabSlice, ok := gameSession.Values["answeredVocabSlice"].([]models.AnsweredVocabID)
	if !ok {
		answeredVocabSlice = []models.AnsweredVocabID{}
	}

	if answeredVocabSlice == nil {
		http.Redirect(w, r, "/pompeii/", http.StatusFound)
	}

	var vocabIDStack []int
	var answeredVocabData []models.AnsweredVocab

	for _, value := range answeredVocabSlice {
		vocabID, correct := value.VocabID, value.Correct

		if !correct {
			vocabIDStack = append(vocabIDStack, vocabID)
		}

		vocab, err := helpers.GetVocabByID(vocabID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		answeredVocab := models.AnsweredVocab{
			Vocab:   *vocab,
			Correct: correct,
		}
		answeredVocabData = append(answeredVocabData, answeredVocab)
	}

	retryIncorrect := vocabIDStack != nil
	if retryIncorrect {
		rand.Shuffle(len(vocabIDStack), func(i, j int) {
			vocabIDStack[i], vocabIDStack[j] = vocabIDStack[j], vocabIDStack[i]
		})

		gameSession.Values["vocabIDStack"] = vocabIDStack
	}

	var score string
	if gameSession.Values["classicMode"].(bool) {
		score = fmt.Sprintf("%d / %d", gameSession.Values["numCorrect"].(int), gameSession.Values["numQuestions"].(int))
	} else {
		score = strconv.Itoa(gameSession.Values["gameScore"].(int))
	}

	gameOverData := models.GameOverData{
		AnsweredVocabData: answeredVocabData,
		Score:             score,
		RetryIncorrect:    retryIncorrect,
	}

	err := gameSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "game_over.html", gameOverData)
}

func GetVocabHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")
	vocabIDStack, ok := gameSession.Values["vocabIDStack"].([]int)
	if !ok {
		vocabIDStack = []int{}
	}

	if len(vocabIDStack) == 0 {
		respondJSON(w, map[string]bool{"gameOver": true})
		return
	}

	n := len(vocabIDStack) - 1
	vocabID := vocabIDStack[n]
	vocabIDStack = vocabIDStack[:n]

	gameSession.Values["vocabIDStack"] = vocabIDStack

	vocab, err := helpers.GetVocabByID(vocabID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = gameSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]models.Vocab{"vocab": *vocab})
}

func GetRandomVocabHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")
	vocabGroup, ok := gameSession.Values["vocabGroup"].(string)
	if !ok {
		log.Printf("failed to get vocabGroup value from gameSession")
		return
	}

	vocab, err := helpers.GetRandomVocabByGroup(vocabGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]models.Vocab{"vocab": *vocab})
}

func CheckAnswerHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")
	classicMode, ok := gameSession.Values["classicMode"].(bool)
	if !ok {
		log.Printf("failed to get classicMode value from gameSession")
		return
	}

	checkAnswerQuery := r.URL.Query()
	vocabID, err := strconv.Atoi(checkAnswerQuery.Get("vocab_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	answer := checkAnswerQuery.Get("answer")

	isCorrect, err := helpers.CheckAnswer(vocabID, answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if classicMode {
		numQuestions, ok := gameSession.Values["numQuestions"].(int)
		if !ok {
			log.Printf("failed to get numQuestions value from gameSession")
			return
		}

		numCorrect, ok := gameSession.Values["numCorrect"].(int)
		if !ok {
			log.Printf("failed to get numCorrect value from gameSession")
			return
		}

		numQuestions++
		if isCorrect {
			numCorrect++
		}

		gameSession.Values["numQuestions"] = numQuestions
		gameSession.Values["numCorrect"] = numCorrect
	} else {
		if isCorrect {
			gameScore, ok := gameSession.Values["gameScore"].(int)
			if !ok {
				log.Printf("failed to get gameScore value from gameSession")
				return
			}

			correctScore, ok := gameSession.Values["correctScore"].(int)
			if !ok {
				log.Printf("failed to get correctScore value from gameSession")
				return
			}

			gameScore += correctScore
			gameSession.Values["gameScore"] = gameScore

		}
	}

	answeredVocabSlice, ok := gameSession.Values["answeredVocabSlice"].([]models.AnsweredVocabID)
	if !ok {
		log.Printf("failed to get answeredVocabSlice value from gameSession")
		return
	}

	answeredVocab := models.AnsweredVocabID{
		VocabID: vocabID,
		Correct: isCorrect,
	}

	answeredVocabSlice = append(answeredVocabSlice, answeredVocab)
	gameSession.Values["answeredVocabSlice"] = answeredVocabSlice

	err = gameSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]bool{"isCorrect": isCorrect})
}

func GetScoreHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")
	numQuestions, ok := gameSession.Values["numQuestions"].(int)
	if !ok {
		log.Printf("failed to get numQuestions value from gameSession")
		return
	}

	numCorrect, ok := gameSession.Values["numCorrect"].(int)
	if !ok {
		log.Printf("failed to get numCorrect value from gameSession")
		return
	}

	respondJSON(w, map[string]int{"numQuestions": numQuestions, "numCorrect": numCorrect})
}

func GetGameScoreHandler(w http.ResponseWriter, r *http.Request) {
	gameSession, _ := store.Get(r, "game-session")
	gameScore, ok := gameSession.Values["gameScore"].(int)
	if !ok {
		log.Printf("failed to get gameScore value from gameSession")
		return
	}

	respondJSON(w, map[string]int{"gameScore": gameScore})
}
