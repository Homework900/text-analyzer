package main

import (
	"flag"
	"fmt"
	"os"
)

type WordStat struct {
	Word  string
	Count int
}

type TextStats struct {
	LineCount int
	WordCount int
	CharCount int
	wordFreq  map[string]int
}

func main() {

	topN := flag.Int("top", 5, "сколько слов выводить в топе")
	linesOnly := flag.Bool("l", false, "выводить только количество строк")

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("Ошибка: не указан файл для анализа.")
		return
	}
	fileName := args[0]

	stats, err := analyzeFile(fileName)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	if *linesOnly {

		fmt.Printf("%d\n", stats.LineCount)
		return
	}

	printStats(stats)

	topWords := getTopWords(stats.wordFreq, *topN)

	fmt.Printf("\n--- Топ-%d слов ---\n", *topN)
	for i, ws := range topWords {
		fmt.Printf("%d. %s: %d\n", i+1, ws.Word, ws.Count)
	}
}

func analyzeFile(fileName string) (*TextStats, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	stats := &TextStats{
		wordFreq: make(map[string]int),
	}

	var currentWord []byte

	processWord := func(wordBytes []byte) {
		if len(wordBytes) == 0 {
			return
		}
		word := string(wordBytes)

		lowerWord := make([]byte, len(word))
		for i := 0; i < len(word); i++ {
			c := word[i]
			if c >= 'A' && c <= 'Z' {
				c = c + ('a' - 'A')
			}
			lowerWord[i] = c
		}

		cleanWord := string(lowerWord)

		stats.WordCount++
		stats.wordFreq[cleanWord]++
	}

	for _, b := range data {
		if b == '\n' {
			stats.LineCount++
			stats.CharCount++
			processWord(currentWord)
			currentWord = currentWord[:0]
		} else if b == ' ' || b == '\t' || b == '\r' {
			stats.CharCount++
			processWord(currentWord)
			currentWord = currentWord[:0]
		} else {
			stats.CharCount++
			currentWord = append(currentWord, b)
		}
	}
	processWord(currentWord)

	if len(data) > 0 && data[len(data)-1] != '\n' {
		stats.LineCount++
	}

	return stats, nil
}

func getTopWords(freq map[string]int, topN int) []WordStat {
	var result []WordStat

	for word, count := range freq {
		result = append(result, WordStat{Word: word, Count: count})
	}
	n := len(result)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if result[j].Count < result[j+1].Count {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}

	if len(result) > topN {
		return result[:topN]
	}
	return result
}

func printStats(stats *TextStats) {
	fmt.Println("Статистика текста:")
	fmt.Printf("Строк: %d\n", stats.LineCount)
	fmt.Printf("Слов: %d\n", stats.WordCount)
	fmt.Printf("Символов: %d\n", stats.CharCount)
}