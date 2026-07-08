package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
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

var (
	artists   []Artist
	locations LocationIndex
	dates     DatesIndex
	relations RelationIndex
)

type ArtistPageData struct {
	ID            int
	Artist_info   Artist
	Date_Location map[string][]string
}

var tmpl *template.Template

func main() {

	var err error
	tmpl = template.Must(template.ParseGlob("templates/*.html"))

	artists, err = FetchArtists("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}

	locations, err = FetchLocations("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatal(err)
	}

	dates, err = FetchDates("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		log.Fatal(err)
	}

	relations, err = FetchRelation("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/artists", artistsHandler)
	mux.HandleFunc("/artist", artistHandler)

	log.Println("starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func FetchArtists(url string) ([]Artist, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching artists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status fetching artists: %s", resp.Status)
	}

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

	if resp.StatusCode != http.StatusOK {
		return LocationIndex{}, fmt.Errorf("unexpected status fetching locations: %s", resp.Status)
	}

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

	if resp.StatusCode != http.StatusOK {
		return RelationIndex{}, fmt.Errorf("unexpected status fetching relation: %s", resp.Status)
	}

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

	if resp.StatusCode != http.StatusOK {
		return DatesIndex{}, fmt.Errorf("unexpected status fetching date: %s", resp.Status)
	}

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

func homeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := tmpl.ExecuteTemplate(w, "index.html", artists); err != nil {
		log.Printf("executing template: %v", err)
	}
}

func artistsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	fmt.Fprintln(w, "Artists page")

	for i, artist := range artists {
		fmt.Fprintf(w, "%d %v\n", i, artist.Name)
	}
}

func artistHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path != "/artist" {
		http.NotFound(w, req)
		return
	}

	id := req.URL.Query().Get("id")
	id_int, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "error converting id", http.StatusBadRequest)
		return
	}
	data := ArtistPageData{
		ID: id_int,
	}

	if err := tmpl.ExecuteTemplate(w, "artist.html", data); err != nil {
		log.Printf("executing template: %v", err)
	}
}
