import helpers
import pandas
import time
import re
import internal.database.reverse_lemmatize as reverse_lemmatize


def import_data(csv_file, vocab_group):
    start_time = time.perf_counter()  # timer to get runtime
    imported_vocab = 0
    csv = csv_file
    igcse_df = pandas.read_csv(csv)
    igcse_vocab_list_of_dicts = igcse_df.to_dict(orient='records')

    for vocab_dict in igcse_vocab_list_of_dicts:
        vocab_word = vocab_dict['vocab_word']
        shown_english_translation = vocab_dict['english_translation']

        shown_english_translation_list = re.split(r',(?![^(]*\))', shown_english_translation)
        stripped_shown_english_translation_list = []

        for word in shown_english_translation_list:
            stripped_shown_english_translation_list.append(word.strip())

        all_english_translation_list = []
        for word in stripped_shown_english_translation_list:  # absolute hell on earth
            reverse_lemmatized_word_list = reverse_lemmatize.reverse_lemmatize(word)
            all_english_translation_list.extend(reverse_lemmatized_word_list)

        helpers.new_vocab(vocab_word.strip(), shown_translations=stripped_shown_english_translation_list,
                          all_translations=all_english_translation_list, vocab_groups=[vocab_group])

        imported_vocab += 1

    end_time = time.perf_counter()
    elapsed_time = end_time - start_time
    print(f"{imported_vocab} vocab imported from {csv_file} in {elapsed_time:.4f} seconds")


import_data('assets/vocab_list_igcse_full.csv', 'Latin IGCSE')
import_data('assets/vocab_list_gcse_dvl.csv', 'Latin GCSE DVL')
