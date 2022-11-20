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
		Id    string
		Name  string
		Info  string
		Price string `json:"price_rub"`
		Curr  string `json:"base_currency"`
	}
)

var (
	htmlReplacements = map[string]string{
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

	return fmt.Sprintf("%s\n\n%s\n\n%s %s\n<a href='%s'>&#8205;</a>", name, info, p.Price, "RUB", imageUrl)
}

func (p Product) PaymentURL(sellerId string) string {
	return fmt.Sprintf(paymentURL, p.Id, sellerId)
}
