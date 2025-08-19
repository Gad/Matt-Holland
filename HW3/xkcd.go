package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Xkcd struct {
	Num              int `json:"num"`
	Url              string
	Year             string `json:"year"`
	Month            string `json:"month"`
	Day              string `json:"day"`
	Title            string `json:"safe_title"`
	TitleFields      []string
	Transcript       string `json:"transcript"`
	TranscriptFields []string
}

const MAX int = 100

func processUrl(url string) (Xkcd, error) {
	joke, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP protocol Error while reading joke : %v\n", err)
		os.Exit(-1)
	}

	defer joke.Body.Close()

	if joke.StatusCode != http.StatusOK {

		return Xkcd{}, fmt.Errorf("HTTP Status Error while reading joke : %v\n", joke.Status)
	} else {
		var t []byte
		t, err := io.ReadAll(joke.Body)

		if err != nil {
			return Xkcd{}, fmt.Errorf("Error reading body of %s : %v\n", url, err)
		}
		var j = Xkcd{Url: url}
		if err := json.Unmarshal(t, &j); err != nil {

			return Xkcd{}, fmt.Errorf("Error decoding json from %s : %v\n", url, err)
		}

		return j, nil
	}
}

func buildCollection(jokesCollection *[]Xkcd) {

	for i := 1; i < 1+MAX; i++ {

		x, err := processUrl("https://xkcd.com/" + strconv.Itoa(i+1) + "/info.0.json")

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error from processUrl : %v", err)
			break
		}

		*jokesCollection = append(*jokesCollection, x)

	}
}

func cleanTranscripts(jokesCollection *[]Xkcd) (*[]Xkcd, error) {

	cleanJokes := make([]Xkcd, MAX)
	copy(cleanJokes, *jokesCollection)
	re, err := regexp.Compile(`[^a-zA-Z0-9]|Alt:|Alt-title:`)
	if err != nil {
		return nil, fmt.Errorf("Error compiling regexp for cleaning transcript : %v\n", err)
	}

	for i := range cleanJokes {
		cleanJokes[i].Transcript = re.ReplaceAllString(cleanJokes[i].Transcript, " ")
		cleanJokes[i].Title = re.ReplaceAllString(cleanJokes[i].Title, " ")
		cleanJokes[i].Transcript = strings.ToLower(cleanJokes[i].Transcript)
		cleanJokes[i].Title = strings.ToLower(cleanJokes[i].Title)

	}

	return &cleanJokes, nil
}

func splitTranscript(jokesCollection *[]Xkcd) {

	for i := range *jokesCollection {
		(*jokesCollection)[i].TranscriptFields = strings.Fields((*jokesCollection)[i].Transcript)
	}

}

func splitTitle(jokesCollection *[]Xkcd) {

	for i := range *jokesCollection {
		(*jokesCollection)[i].TitleFields = strings.Fields((*jokesCollection)[i].Title)
	}

}

func createIndex(JokesCollection *[]Xkcd) map[string][]int {

	// for each key (a word), provides the list of joke numbers where the word is present in the
	// joke transcript or title
	var index = make(map[string][]int)

	for _, joke := range *JokesCollection {
		
		growIndex := func(fields []string, index *map[string][]int) {
			for _, word := range fields {

				s, exists := (*index)[word]
				if exists {
					// check the joke number is already in the slice
					if slices.Contains(s, joke.Num) {
						continue
					}
					(*index)[word] = append((*index)[word], joke.Num)
					continue
				}
				(*index)[word] = make([]int, 0)
				(*index)[word] = append((*index)[word], joke.Num)

			}
		}

		growIndex(joke.TranscriptFields, &index)
		growIndex(joke.TitleFields, &index)

	}
	return index
}

func indexSearch(searchTerm string, index map[string][]int) []int {

	val, exists := index[searchTerm]
	if exists {
		return val
	}
	return nil

}

func printResults(searchTerm string, searchResults []int, jokesCollection []Xkcd) {
	fmt.Printf("Search for \"%s\" returns :\n", searchTerm)

	for _, num := range searchResults {
		for j := range jokesCollection {
			if jokesCollection[j].Num == num {
				fmt.Printf("https://xkcd.com/%d/ %s/%s/%s \"%s\"\n",
					num,
					jokesCollection[j].Day,
					jokesCollection[j].Month,
					jokesCollection[j].Year,
					jokesCollection[j].Title)
			}
		}
	}
}

func main() {

	jokesCollection := make([]Xkcd, 0, MAX)

	buildCollection(&jokesCollection)
	log.Printf("Done with retrieving %d jokes", MAX)

	cleanJokesCollection, err := cleanTranscripts(&jokesCollection)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from cleanTranscripts : %v", err)
		os.Exit(-1)
	}
	log.Println("Done with cleaning jokes entries")

	splitTranscript(cleanJokesCollection)
	splitTitle(cleanJokesCollection)
	log.Println("Done with analysing jokes titles and transcripts")

	index := createIndex(cleanJokesCollection)
	log.Println("Done with creating index")



	if l := len(os.Args); l != 1 {
		for i := 1; i < l; i++ {
			searchTerm := os.Args[i]
			r := indexSearch(searchTerm, index)
			if r != nil {
				printResults(searchTerm, r, jokesCollection)
			}
		}

	}

}
