package models

type Vocab struct {
	VocabID            int
	VocabWord          string
	EnglishTranslation string
	VocabGroup         string
}

type AnsweredVocab struct {
	Vocab
	Correct bool
}

type AnsweredVocabID struct {
	VocabID int
	Correct bool
}

type GameData struct {
	VocabGroup     string
	DifficultyMode string
}

type GameOverData struct {
	AnsweredVocabData []AnsweredVocab
	Score             string
	RetryIncorrect    bool
}
