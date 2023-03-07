package desc

import (
	"fmt"
	"strings"

	"github.com/mispon/xbox-store-bot/bot/digi"
)

type (
	Products struct {
		Items []Product `json:"product"`
		Pages string    `json:"totalPages"`
	}

	Product struct {
		Id      string
		Name    string
		Info    string `json:"-"`
		AddInfo string `json:"-"`
		Price   string `json:"price_rub"`
		Curr    string `json:"base_currency"`
	}

	ProductFull struct {
		Product struct {
			Info    string `json:"info"`
			AddInfo string `json:"add_info"`
		}
	}
)

var (
	htmlReplacements = map[string]string{
		"<br />":       "\n",
		"<br/>":        "\n",
		"<br>":         "\n",
		"<attention>":  "<b>",
		"</attention>": "</b>",
		"<delivery>":   "<i>",
		"</delivery>":  "</i>",
	}
)

func (p Product) String() string {
	imageUrl := digi.ProductImageUrl(p.Id)

	name := p.Name
	info := p.Info
	for k, v := range htmlReplacements {
		name = strings.ReplaceAll(name, k, v)
		info = strings.ReplaceAll(info, k, v)
	}

	res := fmt.Sprintf("%s\n\n%s\n\n%s %s\n<a href='%s'>&#8205;</a>", name, info, p.Price, "RUB", imageUrl)
	return res
}

func (p Product) Instruction() string {
	addInfo := p.AddInfo
	for k, v := range htmlReplacements {
		addInfo = strings.ReplaceAll(addInfo, k, v)
	}
	return addInfo
}

func (p Product) PaymentURL(sellerId string) string {
	return digi.ProductPaymentURL(p.Id, sellerId)
}
