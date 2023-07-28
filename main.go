package main

import (
	"bytes"
	"encoding/json"

	"fmt"
	"log"
	"time"
	"os"
	"io"
	"net/http"
	"bufio"
	"strings"

	"github.com/dghubble/oauth1"
)

type BotKeys struct {
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
	bearerToken    string
}

func readBotKeysFromFile(filePath string) (*BotKeys, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	botKeys := &BotKeys{}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "consumerKey":
			botKeys.consumerKey = value
		case "consumerSecret":
			botKeys.consumerSecret = value
		case "accessToken":
			botKeys.accessToken = value
		case "accessSecret":
			botKeys.accessSecret = value
		case "bearerToken":
			botKeys.bearerToken = value
		}
	}

	return botKeys, nil
}

func oauth1Config(consumerKey, consumerSecret, accessToken, accessSecret string) *http.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	return httpClient
}

func postTweet(endpoint string, tweetJSON []byte, client *http.Client,Keys *BotKeys) ([]byte, error) {
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(tweetJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + Keys.bearerToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

func textToTweetJSON(tweetText string) []byte {
	tweetData := map[string]string{"text": tweetText}
	tweetJSON, err := json.Marshal(tweetData)
	if err != nil {
		log.Fatalf("Error marshaling tweet data: %v\n", err)
	}
	return tweetJSON
}

func writeResponseToFile(responseData []byte) error {
	// Create a file with the current timestamp as its name
	fileName := "responses/"+time.Now().Format("2006-01-02-15-04-05") + ".json"

	// Open the file in write-only mode, create if it doesn't exist
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the API response data to the file
	_, err = file.Write(responseData)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Iniciando programa...")

	botKeys, err := readBotKeysFromFile("./claves.txt")
	if err != nil {
		fmt.Println("Error reading bot keys:", err)
		return
	}

	config := oauth1Config(botKeys.consumerKey, botKeys.consumerSecret, botKeys.accessToken, botKeys.accessSecret)
	endpoint := "https://api.twitter.com/2/tweets"

	var tweetText string
	//tweetText = "Hola Twitter! Este es mi primer tweet automatizado de la cuenta 😎"
	
	fmt.Print("Ingrese el tweet que desea publicar: ")
	inputReader := bufio.NewReader(os.Stdin)
	tweetText,_ = inputReader.ReadString('\n')
	fmt.Printf("You entered: %s\n", tweetText)

	resp, err := postTweet(endpoint, textToTweetJSON(tweetText), config, botKeys)
	if err != nil {
		log.Fatalf("Error posting tweet: %v\n", err)
	}
	//Escribo la respuesta en un archivo
	err = writeResponseToFile(resp)
	if err != nil {
		log.Fatalf("Error writing response to file: %v\n", err)
	}

	fmt.Printf("Respuesta del servidor: %s\n", resp)
	fmt.Println("Finalizando ejecución, adiós!")
}