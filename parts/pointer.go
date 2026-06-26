package parts

func PointAddr() *int {
	year := 2026

	yearPtr := &year

	*yearPtr = 2000

	return yearPtr
}
