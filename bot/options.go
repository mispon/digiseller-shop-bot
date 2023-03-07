package bot

type (
	searchOpt struct {
		enabled                  bool
		url                      string
		universalProductId       int
		universalProductOptionId int
	}
	options struct {
		sellerId string
		debug    bool
		search   searchOpt
	}

	Option func(o *options)
)

// WithSeller sets seller id
func WithSeller(sellerId string) Option {
	return func(o *options) {
		o.sellerId = sellerId
	}
}

// WithDebug enables debug output
func WithDebug(debug bool) Option {
	return func(o *options) {
		o.debug = debug
	}
}

func WithSearch(searchUrl string, productId, optionId int) Option {
	return func(o *options) {
		if searchUrl != "" {
			o.search.enabled = true
			o.search.url = searchUrl
			o.search.universalProductId = productId
			o.search.universalProductOptionId = optionId
		}
	}
}
