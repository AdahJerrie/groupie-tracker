package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	rawJSON := `{"id":1,"image":"https://example.com/queen.jpeg","name":"Queen","members":["Freddie Mercury","Brian May","Roger Taylor","John Deacon"],"creationDate":1970,"firstAlbum":"14-12-1973"}`
	var artist Artist
	if err := json.Unmarshal([]byte(rawJSON), &artist); err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%+v\n", artist)
}
