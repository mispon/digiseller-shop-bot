package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type callbackType int

const (
	Categories callbackType = iota
	SubCategory
	Products
	Product
	ProductInstruction
	Search
	SearchParams
	SearchInstruction
	SearchSubCategory
	SearchProduct
	Back
)

type (
	callbackEntity struct {
		cbType     callbackType
		id         string
		parentType callbackType
		parentIds  []string
		page       int
		skip       int
	}

	callbackFn func(upd tgbotapi.Update, entity callbackEntity)
)

func (c callbackEntity) Clone() callbackEntity {
	return callbackEntity{
		cbType:     c.cbType,
		id:         c.id,
		parentType: c.parentType,
		parentIds:  c.parentIds,
		page:       c.page,
	}
}

func (b *bot) initCallbacks() {
	b.callbacks = map[callbackType]callbackFn{
		SubCategory:        b.SubCategoryCallback,
		Products:           b.ProductsCallback,
		Product:            b.ProductCallback,
		ProductInstruction: b.ProductInstructionCallback,
		SearchSubCategory:  b.SearchSubCategoryCallback,
		SearchParams:       b.SearchParamsCallback,
		SearchProduct:      b.SearchProductCallback,
		Back:               b.BackCallback,
		Search:             b.SearchCallback,
		SearchInstruction:  b.SearchInstructionCallback,
	}
}

func marshallCb(data callbackEntity) string {
	return fmt.Sprintf(
		"%d;%s;%d;%s;%d;%d",
		data.cbType,
		data.id,
		data.parentType,
		strings.Join(data.parentIds, "."),
		data.page,
		data.skip,
	)
}

func unmarshallCb(data string) callbackEntity {
	d := strings.Split(data, ";")

	var cbType int
	if len(d) > 0 {
		cbType, _ = strconv.Atoi(d[0])
	}

	var id string
	if len(d) > 1 {
		id = d[1]
	}

	var pType int
	if len(d) > 2 {
		pType, _ = strconv.Atoi(d[2])
	}

	var parentIds []string
	if len(d) > 3 {
		parentIds = strings.Split(d[3], ".")
	}

	var page int
	if len(d) > 4 {
		page, _ = strconv.Atoi(d[4])
	}

	var skip int
	if len(d) > 5 {
		skip, _ = strconv.Atoi(d[5])
	}

	return callbackEntity{
		cbType:     callbackType(cbType),
		id:         id,
		parentType: callbackType(pType),
		parentIds:  parentIds,
		page:       page,
		skip:       skip,
	}
}
