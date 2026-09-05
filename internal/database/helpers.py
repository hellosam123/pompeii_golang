import sqlite3


def get_db_connection():
    conn = sqlite3.connect("database.db")
    conn.row_factory = sqlite3.Row  # converts row to dictionary object
    return conn


def get_vocab_id_from_vocab_word(vocab_word: str) -> int:
    conn = get_db_connection()

    vocab = conn.execute(
        """
        SELECT *
        FROM vocab
        WHERE vocab_word = ?
        """,
        (vocab_word,),
    ).fetchone()

    conn.close()

    if vocab:
        return vocab["vocab_id"]
    else:
        return 0


def new_vocab(
    vocab_word: str,
    shown_translations: list[str] = [],
    all_translations: list[str] = [],
    vocab_groups: list[str] = [],
):
    """Inserts a new record into the database
    Returns the new record"""

    conn = get_db_connection()

    # 💻  EDIT THIS ↓ EXECUTE STATEMENT TO MATCH YOUR DATABASE
    vocab_id = get_vocab_id_from_vocab_word(
        vocab_word
    )  # makes sure latin_vocab isn't duplicated
    if not vocab_id:
        vocab_id = conn.execute(
            """
            INSERT INTO
            vocab (vocab_word)
            VALUES (?)""",
            (vocab_word,),
        ).lastrowid

    if shown_translations:
        for shown_translation in shown_translations:
            conn.execute(
                """
                INSERT INTO
                shown_translations (vocab_id, english_translation)
                VALUES (?, ?)""",
                (vocab_id, shown_translation),
            )

    if all_translations:
        for all_translation in all_translations:
            conn.execute(
                """
                INSERT INTO
                all_translations (vocab_id, english_translation)
                VALUES (?, ?)""",
                (vocab_id, all_translation),
            )

    if vocab_groups:  # this doesn't account for duplicate groups, but that is accounted for with DISTINCT()
        for vocab_group in vocab_groups:
            conn.execute(
                """
                INSERT INTO
                vocab_groups (vocab_id, vocab_group)
                VALUES (?, ?)""",
                (vocab_id, vocab_group),
            )
    conn.commit()

    conn.close()
