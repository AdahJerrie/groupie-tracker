package parts

import "fmt"

func Relation() []string {

	relate := map[string][]string{
		"Lagos": {"12-03-2002", "23-07-2004", "15-02-2015"},
		"Benue": {"28-12-2022", "21-04-2023", "10-01-2024"},
		"Abuja": {"20-08-2024", "11-05-2025"},
	}

	var location string
	fmt.Println("enter location")
	fmt.Scanln(&location)

	dates, ok := relate[location]
	if !ok {
		return []string{"not found"}
	}
	return dates
}

// func main() {

// }
