package parts

type Artist struct {
	Name         string `json: "name"`
	CreationDate int
	Members      []string
}

func GetArtist() []Artist {
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

	return artists

}
