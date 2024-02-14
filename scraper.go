package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
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

	// Take user input for artist
	artist := getUserInput(reader, "Enter the artist: ")

	// Construct the search URL with the user input
	searchURL := constructSearchURL(artist)

	songTitles := filterHTML(searchURL, artist)

	for _, song := range songTitles {
		// Construct the URL with user input
		url := ArtistSongToURL(artist, song)

		// Make the HTTP request and get the response
		response := makeRequest(url)

		// Parse the response to extract lyrics
		lyrics := parseLyrics(response)

		// Save lyrics to a file
		if lyrics==""{
			fmt.Printf("Lyrics not found in response. URL:%s\n",url)
			continue
		}
		saveLyricsToFile(artist, song, lyrics)
	}
}

func getUserInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
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
		return ""
	}

	// Extract text between the <Lyric> tags
	lyrics := response[startIndex+len(startTag) : endIndex]

	return lyrics
}

func saveLyricsToFile(artist, song, lyrics string) {
	song = strings.Title(strings.ToLower(song))
	re := regexp.MustCompile(`[?<>:"/\\|*]`)
	song = re.ReplaceAllString(song, "")
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

func constructSearchURL(artist string) string {
	// Replace spaces in the artist with '+'
	artist = strings.ReplaceAll(artist, " ", "+")
	return fmt.Sprintf("http://www.chartlyrics.com/search.aspx?q=%s", artist)
}

func printHTMLPage(urlString string) string {
	// Make an HTTP request to the URL
	response, err := http.Get(urlString)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer response.Body.Close()

	// Read the HTML page
	htmlPage, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	return string(htmlPage)
}

func filterHTML(urlString, artist string) []string {
	// Get the HTML content
	htmlContent := printHTMLPage(urlString)

	// Define a regular expression for the desired format
	pattern := fmt.Sprintf(`([^\/]+).aspx">%s`,artist)
	re := regexp.MustCompile(pattern)

	// Find all matches in the HTML content
	matches := re.FindAllStringSubmatch(htmlContent, -1)

	// Extract the matched text
	var songTitles []string
	for _, match := range matches {
		songTitle, err := url.QueryUnescape(match[1])
		if err != nil {
			fmt.Println("Error decoding URL:", err)
			return nil
		}
		songTitles = append(songTitles, songTitle)
	}

	return songTitles
}