package controller

import (
	"WBABEProject-09/model"
	ut "WBABEProject-09/type/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 공통적으로 error에 대한 처리가 부실함
// - TODO - 로그: 상세하게 강화, http: statusCode를 활용해 각 상황에 대한 error 구분 필요, error: 에러 발생사항 관련된 개별 정보(발생 위치 등)를 좀더 상세하게 명시
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
func (p *Controller) CheckUser(userId int) (int, error) {
	user, err := p.md.GetUserTypeByIdModel(userId)
	if err != nil {
		return 0, err
	}
	return user.Type, err
}

// user ID 및 Type에대한 유효성 검사를 공통적으로 수행하기 위해 선언
func (p *Controller) UserValidation(c *gin.Context, targetUserType int, userId int) bool {

	if userId == 0 {
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
func (p *Controller) MenuValidation(c *gin.Context, menuId int) bool {
	if menuId == 0 {
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

func (p *Controller) OrderBind(c *gin.Context, order *model.Order) bool {

	if err := c.ShouldBindJSON(&order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "오더 정보가 잘못됬습니다!",
			"error":   err.Error(),
		})
		return false
	}
	return true
}
func (p *Controller) ReviewBind(c *gin.Context, review *model.Review) bool {

	if err := c.ShouldBindJSON(&review); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "리뷰 정보가 잘못됬습니다!",
			"error":   err.Error(),
		})
		return false
	}
	return true
}

// InsertUserControl godoc
//
//	@Summary		call InsertUserControl, return menu data by Json.
//	@Description	user data 추가를 위한 기능.
//	@name			InsertUserControl
//	@Accept			json
//	@Produce		json
//	@Param			user	body	model.User	true	"{userId, name, email, phone, address, type}"
//	@Router			/user [post]
//	@Success		200	{object}	string
func (p *Controller) InsertUserControl(c *gin.Context) {

	user := model.NewUser()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저 정보가 잘못됬습니다!",
			"error":   err.Error(),
		})
		return
	}
	userResult, err := p.md.InsertUserModel(user)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "유저를 추가하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, userResult)
	return
}

// GetMenuControl godoc
//
//	@Summary		call GetMenuControl, return menu data by []model.Menu
//	@Description	menu data 조회를 위한 기능.
//	@name			GetMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			sortBy	query	string	false	"recommend, star, orderCount, date"
//	@Param			checkReview	query	int	false	"리뷰 확인 여부"
//	@Router			/customer/menu [get]
//	@Router			/owner/menu [get]
//	@Success		200	{object}	[]model.Menu
func (p *Controller) GetMenuControl(c *gin.Context) {

	sortBy := c.Query("sortBy")
	checkReview, _ := strconv.Atoi(c.GetHeader("checkReview"))

	if sortBy == "" {
		sortBy = "_id"
	}
	menuResult, err := p.md.GetMenuModel(sortBy, checkReview)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴를 검색하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, menuResult)
	return
}

// GetMenuDetailControl godoc
//
//	@Summary		call GetMenuDetailControl, return menu data by []model.Menu
//	@Description	menu data 조회를 위한 기능.
//	@name			GetMenuControl
//	@Accept			json
//	@Produce		json
//	@Param			menuId	header	string	true	"Menu Id"
//	@Router			/customer/menu/detail [get]
//	@Router			/owner/menu/detail [get]
//	@Success		200	{object}	[]model.Menu
func (p *Controller) GetMenuDetailControl(c *gin.Context) {

	menuId, _ := strconv.Atoi(c.GetHeader("menuId"))
	if !p.MenuValidation(c, menuId) {
		return
	}
	menuResult, err := p.md.GetMenuDetailModel(menuId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "메뉴를 검색하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, menuResult)
	return
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

	userId, _ := strconv.Atoi(c.GetHeader("userId"))
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
//	@Success		200	{object}	string
func (p *Controller) UpdateMenuControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))
	menuId, _ := strconv.Atoi(c.GetHeader("menuId"))

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
//	@Success		200	{object}	string
func (p *Controller) DeleteMenuControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))
	menuId, _ := strconv.Atoi(c.GetHeader("menuId"))

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

// GetOrderControl godoc
//
//	@Summary		call GetOrderControl, return order data by []model.Order
//	@Description	order data 조회를 위한 기능.
//	@name			GetOrderControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Router			/customer/order [get]
//	@Router			/owner/order [get]
//	@Success		200	{object}	[]model.Order
func (p *Controller) GetOrderControl(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if userId == 0 {
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
	}

	inOrderResult, inOrderErr := p.md.GetInOrderModel(userId, userType)
	doneOrderResult, doneOrderErr := p.md.GetDoneOrderModel(userId, userType)

	if inOrderErr != nil || doneOrderErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "오더를 검색하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	// 진행오더와 완료오더 확인 api를 분리할 수 있지만,
	c.JSON(200, gin.H{"진행오더": inOrderResult, "완료오더": doneOrderResult})
	return
}

