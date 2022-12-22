package controller

import (
	"WBABEProject-09/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	md *model.Model
}

func NewCTL(rep *model.Model) (*Controller, error) {
	r := &Controller{md: rep}
	return r, nil
}

func (p *Controller) GetOK(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
	return
}
func (p *Controller) CheckUser(userId string) (int, error) {
	user, err := p.md.GetUserTypeByIdModel(userId)
	if err != nil {
		return 0, err
	}
	return user.Type, err
}

// InsertMenuControl godoc
//
//	@Summary		call InsertMenuControl, return menu data by model.Menu.
//	@Description	menu data 추가를 위한 기능.
//	@name			InsertMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	query	string	true	"User ID"
//	@Param			menu	body	model.Menu	true	"{category, name, price, recommend, orderState, orderDailyLimit}"
//	@Router			/owner/menu [post]
//	@Success		200	{object}	model.Menu
func (p *Controller) InsertMenuControl(c *gin.Context) {

	userId := c.GetHeader("userId")
	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "user ID가 유효하지 않습니다",
		})
		return
	}
	userType, err := p.CheckUser(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저 정보를 확인할 수 없습니다!",
			"error":   err.Error(),
		})
		return
	} else if userType != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저가 피주문자가 아닙니다!",
			"error":   err.Error(),
		})
		return
	}

	menu := model.NewMenu()
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴 정보가 잘못됬습니다!",
			"error":   err.Error(),
		})
		return
	}

	menuResult, err := p.md.InsertMenuModel(menu)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴를 추가하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, menuResult)
	return
}

// UpdateMenuControl godoc
//
//	@Summary		call UpdateMenuControl, return menu data by model.Menu.
//	@Description	menu data 수정을 위한 기능.
//	@name			UpdateMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	query	string	true	"User ID"
//	@Param			menu	body	model.Menu	true	"{id, category, name, price, recommend, orderState, orderDailyLimit}"
//	@Router			/owner/menu [put]
//	@Success		200	{object}	model.Menu
func (p *Controller) UpdateMenuControl(c *gin.Context) {

	userId := c.GetHeader("userId")
	menuId := c.GetHeader("menuId")
	if menuId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "menu ID가 유효하지 않습니다",
		})
		return
	} else if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "user ID가 유효하지 않습니다",
		})
		return
	}
	userType, err := p.CheckUser(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저 정보를 확인할 수 없습니다!",
			"error":   err.Error(),
		})
		return
	} else if userType != 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저가 피주문자가 아닙니다!",
			"error":   err.Error(),
		})
		return
	}

	menu := model.Menu{}
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴 정보가 잘못됬습니다!",
			"error":   err.Error(),
		})
		return
	}

	if err := p.md.UpdateMenuModel(menuId, menu); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴를 수정하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}
