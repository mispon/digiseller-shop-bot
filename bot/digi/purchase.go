package digi

import (
	"fmt"
	"net/http"
	"strconv"

	uhttp "github.com/mispon/xbox-store-bot/utils/http"
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
