package controller

import (
	"WBABEProject-09/model"
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

// main.go 실행과 별개적으로 유닛 테스트를 할 수 있는 구조가 필요 - TODO -
// controller 대신 model unit 테스트를 활용하는 방안 고려

// var ownerUserId = "63a49555681a8bac48b5e47b"
// var customerUserId = "63a49555681a8bac48b5e47c"
// var menuId = "63a4951a681a8bac48b5e47a"

// func SetUpRouter() *gin.Engine {
// 	router := gin.Default()
// 	return router
// }

// 초기에는 gin.Default에 실제 Controller를 연결하는 방식으로 고안됨
// 고민하는데 시간이 너무 소요되서 실제 동작하는 서버에 http 요청을 날리는 형식으로 변경
// 개선 및 고민이 필요하지만 시간이 모잘라 일단 기능 개발에 집중
// func TestInsertMenuControl(t *testing.T) {
// 	menu := model.Menu{
// 		Category:        "중식",
// 		Name:            "짜장면",
// 		Price:           1000,
// 		Recommend:       false,
// 		OrderState:      1,
// 		OrderDailyLimit: 5,
// 	}
// 	jsonValue, _ := json.Marshal(menu)
// 	req, _ := http.NewRequest("POST", "http://localhost:8088/owner/menu", bytes.NewBuffer(jsonValue))
// 	req.Header.Set("userId", ownerUserId)
// 	req.Header.Set("menuId", menuId)
// 	req.Header.Set("Content-Type", "application/json")

// 	client := http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("TestInsertMenuControl Error: %s", err)
// 	}
// 	body, err := ioutil.ReadAll(res.Body)
// 	defer res.Body.Close()
// 	if err != nil {
// 		t.Errorf("TestInsertMenuControl Error: %s", err)
// 	}

// 	assert.Equal(t, http.StatusOK, res.StatusCode, string(body))
// }

// 초기 환경 구성시 user, menu, order 기본 데이터를 넣기위한 test
// 시간 문제로 Review 및 수정 init는 - TODO -
func TestInitControl(t *testing.T) {
	userList := []model.User{
		model.User{
			UserId:  1,
			Name:    "김주인",
			Email:   "email@email.email",
			Phone:   "01000000000",
			Address: "서울",
			Type:    1,
		},
		model.User{
			UserId:  2,
			Name:    "김손님1",
			Email:   "email@email.email",
			Phone:   "01000000000",
			Address: "서울",
			Type:    2,
		},
		model.User{
			UserId:  3,
			Name:    "김손님2",
			Email:   "email@email.email",
			Phone:   "01000000000",
			Address: "서울",
			Type:    2,
		},
	}

	for _, user := range userList {

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "http://localhost:8088/v1/user", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		_, err := client.Do(req)
		if err != nil {
			t.Errorf("TestInitControl user Error: %s", err)
		}
	}

	menuList := []model.Menu{
		model.Menu{
			Category:        "중식",
			Name:            "짜장면",
			Price:           5000,
			Recommend:       true,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "중식",
			Name:            "짬뽕면",
			Price:           5000,
			Recommend:       false,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "중식",
			Name:            "짬짜면",
			Price:           5500,
			Recommend:       true,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "중식",
			Name:            "탕수육",
			Price:           10000,
			Recommend:       false,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "한식",
			Name:            "김치찌개",
			Price:           8000,
			Recommend:       true,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "한식",
			Name:            "된장찌개",
			Price:           8000,
			Recommend:       false,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "한식",
			Name:            "삼겹살",
			Price:           15000,
			Recommend:       true,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "한식",
			Name:            "부대찌개",
			Price:           10000,
			Recommend:       false,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "양식",
			Name:            "피자",
			Price:           20000,
			Recommend:       true,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "양식",
			Name:            "스테이크",
			Price:           30000,
			Recommend:       false,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "양식",
			Name:            "치킨",
			Price:           20000,
			Recommend:       true,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
		model.Menu{
			Category:        "양식",
			Name:            "리조또",
			Price:           15000,
			Recommend:       false,
			OrderState:      1,
			OrderDailyLimit: 5,
		},
	}

	for _, menu := range menuList {

		jsonValue, _ := json.Marshal(menu)
		req, _ := http.NewRequest("POST", "http://localhost:8088/v1/owner/menu", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("userId", "1")
		client := http.Client{}
		_, err := client.Do(req)
		if err != nil {
			t.Errorf("TestInitControl menu Error: %s", err)
		}
	}

	orderList := []model.Order{
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 4,
					Name:   "탕수육",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 5,
					Name:   "김치찌개",
				},
				model.OrderMenu{
					MenuId: 7,
					Name:   "삼겹살",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 4,
					Name:   "탕수육",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 4,
					Name:   "탕수육",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 7,
					Name:   "삼겹살",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 5,
					Name:   "김치찌개",
				},
				model.OrderMenu{
					MenuId: 7,
					Name:   "삼겹살",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 4,
					Name:   "탕수육",
				},
				model.OrderMenu{
					MenuId: 5,
					Name:   "김치찌개",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
		model.Order{
			UserId: 2,
			Menu: []model.OrderMenu{
				model.OrderMenu{
					MenuId: 1,
					Name:   "짜장면",
				},
				model.OrderMenu{
					MenuId: 2,
					Name:   "짬뽕면",
				},
				model.OrderMenu{
					MenuId: 4,
					Name:   "탕수육",
				},
				model.OrderMenu{
					MenuId: 11,
					Name:   "치킨",
				},
			},
			Phone:   "01000000000",
			Address: "서울",
		},
	}

	for _, order := range orderList {

		jsonValue, _ := json.Marshal(order)
		req, _ := http.NewRequest("POST", "http://localhost:8088/v1/customer/order", bytes.NewBuffer(jsonValue))
		req.Header.Set("userId", "2")
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}
		_, err := client.Do(req)
		if err != nil {
			t.Errorf("TestInitControl order Error: %s", err)
		}
	}
}
