package bot

type (
	options struct {
		sellerId string
		debug    bool
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
