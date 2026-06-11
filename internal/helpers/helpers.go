package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/hellosam123/pompeii_golang/internal/models"

	_ "modernc.org/sqlite"
)

var initOnce sync.Once
var readOnlyDB *sql.DB
var dbInitError error

func getDBPath() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get current file path")
	}

	exeDir := filepath.Dir(exe)
	dbPath := filepath.Join(exeDir, "internal", "database", "database.db")
	dbPath = filepath.Clean(dbPath)

	return dbPath, nil
}

func initReadOnlyDB() {
	dbPath, err := getDBPath()

	if err != nil {
		dbInitError = fmt.Errorf("failed to get database path: %w", err)
		return
	}

	// since DB doesn't change, cache=shared speeds up queries
	readOnlyDBPath := fmt.Sprintf("file:%s?mode=ro&cache=shared&immutable=1", dbPath)

	readOnlyDB, err = sql.Open("sqlite", readOnlyDBPath)

	if err != nil {
		dbInitError = fmt.Errorf("failed to open database: %w", err)
		return
	}

	if err := readOnlyDB.Ping(); err != nil {
		readOnlyDB.Close()
		dbInitError = fmt.Errorf("failed to connect to database: %w", err)
		return
	}

	readOnlyDB.SetMaxOpenConns(25)
	readOnlyDB.SetMaxIdleConns(10)
	readOnlyDB.SetConnMaxLifetime(5 * time.Minute)
	readOnlyDB.SetConnMaxIdleTime(1 * time.Minute)
}

func getReadOnlyDB() (*sql.DB, error) {
	initOnce.Do(initReadOnlyDB)
	if dbInitError != nil {
		return nil, dbInitError
	}
	return readOnlyDB, nil
}

func sqlRowsToVocabs(rows *sql.Rows) ([]models.Vocab, error) {

	defer rows.Close()

	var vocabs []models.Vocab
	for rows.Next() {
		var vocab models.Vocab

		err := rows.Scan(&vocab.VocabID, &vocab.VocabWord, &vocab.EnglishTranslation, &vocab.VocabGroup)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		vocabs = append(vocabs, vocab)
	}

	return vocabs, nil
}

func GetAllVocab() ([]models.Vocab, error) {
	db, err := getReadOnlyDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`
		SELECT vocab.vocab_id, vocab.vocab_word,
        GROUP_CONCAT(DISTINCT english_translation) AS english_translation,
        GROUP_CONCAT(DISTINCT vocab_group) AS vocab_group
        FROM vocab
        LEFT JOIN shown_translations ON vocab.vocab_id = shown_translations.vocab_id
        LEFT JOIN vocab_groups ON vocab.vocab_id = vocab_groups.vocab_id
        GROUP BY vocab.vocab_id, vocab.vocab_word
        ORDER BY vocab.vocab_word
		`)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return sqlRowsToVocabs(rows)
}

func GetVocabByGroup(vocabGroup string) ([]models.Vocab, error) {
	db, err := getReadOnlyDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`
        SELECT vocab.vocab_id, vocab.vocab_word,
        GROUP_CONCAT(DISTINCT english_translation) AS english_translation,
        GROUP_CONCAT(DISTINCT vocab_group) AS vocab_group
        FROM vocab
        LEFT JOIN shown_translations ON vocab.vocab_id = shown_translations.vocab_id
        LEFT JOIN vocab_groups ON vocab.vocab_id = vocab_groups.vocab_id
        WHERE vocab_group = ?
        GROUP BY vocab.vocab_id, vocab.vocab_word
        ORDER BY vocab.vocab_word
		`, vocabGroup)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	vocabs, err := sqlRowsToVocabs(rows)

	if err != nil {
		return nil, err
	}

	if len(vocabs) == 0 {
		return nil, fmt.Errorf("no vocab found for group: %s", vocabGroup)
	}

	return vocabs, nil
}

func GetVocabByID(vocabID int) (*models.Vocab, error) {
	db, err := getReadOnlyDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`
        SELECT vocab.vocab_id, vocab.vocab_word,
        GROUP_CONCAT(DISTINCT english_translation) AS english_translation,
        GROUP_CONCAT(DISTINCT vocab_group) AS vocab_group
        FROM vocab
        LEFT JOIN shown_translations ON vocab.vocab_id = shown_translations.vocab_id
        LEFT JOIN vocab_groups ON vocab.vocab_id = vocab_groups.vocab_id
        WHERE vocab.vocab_id = ?
        GROUP BY vocab.vocab_id, vocab.vocab_word
        ORDER BY vocab.vocab_word
		`, vocabID)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	vocabs, err := sqlRowsToVocabs(rows)

	if err != nil {
		return nil, err
	}

	if len(vocabs) == 0 {
		return nil, fmt.Errorf("vocab with ID %d not found", vocabID)
	}

	return &vocabs[0], nil
}

func GetVocabByIDAllTranslations(vocabID int) (*models.Vocab, error) {
	db, err := getReadOnlyDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`
        SELECT vocab.vocab_id, vocab.vocab_word,
        GROUP_CONCAT(DISTINCT english_translation) AS english_translation,
        GROUP_CONCAT(DISTINCT vocab_group) AS vocab_group
        FROM vocab
        LEFT JOIN all_translations ON vocab.vocab_id = all_translations.vocab_id
        LEFT JOIN vocab_groups ON vocab.vocab_id = vocab_groups.vocab_id
        WHERE vocab.vocab_id = ?
        GROUP BY vocab.vocab_id, vocab.vocab_word
        ORDER BY vocab.vocab_word
		`, vocabID)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	vocabs, err := sqlRowsToVocabs(rows)

	if err != nil {
		return nil, err
	}

	if len(vocabs) == 0 {
		return nil, fmt.Errorf("vocab with ID %d not found", vocabID)
	}

	return &vocabs[0], nil
}

func GetRandomVocabByGroup(vocabGroup string) (*models.Vocab, error) {
	db, err := getReadOnlyDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`
        SELECT vocab.vocab_id, vocab.vocab_word,
        GROUP_CONCAT(DISTINCT english_translation) AS english_translation,
        GROUP_CONCAT(DISTINCT vocab_group) AS vocab_group
        FROM vocab
        LEFT JOIN shown_translations ON vocab.vocab_id = shown_translations.vocab_id
        LEFT JOIN vocab_groups ON vocab.vocab_id = vocab_groups.vocab_id
        WHERE vocab_group = ?
        GROUP BY vocab.vocab_id, vocab.vocab_word
        ORDER BY RANDOM()
        LIMIT 1
		`, vocabGroup)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	vocabs, err := sqlRowsToVocabs(rows)

	if err != nil {
		return nil, err
	}

	if len(vocabs) == 0 {
		return nil, fmt.Errorf("no vocab found for group: %s", vocabGroup)
	}

	return &vocabs[0], nil
}

func cleanAnswer(answer string) string {
	re := regexp.MustCompile(`\(.*?\)`)
	cleaned := re.ReplaceAllString(answer, "")
	re = regexp.MustCompile(`[^\w]`)
	cleaned = re.ReplaceAllString(cleaned, "")
	cleaned = strings.ToLower(cleaned)
	return cleaned
}

func CheckAnswer(vocabID int, answer string) (bool, error) {
	vocab, err := GetVocabByIDAllTranslations(vocabID)
	if err != nil {
		return false, err
	}

	correctAnswers := strings.Split(vocab.EnglishTranslation, ",")

	answer = cleanAnswer(answer)
	for _, value := range correctAnswers {
		correctAnswer := cleanAnswer(value)
		if answer == correctAnswer {
			return true, nil
		}
	}

	return false, nil
}
