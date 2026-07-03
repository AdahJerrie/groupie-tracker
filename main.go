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

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type LocationIndex struct {
	Index []Location `json:"index"`
}

type Dates struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}
type DatesIndex struct {
	Index []Dates `json:"index"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type RelationIndex struct {
	Index []Relation `json:"index"`
}

func main() {
	artistsDetails, err := FetchArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Println(err)
		return
	}

	locateDetails, err := FetchLocations("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Println(err)
		return
	}

	dates, err := FetchDates("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		log.Println(err)
		return
	}

	relations, err := FetchRelation("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Println(err)
		return
	}

	for i, date := range dates.Index {
		fmt.Printf("%d %s\n", i, date.Dates)
	}

	for i, relate := range relations.Index {
		fmt.Printf("%d %v\n", i, relate.DatesLocations)
	}

	for i, artist := range artistsDetails {
		fmt.Printf("%d %+v\n", i, artist)
	}

	for i, location := range locateDetails.Index {
		fmt.Printf("%d %s\n", i, location.Locations)
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

func FetchLocations(url string) (LocationIndex, error) {
	resp, err := http.Get(url)
	if err != nil {
		return LocationIndex{}, fmt.Errorf("fetching locations: %w", err)
	}
	defer resp.Body.Close()

	respByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationIndex{}, fmt.Errorf("reading response body: %w", err)
	}

	var locate LocationIndex
	if err := json.Unmarshal(respByte, &locate); err != nil {
		return LocationIndex{}, fmt.Errorf("decoding locate JSON: %w", err)
	}
	return locate, nil
}

func FetchRelation(url string) (RelationIndex, error) {
	resp, err := http.Get(url)
	if err != nil {
		return RelationIndex{}, fmt.Errorf("fetching relation: %w", err)
	}
	defer resp.Body.Close()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return RelationIndex{}, fmt.Errorf("reading response body: %w", err)
	}

	var relations RelationIndex
	if err := json.Unmarshal(byteBody, &relations); err != nil {
		return RelationIndex{}, fmt.Errorf("decoding relation JSON: %w", err)
	}
	return relations, nil
}

func FetchDates(url string) (DatesIndex, error) {
	resp, err := http.Get(url)
	if err != nil {
		return DatesIndex{}, fmt.Errorf("fetching dates: %w", err)
	}
	defer resp.Body.Close()

	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return DatesIndex{}, fmt.Errorf("reading response body: %w", err)
	}

	var dates DatesIndex
	if err := json.Unmarshal(bodyByte, &dates); err != nil {
		return DatesIndex{}, fmt.Errorf("decoding date JSON: %w", err)
	}
	return dates, nil
}
