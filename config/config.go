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
	/* [코드리뷰]
	 * 해당 코드에는 하나의 function에서 간결한 이중 조건문이 발생하고 있습니다.
	 * 15 line으로 구성된 한눈에 들어오는 함수에서는 사용하기 적합합니다.
	 * 그러나 코드 라인수가 많아지고, 비즈니스 로직이 풍부해 지는 것을 고려한다면
	 * return으로 나갈 수 빠지는 case를 최소화하고, if문을 줄이는 방향으로 개발이 진행되어야 합니다. 
	 * 점진적으로 코드의 가독성이 조금 더 향상하게 됩니다.
	 * as-is: 
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
	 * to-be:
	 if file, err := os.Open(fpath); err == nil {
			defer file.Close()
			//toml 파일 디코딩
			if err := toml.NewDecoder(file).Decode(c); err == nil {
				fmt.Println(c)
				return c, nil
			} 
		}
		return nil, err
	 */
}
