package bot

import "strconv"

type (
	options struct {
		debug bool
	}

	Option func(o *options)
)

// WithDebug enables debug output
func WithDebug(debug string) Option {
	return func(o *options) {
		if v, err := strconv.ParseBool(debug); err != nil {
			o.debug = v
		}
	}
}
