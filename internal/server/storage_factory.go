package server

import (
	"errors"

	"github.com/ngin8-beta/tfbackend/internal/storage"
	"github.com/spf13/viper"
)

func GetStorage() (storage.Storage, error) {
	backend := viper.GetString("storage")
	switch backend {
	case "local":
		s, err := storage.NewLocalStorage(viper.GetString("storage_local_dir"))
		if err != nil {
			return nil, err
		}
		return s, nil
	default:
		return nil, errors.New("unsupported storage backend: " + backend)
	}
}
