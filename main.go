package main

import "fmt"

type Artist struct {
	Name         string
	CreationDate int
	Members      []string
}

func main() {
	var artists []Artist

	artist1 := Artist{
		Name:         "Abraham H",
		CreationDate: 2020,
		Members:      []string{"Jerrie", "Kondu", "Ben"},
	}

	artist2 := Artist{
		Name:         "One Common",
		CreationDate: 2022,
		Members:      []string{"King", "Joe"},
	}

	artist3 := Artist{
		Name:         "Vicky Ivyic",
		CreationDate: 2017,
		Members:      []string{"Canny", "Don", "Becky"},
	}

	artists = append(artists, artist1, artist2, artist3)

	artist4 := Artist{
		Name:         "Maverick",
		CreationDate: 2021,
		Members:      []string{"Mary", "Cor"},
	}

	artists = append(artists, artist4)

	for i, artist := range artists {
		fmt.Printf("%d: %s (since %d)\n", i, artist.Name, artist.CreationDate)
	}

	relate := map[string][]string{
		"Lagos": {"12-03-2002", "23-07-2004", "15-02-2015"},
		"Benue": {"28-12-10-2022", "21-04-2023", "10-01-2024"},
		"Abuja": {"20-08-2024", "11-05-2025"},
	}

	for key, value := range relate {
		fmt.Printf("%s --> %v\n", key, value)
	}

	dates, ok := relate["Abuja"]
	if !ok {
		fmt.Println("Not found")
	} else {
		fmt.Println(dates)
	}

	date, ok := relate["Niger"]
	if !ok {
		fmt.Println("Not found")
	} else {
		fmt.Println(date)
	}
}
