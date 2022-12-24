package config

import (
	"fmt"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Server struct {
		Mode string
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

	if file, err := os.Open(fpath); err != nil {
		return nil, err
	} else {
		defer file.Close()
		//toml 파일 디코딩
		if err := toml.NewDecoder(file).Decode(c); err != nil {
			return nil, err
		} else {
			fmt.Println(c)
			return c, nil
		}
	}
}
