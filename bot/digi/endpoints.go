package digi

const (
	ReviewslogoUrl = "https://my.digiseller.ru/preview/257605/logo_20220513174904.png"
	CategoryUrl    = "https://api.digiseller.ru/api/categories?seller_id=%s"
	ProductListUrl = "https://api.digiseller.ru/api/shop/products?seller_id=%s&category_id=%s&page=%d"
	ProductDataUrl = "https://api.digiseller.ru/api/products/%s/data"

	imageUrlTempl       = "http://graph.digiseller.ru/img.ashx?id_d=%s&maxlength=400"
	paymentURLTempl     = "https://www.digiseller.market/asp2/pay_api.asp?id_d=%s&curr=API_5011_RUB&_subcurr=&lang=ru-RU&_ids_shop=%s&failpage=https://x-box-store.ru"
	purchasesOptionsUrl = "https://api.digiseller.ru/api/purchases/options"
)
