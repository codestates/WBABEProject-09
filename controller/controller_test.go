package controller

import (
	"WBABEProject-09/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var ownerUserId = "63a49555681a8bac48b5e47b"
var customerUserId = "63a49555681a8bac48b5e47c"
var menuId = "63a4951a681a8bac48b5e47a"

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// 초기에는 gin.Default에 실제 Controller를 연결하는 방식으로 고안됨
// 고민하는데 시간이 너무 소요되서 실제 동작하는 서버에 http 요청을 날리는 형식으로 변경
// 개선 및 고민이 필요하지만 시간이 모잘라 일단 기능 개발에 집중
func TestInsertMenuControl(t *testing.T) {
	menu := model.Menu{
		Category:        "중식",
		Name:            "짜장면",
		Price:           1000,
		Recommend:       false,
		OrderState:      1,
		OrderDailyLimit: 5,
	}
	jsonValue, _ := json.Marshal(menu)
	req, _ := http.NewRequest("POST", "http://localhost:8080/owner/menu", bytes.NewBuffer(jsonValue))
	req.Header.Set("userId", ownerUserId)
	req.Header.Set("menuId", menuId)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("TestInsertMenuControl Error: %s", err)
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		t.Errorf("TestInsertMenuControl Error: %s", err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, string(body))
}
