package search

import (
	"fmt"
	"net/http"

	"github.com/mispon/xbox-store-bot/bot/desc"
	uhttp "github.com/mispon/xbox-store-bot/utils/http"
)

var (
	categoriesUrl = "%s/categories"
	productsUrl   = "%s/search"
	productUrl    = "%s/product"
)

func Categories(c *http.Client, host string) ([]desc.SubCategory, error) {
	subCategories, err := uhttp.Get[[]Category](c, fmt.Sprintf(categoriesUrl, host), map[string]string{})
	if err != nil {
		return nil, err
	}
	result := make([]desc.SubCategory, 0, len(subCategories))

	for _, item := range subCategories {
		result = append(result, desc.SubCategory{
			Id:   item.Name,
			Name: item.Description,
		})
	}

	return result, nil
}

func Search(c *http.Client, host, category, query string, count, skip int) ([]desc.Product, int, error) {
	queryParams := map[string]string{
		"categories": category,
		"count":      fmt.Sprintf("%d", count),
		"skip":       fmt.Sprintf("%d", skip),
		"query":      query,
	}
	products, err := uhttp.Get[Products](c, fmt.Sprintf(productsUrl, host), queryParams)
	if err != nil {
		return nil, 0, err
	}
	result := make([]desc.Product, 0, len(products.Items))

	for _, item := range products.Items {
		if item.Product.Prices["ARS"] != 0 || item.Product.Prices["TRY"] != 0 {
			result = append(result, desc.Product{
				Id:   item.Product.ID,
				Name: item.Product.Name,
			})
		}
	}

	return result, products.TotalItems, nil
}

func GetProduct(c *http.Client, host string, id string) (Product, error) {
	queryParams := map[string]string{
		"id": id,
	}
	product, err := uhttp.Get[Product](c, fmt.Sprintf(productUrl, host), queryParams)
	if err != nil {
		return Product{}, err
	}
	return product, nil
}

func (p Product) String() string {
	res := fmt.Sprintf("%s\n<a href='%s'>&#8205;</a>", p.Name, p.Img)
	return res
}

func (p Product) IsBackwardCompatibil() bool {
	if p.CategoryName == "Xbox360BackwardCompatibil" {
		return true
	}

	for _, gen := range p.Gens {
		if gen != "ConsoleGen7" {
			return true
		}
	}

	return false
}
