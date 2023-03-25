package search

type Category struct {
	Name        string
	Description string
}

type Product struct {
	ID             string
	Name           string
	Img            string
	Gens           []string
	GensCompatible []string
	Prices         map[string]float64
	CategoryName   string
	Type           string
}

type Products struct {
	Items []struct {
		Weight  int
		Product Product
	}
	TotalItems int
	LastPage   bool
}
