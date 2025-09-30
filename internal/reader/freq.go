package reader

import (
	"encoding/csv"
	"reader/data"
	"strconv"
)

type WordFreqRecord struct {
	Word     string
	Freq     int
	Position int
}

type WordFreq struct {
	Words   []WordFreqRecord
	WordMap map[string]*WordFreqRecord
}

func NewWordFreq() WordFreq {
	f, err := data.DataFiles.Open("unigram_freq.csv")
	if err != nil {
		panic(err)
	}
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	records = records[1:]
	wordFreq := WordFreq{
		Words:   make([]WordFreqRecord, 0),
		WordMap: make(map[string]*WordFreqRecord),
	}
	for i, record := range records {
		word := record[0]
		freq, err := strconv.Atoi(record[1])
		if err != nil {
			panic(err)
		}
		wordFreq.Words = append(wordFreq.Words, WordFreqRecord{
			Position: i + 1,
			Word:     word,
			Freq:     freq,
		})
	}
	for i, word := range wordFreq.Words {
		wordFreq.WordMap[word.Word] = &wordFreq.Words[i]
	}
	return wordFreq
}

func (wf *WordFreq) Get(word string) *WordFreqRecord {
	return wf.WordMap[word] 
}
