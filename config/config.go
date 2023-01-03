package config

import (
	"fmt"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Server struct {
		Port string
	}

	DB struct {
		Host                 string
		DB                   string
		UserCollection       string
		OrderCollection      string
		OrderSaveCollection  string
		MenuCollection       string
		ReviewCollection     string
		ReviewSaveCollection string
		IdSequence           string
		OrderCountSequence   string
	}

	Log struct {
		Level   string
		Fpath   string
		Msize   int
		Mage    int
		Mbackup int
	}
}

func NewConfig(fpath string) (*Config, error) {
	c := new(Config)
	file, err := os.Open(fpath)
	if err == nil {
		defer file.Close()
		//toml 파일 디코딩
		err := toml.NewDecoder(file).Decode(c)
		if err == nil {
			fmt.Println(c)
			return c, nil
		}
	}
	return nil, err
}
