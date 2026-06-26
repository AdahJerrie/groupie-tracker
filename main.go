package main

import (
	"fmt"
	//"groupie/parts"
)

func main() {

	// getartist := parts.GetArtist()

	// for i, artist := range getartist {
	// 	fmt.Printf("%d: %s (since %d)\n", i, artist.Name, artist.CreationDate)
	// }

	// related := parts.Relation()

	// for _, value := range related {

	// 	fmt.Println(value)
	// }

	// pointer := parts.PointAddr()
	// fmt.Println(*pointer)

	year := 2026
	yearPtr := &year
	fmt.Println(yearPtr)
	fmt.Println(*yearPtr)

	*yearPtr = 2021
	fmt.Println(year)

}
