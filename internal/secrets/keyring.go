package secrets

import (
	"github.com/zalando/go-keyring"
)

const service = "phoneinfoga-desktop"

func Set(key, value string) error {
	return keyring.Set(service, key, value)
}

func Get(key string) (string, error) {
	return keyring.Get(service, key)
}

func Delete(key string) error {
	return keyring.Delete(service, key)
}
