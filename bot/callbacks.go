package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

type callbackType int

const (
	Categories callbackType = iota
	SubCategory
	Products
	Product
	ProductInstruction
	Back
)

type (
	callbackEntity struct {
		cbType     callbackType
		id         string
		parentType callbackType
		parentIds  []string
	}

	callbackFn func(upd tgbotapi.Update, entity callbackEntity)
)

func (b *bot) initCallbacks() {
	b.callbacks = map[callbackType]callbackFn{
		SubCategory:        b.SubCategoryCallback,
		Products:           b.ProductsCallback,
		Product:            b.ProductCallback,
		ProductInstruction: b.ProductInstructionCallback,
		Back:               b.BackCallback,
	}
}

func marshallCb(data callbackEntity) string {
	return fmt.Sprintf(
		"%d;%s;%d;%s",
		data.cbType,
		data.id,
		data.parentType,
		strings.Join(data.parentIds, "."),
	)
}

func unmarshallCb(data string) callbackEntity {
	d := strings.Split(data, ";")

	cbType, _ := strconv.Atoi(d[0])
	pType, _ := strconv.Atoi(d[2])

	return callbackEntity{
		cbType:     callbackType(cbType),
		id:         d[1],
		parentType: callbackType(pType),
		parentIds:  strings.Split(d[3], "."),
	}
}
