from pdfminer.high_level import extract_text_to_fp
from io import StringIO
from csv import reader as csv_reader
from csv import DictWriter as csv_writer
from json import load as json_load
from collections import defaultdict

from text_preprocess import get_cleared_words_from_text

PATH_TO_CSV = "../../connector/output/records.csv"
PATH_TO_PDF = "../../connector/output/VKR/pdfs/"
PATH_TO_JSON = "../../connector/output/VKR/jsons/"


def convert_pdf_to_text(path):
    try:
        output = StringIO()

        with open(path, "rb") as f:
            extract_text_to_fp(f, output)

        return output.getvalue()
    except:
        return None


def get_cleared_pdf(path):
    text = convert_pdf_to_text(path)
    if text is not None:
        return get_cleared_words_from_text(text)
    else:
        return None


def get_keywords(path):
    with open(path, "rb") as f:
        json_dict = json_load(f)

    keywords = json_dict["keyWordsRu"]

    return keywords


def get_cleared_keywords(keywords):
    return [get_cleared_words_from_text(keyword) for keyword in keywords]


def match_lists(list1, list2):
    return set(list1) == set(list2)


def get_keywords_freq(words, keywords):
    keyword_sets = [set(keyword) for keyword in keywords]
    output_dict = defaultdict(int)

    for i in range(len(words)):
        for j in range(len(keyword_sets)):
            if words[i] in keyword_sets[j]:
                if set(words[i:i + len(keyword_sets[j])]) == keyword_sets[j]:
                    output_dict[j] += 1

    return output_dict


def main():
    csv_input_file = open(PATH_TO_CSV, 'r')
    reader = csv_reader(csv_input_file, delimiter=',')

    csv_output_file = open("out.csv", 'w', encoding='utf-8-sig', newline='')
    fieldnames = ['keyword', 'freq']
    writer = csv_writer(csv_output_file, fieldnames=fieldnames)
    writer.writeheader()

    k = 0
    for row in reader:
        print(str(k) + " : ", row[0])

        pdf_path = PATH_TO_PDF + row[0] + ".pdf"
        json_path = PATH_TO_JSON + row[0] + ".json"

        keywords = get_keywords(json_path)
        cleared_keywords = get_cleared_keywords(keywords)
        pdf_words = get_cleared_pdf(pdf_path)

        if pdf_words is None:
            continue

        keywords_freq = get_keywords_freq(pdf_words, cleared_keywords)

        for i in range(len(keywords)):
            writer.writerow({'keyword': keywords[i], 'freq': keywords_freq[i]})

        k = k + 1

    csv_input_file.close()
    csv_output_file.close()

if __name__ == "__main__":
    main()
