package filter

// config represents the internal config of the parser functions.
type config struct {
	// useNumber indicates that json.Number needs to be returned instead of int/float64 values.
	useNumber bool
}

type configOption func(config) config

// WithUseNumber set use number to true
func WithUseNumber() configOption {
	return func(c config) config {
		c.useNumber = true
		return c
	}
}

func getConfig(options ...configOption) config {
	c := config{}
	for _, opt := range options {
		c = opt(c)
	}
	return c
}
