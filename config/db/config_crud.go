package db

func GetConfig[T any](key string, def ...T) T {
	var defValue T

	// If a default value is provided, use it
	if len(def) > 0 {
		defValue = def[0]
	}

	var config Config[T]
	dbi.Where("key = ?", key).First(&config)

	if config.Key == "" {
		return defValue
	}

	return config.Value
}

func SetConfig[T any](key string, value T) {
	config := Config[T]{
		Key:   key,
		Value: value,
	}

	dbi.Save(&config)
}

func DeleteConfig(key string) {
	dbi.Delete(&Config[any]{Key: key})
}