// InsertCustomerOrderControl godoc
//
//	@Summary		call InsertCustomerOrderControl, return result by json.
//	@Description	order data 추가을 위한 기능.
//	@name			InsertCustomerOrderControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menu	body	model.Order	true	"{userId, menu[{menuID, name}], phone, address}"
//	@Router			/customer/order [post]
//	@Success		200	{object}	string
func (p *Controller) InsertCustomerOrderControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if !p.UserValidation(c, ut.TypeCustomer, userId) {
		return
	}

	NewOrder := model.NewOrder()
	if !p.OrderBind(c, &NewOrder) {
		return
	}
	NewOrder.UserId = userId

	orderResult, err := p.md.InsertOrderModel(NewOrder)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "오더를 추가하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	orderNumber := fmt.Sprintf("%s_%d", orderResult.OrderDay, orderResult.OrderId)
	c.JSON(200, gin.H{"msg": "오더 추가 완료", "오더번호: ": orderNumber})
	return
}

// UpdateCustomerOrderControl godoc
//
//	@Summary		call UpdateCustomerOrderControl, return result by json.
//	@Description	주문자의 order data 수정을 위한 기능.
//	@name			UpdateCustomerOrderControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menu	body	model.Order	true	"{userId, orderDate, orderID , menu[{menuID, name}], phone, address}"
//	@Router			/customer/order [put]
//	@Success		200	{object}	string
func (p *Controller) UpdateCustomerOrderControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if !p.UserValidation(c, ut.TypeCustomer, userId) {
		return
	}

	order := model.Order{}
	if !p.OrderBind(c, &order) {
		return
	}
	order.UserId = userId

	if err := p.md.UpdateCustomerOrderModel(order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "오더를 수정하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}

// UpdateOwnerOrderControl godoc
//
//	@Summary		call UpdateOwnerOrderControl, return result by json.
//	@Description	오너가 order state 수정을 위한 기능.
//	@name			UpdateOwnerOrderControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menu	body	model.Order	true	"{orderDate, orderID , state}"
//	@Router			/owner/order [put]
//	@Success		200	{object}	string
func (p *Controller) UpdateOwnerOrderControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if !p.UserValidation(c, ut.TypeOwner, userId) {
		return
	}

	order := model.Order{}
	if !p.OrderBind(c, &order) {
		return
	}
	order.UserId = userId

	if err := p.md.UpdateOwnerOrderModel(order); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "오더 상태를 수정하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}

// GetReviewControl godoc
//
//	@Summary		call GetReviewControl, return result by []model.Review.
//	@Description	review data 확인을 위한 기능.
//	@name			GetReviewControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			sortBy	query	string	false	"정렬할 컬럼명"
//	@Router			/customer/order/review [get]
//	@Success		200	{object}	[]model.Review
func (p *Controller) GetReviewControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if !p.UserValidation(c, ut.TypeCustomer, userId) {
		return
	}

	sortBy := c.Query("sortBy")

	if sortBy == "" {
		sortBy = "_id"
	}

	reviewResult, err := p.md.GetReviewModel(userId, sortBy)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "리뷰를 조회하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, reviewResult)
	return
}

// InsertReviewControl godoc
//
//	@Summary		call InsertReviewControl, return result by json.
//	@Description	review data 추가을 위한 기능.
//	@name			InsertReviewControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			review	body	model.Review	true	"{orderDay, orderId, star, content}"
//	@Router			/customer/order/review [post]
//	@Success		200	{object}	model.Review
func (p *Controller) InsertReviewControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if !p.UserValidation(c, ut.TypeCustomer, userId) {
		return
	}

	newReview := model.NewReview()
	if !p.ReviewBind(c, &newReview) {
		return
	}
	newReview.UserId = userId

	reviewResult, err := p.md.InsertReviewModel(newReview)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "리뷰를 추가하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, reviewResult)
	return
}

// UpdateReviewControl godoc
//
//	@Summary		call UpdateReviewControl, return result by json.
//	@Description	review data 수정을 위한 기능.
//	@name			UpdateReviewControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			menu	body	model.Review	true	"{orderDay, orderId, star, content}"
//	@Router			/customer/order/review [put]
//	@Success		200	{object}	string
func (p *Controller) UpdateReviewControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))

	if !p.UserValidation(c, ut.TypeCustomer, userId) {
		return
	}

	review := model.Review{}
	if !p.ReviewBind(c, &review) {
		return
	}
	review.UserId = userId

	if err := p.md.UpdateReviewModel(review); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "리뷰를 수정하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}

// DeleteReviewControl godoc
//
//	@Summary		call DeleteReviewControl, return result by json.
//	@Description	review data 삭제를 위한 기능.
//	@name			DeleteReviewControl
//	@Accept			json
//	@Produce		json
//	@Param			userId	header	string	true	"User ID"
//	@Param			orderDay	query	string	true	"Order Day"
//	@Param			orderId	query	int	true	"Order Id"
//	@Router			/customer/order/review [delete]
//	@Success		200	{object}	string
func (p *Controller) DeleteReviewControl(c *gin.Context) {

	userId, _ := strconv.Atoi(c.GetHeader("userId"))
	orderDay := c.Query("orderDay")
	orderId, _ := strconv.Atoi(c.Query("orderId"))
	if !p.UserValidation(c, ut.TypeCustomer, userId) {
		return
	}

	review := model.Review{}

	review.UserId = userId

	if orderDay == "" || orderId == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":  "삭제 정보가 잘못됬습니다!",
			"orderDay": orderDay,
			"orderId":  orderId,
		})
		return
	}

	review.OrderDay = orderDay
	review.OrderId = orderId
	if err := p.md.DeleteReviewModel(review); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "리뷰를 삭제하지 못했습니다!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"msg": "ok"})
	return
}
