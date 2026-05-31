-- Use this file to define your SQL table structure
-- 1) Replace TABLE_NAME 
-- 2) Add appropriate fields with the data type


CREATE TABLE IF NOT EXISTS vocab(
    vocab_id INTEGER PRIMARY KEY AUTOINCREMENT,
    vocab_word TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS shown_translations(
    translation_id INTEGER PRIMARY KEY AUTOINCREMENT,
    vocab_id INTEGER,
    english_translation TEXT,
    FOREIGN KEY (vocab_id) REFERENCES vocab(vocab_id)
);

CREATE TABLE IF NOT EXISTS all_translations(
    translation_id INTEGER PRIMARY KEY AUTOINCREMENT,
    vocab_id INTEGER,
    english_translation TEXT,
    FOREIGN KEY (vocab_id) REFERENCES vocab(vocab_id)
);

CREATE TABLE IF NOT EXISTS vocab_groups(
    group_id INTEGER PRIMARY KEY AUTOINCREMENT,
    vocab_id INTEGER,
    vocab_group TEXT,
    FOREIGN KEY (vocab_id) REFERENCES vocab(vocab_id)
);