/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

// randomCmd represents the random command
var randomCmd = &cobra.Command{
	Use:   "random",
	Short: "Returns a random dadjoke.",
	Long:  `Returns a random dadjoke.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(randomCmd)
	getRandomJoke()
}

type Joke struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int    `json:"status"`
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
