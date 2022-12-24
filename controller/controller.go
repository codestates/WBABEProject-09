package controller

import (
	"WBABEProject-09/model"
	ut "WBABEProject-09/type/user"
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

// user ID 및 Type에대한 유효성 검사를 공통적으로 수행하기 위해 선언
func (p *Controller) UserValidation(c *gin.Context, targetUserType int, userId string) bool {

	if userId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "user ID가 유효하지 않습니다",
		})
		return false
	}
	userType, err := p.CheckUser(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저 정보를 확인할 수 없습니다!",
			"error":   err.Error(),
		})
		return false
	} else if userType != targetUserType {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":  "유저 타입이 일치하지 않습니다!",
			"userType": ut.GetUserTypeText(userType),
		})
		return false
	}

	return true
}

// menu ID에 대한 유효성 검사를 공통적으로 수행하기 위해 선언
func (p *Controller) MenuValidation(c *gin.Context, menuId string) bool {
	if menuId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "menu ID가 유효하지 않습니다",
		})
		return false
	}

	return true
}

// menu 관련 controller에서 사용되는 Bind이 중복되기에 별도로 분리
func (p *Controller) MenuBind(c *gin.Context, menu *model.Menu) bool {

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴 정보가 잘못됬습니다!",
			"error":   err.Error(),
		})
		return false
	}
	return true
}

// InsertMenuControl godoc
//
//	@Summary		call InsertMenuControl, return menu data by model.Menu.
//	@Description	menu data 추가를 위한 기능.
//	@name			InsertMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menu	body	model.Menu	true	"{category, name, price, recommend, orderState, orderDailyLimit}"
//	@Router			/owner/menu [post]
//	@Success		200	{object}	model.Menu
func (p *Controller) InsertMenuControl(c *gin.Context) {

	userId := c.GetHeader("userId")
	if !p.UserValidation(c, ut.TypeOwner, userId) {
		return
	}
	menu := model.NewMenu()
	if !p.MenuBind(c, &menu) {
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
//	@Summary		call UpdateMenuControl, return result by json.
//	@Description	menu data 수정을 위한 기능.
//	@name			UpdateMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menuId	header	string	true	"Menu ID"
//	@Param			menu	body	model.Menu	true	"{category, name, price, recommend, orderState, orderDailyLimit}"
//	@Router			/owner/menu [put]
//	@Success		200	{object}
func (p *Controller) UpdateMenuControl(c *gin.Context) {

	userId := c.GetHeader("userId")
	menuId := c.GetHeader("menuId")

	if !p.UserValidation(c, ut.TypeOwner, userId) {
		return
	}
	if !p.MenuValidation(c, menuId) {
		return
	}

	menu := model.Menu{}
	if !p.MenuBind(c, &menu) {
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

// DeleteMenuControl godoc
//
//	@Summary		call DeleteMenuControl, return result by json.
//	@Description	menu data 삭제을 위한 기능.
//	@name			DeleteMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menuId	header	string	true	"Menu ID"
//	@Router			/owner/menu [delete]
//	@Success		200	{object}	controller
func (p *Controller) DeleteMenuControl(c *gin.Context) {

	userId := c.GetHeader("userId")
	menuId := c.GetHeader("menuId")

	if !p.UserValidation(c, ut.TypeOwner, userId) {
		return
	}
	if !p.MenuValidation(c, menuId) {
		return
	}

	if err := p.md.DeleteMenuModel(menuId); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴를 삭제하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}
