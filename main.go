package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

func main() {
	artistsDetails, err := FetchArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println(err)
	}

	for i, artist := range artistsDetails {
		fmt.Printf("%d %+v\n", i, artist.Name)
	}
}

func FetchArtists(url string) ([]Artist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching artists: %w", err)
	}
	defer resp.Body.Close()

	byteData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var artists []Artist
	if err := json.Unmarshal(byteData, &artists); err != nil {
		return nil, fmt.Errorf("decoding artists JSON: %w", err)
	}

	return artists, nil
}
