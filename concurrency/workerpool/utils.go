package workerpool

// setOption function to set the configuration
func setOption(conf *Config, options ...func(*Config) error) error {
	for _, opt := range options {
		if err := opt(conf); err != nil {
			return err
		}
	}
	return nil
}

// Min function
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max function
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
