package cache

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mispon/xbox-store-bot/bot/desc"
	uhttp "github.com/mispon/xbox-store-bot/utils/http"
	"go.uber.org/zap"
)

const (
	categoryUrl    = "https://api.digiseller.ru/api/categories?seller_id=%s"
	productListUrl = "https://api.digiseller.ru/api/shop/products?seller_id=%s&category_id=%s&page=%d"
	productDataUrl = "https://api.digiseller.ru/api/products/%s/data"
)

type (
	cache struct {
		logger   *zap.Logger
		client   *http.Client
		sellerId string

		mu   sync.RWMutex
		data data
	}

	data struct {
		Categories []desc.Category
		Products   map[string][]desc.Product
	}
)

func New(logger *zap.Logger, sellerId string, load bool) (*cache, error) {
	c := &cache{
		logger:   logger.Named("cache"),
		client:   http.DefaultClient,
		sellerId: sellerId,
	}

	startTime := time.Now()
	if load {
		if err := c.load(); err != nil {
			return nil, err
		}
		go c.refresh()
	}

	c.logger.Info("cache loaded", zap.Duration("duration", time.Since(startTime)))
	return c, nil
}

func (c *cache) Categories() []desc.Category {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.data.Categories
}

func (c *cache) SubCategory(categoryId string) (string, []desc.SubCategory, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, category := range c.data.Categories {
		if category.Id == categoryId {
			return category.Name, category.Sub, true
		}
	}

	return "", nil, false
}

func (c *cache) Products(subCategoryId string, page, total int) (string, []desc.Product, bool, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var subCategoryName string
LOOP:
	for _, category := range c.data.Categories {
		for _, subCategory := range category.Sub {
			if subCategory.Id == subCategoryId {
				subCategoryName = subCategory.Name
				break LOOP
			}
		}
	}

	products, ok := c.data.Products[subCategoryId]
	if !ok {
		return subCategoryName, nil, false, false
	}

	hasMore := false
	result := make([]desc.Product, 0, total)
	for i := page * total; i < len(products); i++ {
		if len(result) == total {
			hasMore = true
			break
		}
		result = append(result, products[i])
	}

	return subCategoryName, result, hasMore, true
}

func (c *cache) Product(subCategoryId, productId string) (desc.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if products, ok := c.data.Products[subCategoryId]; ok {
		for _, product := range products {
			if product.Id == productId {
				return product, true
			}
		}
	}

	return desc.Product{}, false
}

func (c *cache) Search(text string) ([]desc.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var products []desc.Product
	for _, productsList := range c.data.Products {
		for _, product := range productsList {
			name := strings.ToLower(product.Name)
			if strings.Contains(name, text) {
				products = append(products, product)
			}
		}
	}

	if len(products) > 0 {
		return products, true
	}
	return nil, false
}

func (c *cache) load() error {
	categoriesResp, err := uhttp.Get[desc.Categories](c.client, fmt.Sprintf(categoryUrl, c.sellerId))
	if err != nil {
		c.logger.Error("failed to load categories", zap.Error(err))
		return err
	}

	productsMap := make(map[string][]desc.Product)
	for _, category := range categoriesResp.Items {
		for _, sc := range category.Sub {
			page := 1
			for {
				productsResp, pErr := uhttp.Get[desc.Products](c.client, fmt.Sprintf(productListUrl, c.sellerId, sc.Id, page))
				if pErr != nil {
					c.logger.Error("failed to load products", zap.Error(err))
					return pErr
				}

				items := productsMap[sc.Id]
				items = append(items, productsResp.Items...)
				productsMap[sc.Id] = items

				page++
				if tp, cErr := strconv.Atoi(productsResp.Pages); cErr != nil || page > tp {
					break
				}
			}
		}
	}

	for sc, productsList := range productsMap {
		for i, product := range productsList {
			info, pErr := uhttp.Get[desc.ProductFull](c.client, fmt.Sprintf(productDataUrl, product.Id))
			if pErr != nil {
				c.logger.Error("failed to get full product data", zap.String("product", product.Name), zap.Error(pErr))
				continue
			}
			product.Info = info.Product.Info
			product.AddInfo = info.Product.AddInfo

			productsList[i] = product
		}
		productsMap[sc] = productsList
	}

	c.mu.Lock()
	c.data.Categories = categoriesResp.Items
	c.data.Products = productsMap
	c.mu.Unlock()

	return nil
}

func (c *cache) refresh() {
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		if err := c.load(); err != nil {
			c.logger.Error("failed to load data", zap.Error(err))
		}
	}
}
