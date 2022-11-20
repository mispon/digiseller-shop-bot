package desc

type (
	Categories struct {
		Items []Category `json:"category"`
	}

	Category struct {
		Id   string
		Name string
		Sub  []SubCategory
	}

	SubCategory struct {
		Id   string
		Name string
	}
)
