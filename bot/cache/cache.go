package cache

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mispon/xbox-store-bot/bot/desc"
	uhttp "github.com/mispon/xbox-store-bot/utils/http"
	"go.uber.org/zap"
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

func New(logger *zap.Logger, sellerId string) (*cache, error) {
	c := &cache{
		logger:   logger.Named("cache"),
		client:   http.DefaultClient,
		sellerId: sellerId,
	}

	if err := c.load(); err != nil {
		return nil, err
	}
	go c.refresh()

	c.logger.Info("cache loaded")
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

func (c *cache) Products(subCategoryId string) (string, []desc.Product, bool) {
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
	return subCategoryName, products, ok
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

func (c *cache) load() error {
	categoriesResp, err := uhttp.Get[desc.Categories](c.client, fmt.Sprintf(categoryUrl, c.sellerId))
	if err != nil {
		c.logger.Error("failed to load categories", zap.Error(err))
		return err
	}

	products := make(map[string][]desc.Product)
	for _, category := range categoriesResp.Items {
		for _, sc := range category.Sub {
			page := 1
			for {
				productsResp, pErr := uhttp.Get[desc.Products](c.client, fmt.Sprintf(productUrl, c.sellerId, sc.Id, page))
				if pErr != nil {
					c.logger.Error("failed to load products", zap.Error(err))
					return pErr
				}

				items := products[sc.Id]
				items = append(items, productsResp.Items...)
				products[sc.Id] = items

				page++
				if tp, cErr := strconv.Atoi(productsResp.Pages); cErr != nil || page > tp {
					break
				}
			}
		}
	}

	c.mu.Lock()
	c.data.Categories = categoriesResp.Items
	c.data.Products = products
	c.mu.Unlock()

	return nil
}

func (c *cache) refresh() {
	for range time.After(10 * time.Minute) {
		if err := c.load(); err != nil {
			c.logger.Error("failed to load data", zap.Error(err))
		}
	}
}
