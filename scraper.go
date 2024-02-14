package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Create "Lyrics" folder if it doesn't exist
	err := os.MkdirAll("lyrics", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating Lyrics folder:", err)
		return
	}

	// Create a bufio reader for user input
	reader := bufio.NewReader(os.Stdin)

	// Take user input for artist and song
	artist := getUserInput(reader, "Enter the artist: ")
	song := getUserInput(reader, "Enter the song: ")

	// Construct the URL with user input
	url := ArtistSongToURL(artist, song)

	// Make the HTTP request and get the response
	response := makeRequest(url)

	// Parse the response to extract lyrics
	lyrics := parseLyrics(response)

	// Save lyrics to a file
	if lyrics!=""{
		saveLyricsToFile(artist, song, lyrics)
	}
}

func getUserInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.Title(strings.ToLower(strings.TrimSpace(input)))
}

func ArtistSongToURL(artist, song string) string {
	url := fmt.Sprintf("http://api.chartlyrics.com/apiv1.asmx/SearchLyricDirect?artist=%s&song=%s", artist, song)
	return strings.TrimSpace(strings.ReplaceAll(url, " ", "%"))
}

func makeRequest(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Request failed with status:", res.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		os.Exit(1)
	}

	return string(body)
}

func parseLyrics(response string) string {
	startTag := "<Lyric>"
	endTag := "</Lyric>"

	startIndex := strings.Index(response, startTag)
	endIndex := strings.Index(response, endTag)

	if startIndex == -1 || endIndex == -1 {
		fmt.Println("Lyrics not found in response.")
		return ""
	}

	// Extract text between the <Lyric> tags
	lyrics := response[startIndex+len(startTag) : endIndex]

	return lyrics
}

func saveLyricsToFile(artist, song, lyrics string) {
	fileName := strings.ReplaceAll(fmt.Sprintf("lyrics/%s - %s.txt", artist, song)," ", "")

	// Create or overwrite the file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Write lyrics to the file
	_, err = file.WriteString(lyrics)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Lyrics saved to", fileName)
}
