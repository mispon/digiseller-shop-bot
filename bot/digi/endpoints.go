package digi

import (
	"fmt"
	"net/http"
	"strconv"

	uhttp "github.com/mispon/digiseller-shop-bot/utils/http"
)

const (
	ReviewsLogoUrl = "https://my.digiseller.ru/preview/257605/logo_20220513174904.png"
	CategoryUrl    = "https://api.digiseller.ru/api/categories?seller_id=%s"
	ProductListUrl = "https://api.digiseller.ru/api/shop/products?seller_id=%s&category_id=%s&page=%d"
	ProductDataUrl = "https://api.digiseller.ru/api/products/%s/data"

	imageUrlTempl       = "http://graph.digiseller.ru/img.ashx?id_d=%s&maxlength=400"
	paymentURLTempl     = "https://www.digiseller.market/asp2/pay_api.asp?id_d=%s&curr=API_5011_RUB&_subcurr=&lang=ru-RU&_ids_shop=%s&failpage=https://x-box-store.ru"
	purchasesOptionsUrl = "https://api.digiseller.ru/api/purchases/options"
)

func CustomProductPaymentURL(
	client *http.Client,
	sellerId, productName string,
	paymentProductID, paymentProductOption, price int) string {

	if paymentProductID != 0 {
		ppo := ProductPurchasesOptions{
			ProductID: paymentProductID,
			UnitCnt:   price,
			IP:        "127.0.0.1",
		}
		ppo.Options = append(ppo.Options, ProductPurchasesOption{
			ID:    paymentProductOption,
			Value: PPOValue{Text: productName}})
		if info, pErr := uhttp.Post[ProductPaymentData](client, purchasesOptionsUrl, ppo); pErr == nil {
			if info.Retval == 0 {
				PaymentProductIDString := strconv.Itoa(paymentProductID)
				return fmt.Sprintf(paymentURLTempl, PaymentProductIDString, sellerId) + fmt.Sprintf("&id_po=%d", info.IDPo)
			}
		}
	}
	return ""
}

func ProductPaymentURL(productID, sellerID string) string {
	return fmt.Sprintf(paymentURLTempl, productID, sellerID)
}

func ProductImageUrl(id string) string {
	return fmt.Sprintf(imageUrlTempl, id)
}
