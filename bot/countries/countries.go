package countries

import "errors"

type country struct {
	Flag     string
	Currency string
}

var (
	countries = map[string]country{"AR": {"🇦🇷", "ARS"}, "TR": {"🇹🇷", "TRY"}}

	countryNotFoundErr = errors.New("Country not found")
)

func GetCountry(countryName string) (country, error) {
	if country, ok := countries[countryName]; ok {
		return country, nil
	}
	return country{}, countryNotFoundErr
}
