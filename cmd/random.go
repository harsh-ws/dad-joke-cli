/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Returns a random dadjoke.",
	Long:  `Returns a random dadjoke.`,
	Run: func(cmd *cobra.Command, args []string) {
		JokeKeyword, _ := cmd.Flags().GetString("keyword")
		if JokeKeyword != "" {
			getRandomJokeWithKeyword(JokeKeyword)
		} else {
			getRandomJoke()
		}
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)

	randomCmd.PersistentFlags().String("keyword", "", "Dadjoke for a specific word.")
}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
}

type SearchResponse struct {
	Results       json.RawMessage `json:"results"`
	SearchKeyword string          `json:"search_term"`
	Status        int             `json:"status"`
	TotalJokes    int             `json:"total_jokes"`
}

func getRandomJoke() {
	url := "https://icanhazdadjoke.com/"
	responseBytes := getJokeData(url)
	joke := Joke{}
	err := json.Unmarshal(responseBytes, &joke)
	if err != nil {
		log.Printf("Unable to do unmarshall of json data - %v", err)
	}
	var randomJoke = joke.Joke
	fmt.Println(randomJoke)
}

func getJokeData(baseAPI string) []byte {
	request, err := http.NewRequest(http.MethodGet, baseAPI, nil)
	if err != nil {
		log.Printf("Unable to call the request %v", err)
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "DADJOKE CLI : https://github.com/harsh-ws/dad-joke-cli")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Unable to call the request %v", err)
	}
	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Unable to read the response body %v", err)
	}
	return responseBytes
}
func getJokeDataWithKeyword(jokeKeyWord string) (totalJokes int, jokeList []Joke) {

	url := fmt.Sprintf("https://icanhazdadjoke.com/search?term=%s", jokeKeyWord)
	responseBytes := getJokeData(url)
	jokeListRaw := SearchResponse{}

	err := json.Unmarshal(responseBytes, &jokeListRaw)
	if err != nil {
		log.Printf("Unable to do unmarshall of json data - %v", err)
	}

	jokes := []Joke{}

	err1 := json.Unmarshal(jokeListRaw.Results, &jokes)
	if err != nil {
		log.Printf("Unable to unmarshal the responses. -%v", err1)
	}
	return jokeListRaw.TotalJokes, jokes

}

func getRandomJokeWithKeyword(jokeKeyword string) {
	total, results := getJokeDataWithKeyword(jokeKeyword)
	randomiser(total, results)
	//fmt.Println(results)
}

func randomiser(length int, jokeList []Joke) {
	rand.Seed(time.Now().Unix())
	min := 0
	max := length - 1

	if length <= 0 {
		err := fmt.Errorf("unable to find a dadjoke with this keyword")
		fmt.Println(err)
	} else {
		randomNumber := min + rand.Intn(max-min)
		fmt.Println(jokeList[randomNumber].Joke)
	}

}
