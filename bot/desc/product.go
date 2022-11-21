package desc

import (
	"fmt"
	"strings"
)

const (
	paymentURL    = "https://www.digiseller.market/asp2/pay_api.asp?id_d=%s&curr=API_5011_RUB&_subcurr=&lang=ru-RU&_ids_shop=%s&failpage=https://x-box-store.ru"
	imageUrlTempl = "http://graph.digiseller.ru/img.ashx?id_d=%s&maxlength=400"
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
	imageUrl := fmt.Sprintf(imageUrlTempl, p.Id)

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

	res := fmt.Sprintf("%s", addInfo)
	return res
}

func (p Product) PaymentURL(sellerId string) string {
	return fmt.Sprintf(paymentURL, p.Id, sellerId)
}
