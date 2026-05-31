# "I walk" -> ["I walk", "I am walking"]
# "I carry away" -> ["I carry away", "I am carrying away"]
# "I run" -> ["I run", "I am running"]
# "dog" -> ["dog"]

import re
import nltk
# nltk.download('punkt_tab')
# nltk.download('averaged_perceptron_tagger_eng')


def edge_case_verb_checker(word: str) -> bool:
    """Checks for edge cases where a word may be both a verb and noun."""
    if word == "order":
        return True
    if word == "cut":
        return True
    return False


def reverse_lemmatize(phrase: str) -> list:
    """
    If a word is a verb, creates a list including the word itself and its present continuous form.
    "I run" -> ["I run", "I am running"]
    """
    tokens = nltk.word_tokenize(phrase)
    pos_tags = nltk.pos_tag(tokens)

    verb_index = next((i for i, (word, tag) in enumerate(pos_tags) if (
        tag.startswith('VB') and word != "am") or edge_case_verb_checker(word)), None)

    if verb_index is not None:
        new_tokens = []
        for i, (word, tag) in enumerate(pos_tags):
            if i == verb_index:
                new_tokens.append("am")
                new_tokens.append(get_present_participle(word))
            else:
                new_tokens.append(word)

        continuous = " ".join(new_tokens)
        return [phrase, continuous]

    else:
        return [phrase]


def get_present_participle(verb: str) -> str:
    """Gets the present participle for any particular verb, including edge cases and exceptions."""
    verb = verb.lower().strip()

    # 1. THE MASTER LEXICON (Hard-coded exceptions)
    # These words defy standard rules due to syllable stress or origin.
    lexicon = {
        # Stress on FIRST syllable (No doubling despite CVC pattern)
        "visit": "visiting", "limit": "limiting", "edit": "editing",
        "offer": "offering", "profit": "profiting", "target": "targeting",
        "budget": "budgeting", "market": "marketing", "focus": "focusing",
        "ballot": "balloting", "benefit": "benefiting", "credit": "crediting",
        "exit": "exiting", "inherit": "inheriting", "orbit": "orbiting",

        # Hard 'C' sound (Requires a 'k')
        "panic": "panicking", "traffic": "trafficking", "mimic": "mimicking",
        "picnic": "picnicking", "frolic": "frolicking", "shellac": "shellacking",

        # Stress on SECOND syllable (Always double)
        "refer": "referring", "defer": "deferring", "infer": "inferring",
        "prefer": "preferring", "transfer": "transferring", "begin": "beginning",
        "forget": "forgetting", "submit": "submitting", "transmit": "transmitting",
        "occur": "occurring", "compel": "compelling", "control": "controlling",

        # True Irregulars & Greek/Latin roots
        "be": "being", "ski": "skiing", "taxi": "taxiing", "age": "aging",
        "dye": "dyeing", "singe": "singeing", "canoe": "canoeing",
        "eye": "eyeing", "hoe": "hoeing", "queue": "queueing"
    }

    if verb in lexicon:
        return lexicon[verb]

    # 2. THE -IE RULE (die -> dying)
    if verb.endswith("ie"):
        return verb[:-2] + "ying"

    # 3. THE SILENT -E RULE (dance -> dancing)
    # We don't drop the 'e' if it's 'ee', 'oe', or 'ye' (handled in lexicon/exceptions)
    if verb.endswith("e") and not (verb.endswith("ee") or verb.endswith("oe")):
        return verb[:-1] + "ing"

    # 4. THE CVC DOUBLING RULE (The core logic)
    # [Consonant][Vowel][Consonant]
    if re.search(r'[^aeiou][aeiou][bcdfghjklmnprstvz]$', verb):
        # Rule: Double for 1-syllable words (run, sit, hop)
        # Note: We exclude 'w', 'x', and 'y' (show -> showing)
        vowels = len(re.findall(r'[aeiouy]+', verb))
        if vowels == 1:
            if verb[-1] not in "wxy":
                return verb + verb[-1] + "ing"

        # Rule: Special British 'L' (Uncomment if you want British style)
        # if verb.endswith('l'): return verb + 'ling'

    # 5. DEFAULT FALLBACK
    return verb + "ing"
