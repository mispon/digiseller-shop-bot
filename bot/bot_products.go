package bot

import (
	"sync"
)

type PurchaseType string
type ProductsDisplayType string

var (
	purchaseTypeKey PurchaseType = "key"
	purchaseTypeAcc PurchaseType = "acc"

	ProductsDisplayTypeMinPrice ProductsDisplayType = "minPrice"
	ProductsDisplayTypeMaxPrice ProductsDisplayType = "maxPrice"
	ProductsDisplayTypeAll      ProductsDisplayType = "all"

	purchaseTypes        = map[PurchaseType]bool{purchaseTypeKey: true, purchaseTypeAcc: true}
	ProductsDisplayTypes = map[ProductsDisplayType]bool{
		ProductsDisplayTypeMinPrice: true,
		ProductsDisplayTypeMaxPrice: true,
		ProductsDisplayTypeAll:      true,
	}
)

type BotProduct struct {
	PurchaseType           PurchaseType `json:"purchaseType"`
	Country                string       `json:"country"`
	MinPrice               int          `json:"minPrice"`
	SkipBackwardCompatibil bool         `json:"skipBackwardCompatibil"`
}

type userConfig struct {
	sync.RWMutex
	ConversionRates     map[string]float64  `json:"conversionRates"`
	BotProducts         []BotProduct        `json:"botProducts"`
	ProductsDisplayType ProductsDisplayType `json:"productsDisplayType"`
}
