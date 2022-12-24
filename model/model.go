package model

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	ot "WBABEProject-09/type/order"

	conf "WBABEProject-09/config"
	log "WBABEProject-09/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client        *mongo.Client
	userCol       *mongo.Collection
	orderCol      *mongo.Collection
	orderSaveCol  *mongo.Collection
	menuCol       *mongo.Collection
	reviewCol     *mongo.Collection
	reviewSaveCol *mongo.Collection
	idSeq         *mongo.Collection
	orderCountSeq *mongo.Collection
}

type User struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
	UserId   int                 `json:"userId" bson:"userId"`
	Name     string              `json:"name" bson:"name"`
	Email    string              `json:"email" bson:"email"`
	Phone    string              `json:"phone" bson:"phone"`
	Address  string              `json:"address" bson:"address"`
	Type     int                 `json:"type" bson:"type"`
	Use      bool                `json:"use" bson:"use"`
	CreateAt time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt time.Time           `json:"modifyAt" bson:"modifyAt"`
}

type Order struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
	UserId   int                 `json:"userId" bson:"userId"`
	OrderDay string              `json:"orderDay" bson:"orderDay"`
	OrderId  int                 `json:"orderId" bson:"orderId"`
	Menu     []OrderMenu         `json:"menu" bson:"menu"`
	Phone    string              `json:"phone" bson:"phone"`
	Address  string              `json:"address" bson:"address"`
	State    int                 `json:"state" bson:"state"`
	CreateAt time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt time.Time           `json:"modifyAt" bson:"modifyAt"`
}
type OrderMenu struct {
	MenuId int    `json:"menuId" bson:"menuId"`
	Name   string `json:"name" bson:"name"`
}

type Review struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
	UserId   int                 `json:"userId" bson:"userId"`
	OrderDay string              `json:"orderDay" bson:"orderDay"`
	OrderId  int                 `json:"orderId" bson:"orderId"`
	Star     float32             `json:"star" bson:"star"`
	Content  string              `json:"content" bson:"content"`
	CreateAt time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt time.Time           `json:"modifyAt" bson:"modifyAt"`
}

type Menu struct {
	Id              *primitive.ObjectID `bson:"_id,omitempty"`
	MenuId          int                 `json:"menuId" bson:"menuId"`
	Category        string              `json:"category" bson:"category"`
	Name            string              `json:"name" bson:"name"`
	Price           int                 `json:"price" bson:"price"`
	Recommend       bool                `json:"recommend" bson:"recommend"`
	Star            float32             `json:"star" bson:"star"`
	OrderState      int                 `json:"orderState" bson:"orderState"`
	OrderCount      int                 `json:"orderCount" bson:"orderCount"`
	OrderDailyLimit int                 `json:"orderDailyLimit" bson:"orderDailyLimit"`
	ReorderCount    int                 `json:"reorderCount" bson:"reorderCount"`
	Use             bool                `json:"use" bson:"use"`
	CreateAt        time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt        time.Time           `json:"modifyAt" bson:"modifyAt"`
}

func NewModel(cfg *conf.Config) (*Model, error) {
	r := &Model{}
	var err error
	mgUrl := cfg.DB.Host
	if r.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(mgUrl)); err != nil {
		return nil, err
	} else if err := r.client.Ping(context.Background(), nil); err != nil {
		return nil, err
	} else {
		db := r.client.Database(cfg.DB.DB)
		r.userCol = db.Collection(cfg.DB.UserCollection)
		r.orderCol = db.Collection(cfg.DB.OrderCollection)
		r.orderSaveCol = db.Collection(cfg.DB.OrderSaveCollection)
		r.menuCol = db.Collection(cfg.DB.MenuCollection)
		r.reviewCol = db.Collection(cfg.DB.ReviewCollection)
		r.reviewSaveCol = db.Collection(cfg.DB.ReviewSaveCollection)
		r.idSeq = db.Collection(cfg.DB.IdSequence)
		r.orderCountSeq = db.Collection(cfg.DB.OrderCountSequence)
	}

	return r, nil
}

func NewMenu() Menu {
	return Menu{
		OrderState:   1,
		OrderCount:   0,
		Star:         0,
		ReorderCount: 0,
		Use:          true,
		CreateAt:     time.Now(),
		ModifyAt:     time.Now(),
	}
}
func NewOrder() Order {
	return Order{
		State:    ot.StateReceiving,
		CreateAt: time.Now(),
		ModifyAt: time.Now(),
	}
}

func NewReview() Review {
	return Review{
		CreateAt: time.Now(),
		ModifyAt: time.Now(),
	}
}

// auto-increment ID를 활용하기 위해 만든 시퀀스
/*
단 조회 순간 바로 update해버리기 때문에 추후에 트랜잭션 등을 통해서 고려가 필요해보임.
조회 후 사용안한다고 해도 버그가 있는건 아니니 보류
*/
func (p *Model) GetAutoId(idType string) (int, error) {

	filter := bson.M{"_id": idType}
	update := bson.D{{"$inc", bson.D{{"seq", 1}}}}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	JSONData := struct {
		Seq int `json:"seq" bson:"seq"`
	}{}

	err := p.idSeq.FindOneAndUpdate(context.TODO(), filter, update, &opt).Decode(&JSONData)

	if err != nil {
		log.Error("Decode error: ", err)
		return 0, err
	}
	return JSONData.Seq, err
}

// 일마다 별도의 auto-increment ID를 산정하기 위한 기능
func (p *Model) GetOrderId(day string) (int, error) {

	filter := bson.M{"_id": day}
	update := bson.D{{"$inc", bson.D{{"seq", 1}}}}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	JSONData := struct {
		Seq int `json:"seq" bson:"seq"`
	}{}

	err := p.orderCountSeq.FindOneAndUpdate(context.TODO(), filter, update, &opt).Decode(&JSONData)

	if err != nil {
		log.Error("Decode error: ", err)
		return 0, err
	}
	return JSONData.Seq, err
}

// order 수정시 메뉴를 비교하기 위한 기능
//
// 0: 같음, 1: 변경, 2: 추가
func CompareMenu(oldOrder []OrderMenu, newOrder []OrderMenu) (map[int]string, int) {
	menuMap := make(map[int]string)
	for _, menuData := range oldOrder {
		menuMap[menuData.MenuId] = "old"
	}
	for _, menuData := range newOrder {
		if menuMap[menuData.MenuId] == "old" {
			menuMap[menuData.MenuId] = "same"
		} else {
			menuMap[menuData.MenuId] = "new"
		}
	}

	compareResult := 0
	for _, compareType := range menuMap {
		if compareType == "old" {
			return menuMap, 1
		} else if compareType == "new" {
			compareResult = 2
		}
	}
	return menuMap, compareResult
}
func (p *Model) GetUserTypeByIdModel(userId int) (*User, error) {

	var user User
	filter := bson.M{"userId": userId}
	project := bson.M{"type": 1}
	opts := options.FindOne().SetProjection(project)
	err := p.userCol.FindOne(context.TODO(), filter, opts).Decode(&user)

	if err != nil {
		log.Error("유저 확인 에러", err.Error())
		return nil, err
	}

	return &user, nil
}

// 오더 주문시 연계된 메뉴들에 대해서 재주문, 주문Count를 처리하기 위한 기능
//
// 트랜잭션 고려는 안되어있음
// 초기 Aggregate count 방식으로 설계되었으나 성능 낭비에 데이터 변환 로직이 복잡해 폐기
func (p *Model) CheckOrderMenuModel(orderData *Order) error {

	for _, orderMenuData := range orderData.Menu {

		filter := bson.M{"menuId": orderMenuData.MenuId}
		updateTarget := bson.D{}

		// 주문자가 일치하고 배달 완료로 오더가 완료된 메뉴를 대상으로 집계
		matchState := bson.D{
			{"userId", orderData.UserId},
			{"menu", bson.D{{"$elemMatch", bson.D{{"menuId", orderMenuData.MenuId}}}}},
			{"state", 7},
		}

		findOpt := options.FindOne()
		var findResult bson.M
		// 배달 완료로 별도에 orderSave 콜렉션에 저장된 과거 주문 내역을 참조
		findErr := p.orderSaveCol.FindOne(context.TODO(), matchState, findOpt).Decode(&findResult)
		if findErr == nil {
			updateTarget = append(updateTarget, bson.E{"reorderCount", 1})
		}
		updateTarget = append(updateTarget, bson.E{"orderCount", 1})
		update := bson.D{{"$inc", updateTarget}}
		result := p.menuCol.FindOneAndUpdate(context.TODO(), filter, update)
		err := result.Err()
		if err != nil {
			log.Error(err.Error())
		}

	}
	return nil
}

// 리뷰 추가시 연계된 메뉴들에 대해서 계산해 업데이트 위한 함수
func (p *Model) UpdateMenuReviewStarModel(orderData *Order) error {

	// 초기 구상은 3개 테이블을 Join하는 방식이었지만, 복잡도와 성능 문제로 폐기
	// 완료된 order와 review를 연결하고 menu ID를 기준으로 포함되는 order 정보로 평균을 내는 방식
	// ERD 구성 및 설계단계부터 잘못된 느낌
	for _, orderMenuData := range orderData.Menu {
		pipeline := mongo.Pipeline{
			{
				{"$lookup", bson.D{
					{"from", "tReview"},
					{"let",
						bson.M{"order_day": "$orderDay", "order_id": "$orderId"},
					},
					{"pipeline", bson.A{bson.D{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", []interface{}{
									bson.M{"$eq": []string{"$orderDay", "$$order_day"}},
									bson.M{"$eq": []string{"$orderId", "$$order_id"}},
								}},
							}},
						}},
					}}},
					{"as", "orderReview"},
				}},
			},
			{{"$unwind", bson.D{{"path", "$orderReview"}}}},

			{{"$match", bson.D{{"state", 7}, {"menu.menuId", orderMenuData.MenuId}}}},
			{{"$group", bson.D{
				{"_id", orderMenuData.MenuId},
				{"avgStar", bson.M{"$avg": "$orderReview.star"}},
			},
			}},
		}
		cursor, _ := p.orderSaveCol.Aggregate(context.TODO(), pipeline)
		if cursor.TryNext(context.TODO()) {

			// aggregate 결과물에서 avg를 추출하는 과정
			// 더 좋은 구조가 있을꺼 같지만, 시간이 너무 소요되서 일단 동작하니 pass
			type avgResult struct {
				AvgStar float64 `bson:"avgStar"`
			}

			bsonResult := avgResult{}
			cursor.Decode(&bsonResult)

			updateStar := math.Round(bsonResult.AvgStar*10) / 10
			if updateStar != 0 {

				updateFilter := bson.M{"menuId": orderMenuData.MenuId}
				update := bson.D{{"$set", bson.D{{"star", updateStar}}}}
				_, err := p.menuCol.UpdateOne(context.TODO(), updateFilter, update)
				if err != nil {
					log.Error("메뉴 star 수정 에러", err.Error())
				}
			}

		}
	}
	return nil
}

func (p *Model) InsertMenuModel(menuData Menu) (*Menu, error) {
	menuId, err := p.GetAutoId("menuId")
	menuData.MenuId = menuId
	res, err := p.menuCol.InsertOne(context.TODO(), menuData)

	if err != nil {
		log.Error("메뉴 추가 에러", err.Error())
	}

	var newMenu Menu
	query := bson.M{"_id": res.InsertedID}
	if err = p.menuCol.FindOne(context.TODO(), query).Decode(&newMenu); err != nil {
		log.Error("메뉴 추가 후 조회 에러", err.Error())
		return nil, err
	}
	return &newMenu, err
}

func (p *Model) UpdateMenuModel(menuId int, menuData Menu) error {

	var oldMenu Menu
	findFilter := bson.M{"menuId": menuId}
	if err := p.menuCol.FindOne(context.TODO(), findFilter).Decode(&oldMenu); err != nil {
		log.Error("메뉴 조회 에러", err.Error())
		return err
	}

	updateTarget := bson.D{}
	switch {
	case menuData.Category != oldMenu.Category:
		updateTarget = append(updateTarget, bson.E{"category", menuData.Category})
		fallthrough
	case menuData.Name != oldMenu.Name:
		updateTarget = append(updateTarget, bson.E{"name", menuData.Name})
		fallthrough
	case menuData.Price != oldMenu.Price:
		updateTarget = append(updateTarget, bson.E{"price", menuData.Price})
		fallthrough
	case menuData.Recommend != oldMenu.Recommend:
		updateTarget = append(updateTarget, bson.E{"recommend", menuData.Recommend})
		fallthrough
	case menuData.OrderState != oldMenu.OrderState:
		updateTarget = append(updateTarget, bson.E{"orderState", menuData.OrderState})
		fallthrough
	case menuData.OrderDailyLimit != oldMenu.OrderDailyLimit:
		updateTarget = append(updateTarget, bson.E{"orderDailyLimit", menuData.OrderDailyLimit})
		fallthrough
	default:
		updateTarget = append(updateTarget, bson.E{"modifyAt", time.Now()})
	}

	updateFilter := bson.M{"menuId": menuId}
	update := bson.D{{"$set", updateTarget}}
	_, err := p.menuCol.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Error("메뉴 수정 에러", err.Error())
	}

	return err
}

func (p *Model) DeleteMenuModel(menuId int) error {

	var oldMenu Menu
	findFilter := bson.M{"menuId": menuId}
	if err := p.menuCol.FindOne(context.TODO(), findFilter).Decode(&oldMenu); err != nil {
		log.Error("메뉴 조회 에러", err)
		return err
	}

	if oldMenu.Use == false {
		log.Info("이미 삭제된 메뉴")
		return errors.New("이미 삭제된 메뉴")
	}

	filter := bson.M{"menuId": menuId}
	delete := bson.D{{"$set", bson.D{{"use", false}}}}
	_, err := p.menuCol.UpdateOne(context.TODO(), filter, delete)

	if err != nil {
		log.Error("메뉴 삭제 에러", err)
	}

	return err
}

func (p *Model) InsertOrderModel(orderData Order) (*Order, error) {

	now := time.Now().UTC()
	day := now.Format("2006-01-02")
	orderId, err := p.GetOrderId(day)

	orderData.OrderDay = day
	orderData.OrderId = orderId

	res, err := p.orderCol.InsertOne(context.TODO(), orderData)

	if err != nil {
		log.Error("오더 추가 에러", err.Error())
		return &orderData, err
	}

	_, err = p.orderSaveCol.InsertOne(context.TODO(), orderData)

	if err != nil {
		log.Error("초기 오더 저장 에러", err.Error())
	}
	err = p.CheckOrderMenuModel(&orderData)

	if err != nil {
		log.Error("오더 메뉴 정산 에러", err.Error())
	}

	var newOrder Order
	query := bson.M{"_id": res.InsertedID}
	if err = p.orderCol.FindOne(context.TODO(), query).Decode(&newOrder); err != nil {
		log.Error("오더 추가 후 조회 에러", err.Error())
		return nil, err
	}
	return &newOrder, err
}

func (p *Model) UpdateCustomerOrderModel(orderData Order) error {

	var oldOrder Order
	findFilter := bson.M{"userId": orderData.UserId, "orderDay": orderData.OrderDay, "orderId": orderData.OrderId}
	if err := p.orderCol.FindOne(context.TODO(), findFilter).Decode(&oldOrder); err != nil {
		log.Error("오더 조회 에러", err.Error())
		return err
	}

	// -TODO- 리턴된 비교map을 활용해 취소된 메뉴와 추가된 메뉴에 대한 menu count 증가 로직이 필요
	// 시간 배분을 고려해 나중에 작업
	_, compareResult := CompareMenu(oldOrder.Menu, orderData.Menu)

	// 배달중 이상의 상태에서는 오더 추가가 불가능
	if compareResult == 2 && orderData.State >= ot.StateInDelivery {
		log.Error("오더 상태 변경 에러")
		errorMsg := fmt.Sprintf("오더를 추가할 수 없습니다. 현재상태: %s", ot.GetOrderStateText(orderData.State))
		return errors.New(errorMsg)
	} else if compareResult == 1 && orderData.State >= ot.StateCooking {
		// 조리중 이상의 상태에서는 오더 변경이 불가능
		log.Error("오더 상태 변경 에러")
		errorMsg := fmt.Sprintf("오더를 변경할 수 없습니다. 현재상태: %s", ot.GetOrderStateText(orderData.State))
		return errors.New(errorMsg)

	}

	updateTarget := bson.D{}
	switch {
	case orderData.Phone != oldOrder.Phone:
		updateTarget = append(updateTarget, bson.E{"phone", orderData.Phone})
		fallthrough
	case orderData.Address != oldOrder.Address:
		updateTarget = append(updateTarget, bson.E{"address", orderData.Address})
		fallthrough
	case compareResult != 0:
		updateTarget = append(updateTarget, bson.E{"menu", orderData.Menu})
		fallthrough
	default:
		updateTarget = append(updateTarget, bson.E{"modifyAt", time.Now()})
	}

	updateFilter := bson.M{"userId": orderData.UserId, "orderDay": orderData.OrderDay, "orderId": orderData.OrderId}
	update := bson.D{{"$set", updateTarget}}
	_, err := p.orderCol.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Error("오더 수정 에러", err.Error())
	}
	return nil
}

func (p *Model) UpdateOwnerOrderModel(orderData Order) error {

	var oldOrder Order
	findFilter := bson.M{"orderDay": orderData.OrderDay, "orderId": orderData.OrderId}
	if err := p.orderCol.FindOne(context.TODO(), findFilter).Decode(&oldOrder); err != nil {
		log.Error("오더 조회 에러", err.Error())
		return err
	}
	switch {
	// 취소 또는 배달 완료의 경우 기존 Order 콜렉션에서 제거 후 Save로 이동
	case orderData.State == 3 || orderData.State == 7:
		oldOrder.State = orderData.State

		if _, err := p.orderSaveCol.InsertOne(context.TODO(), oldOrder); err != nil {
			log.Error("오더 저장 에러", err.Error())
			return err
		}
		if _, err := p.orderCol.DeleteOne(context.TODO(), findFilter); err != nil {
			log.Error("오더 백업 에러", err.Error())
			return err
		}
		return nil
	case orderData.State == 1:
		// 오더를 접수중으로 변경할수는 없음
		log.Error("오더 상태 변경 에러!")
		return errors.New("접수중으로 변경할 순 없습니다!")
	}

	updateFilter := bson.M{"orderDay": orderData.OrderDay, "orderId": orderData.OrderId}
	update := bson.D{{"$set", bson.D{{"state", orderData.State}}}}
	_, err := p.orderCol.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Error("오더 상태 수정 에러", err.Error())
	}
	return nil
}

func (p *Model) InsertReviewModel(reviewData Review) (*Review, error) {

	var targetOrder Order
	orderFindFilter := bson.M{"userId": reviewData.UserId, "orderDay": reviewData.OrderDay, "orderId": reviewData.OrderId, "state": 7}
	if err := p.orderSaveCol.FindOne(context.TODO(), orderFindFilter).Decode(&targetOrder); err != nil {
		log.Error("리뷰 타겟 오더 조회 에러", err.Error())
		return nil, err
	}
	var oldReview Review
	reviewFindFilter := bson.M{"userId": reviewData.UserId, "orderDay": reviewData.OrderDay, "orderId": reviewData.OrderId}
	if err := p.reviewCol.FindOne(context.TODO(), reviewFindFilter).Decode(&oldReview); err == nil {
		log.Error("리뷰 저장 에러: 이미 리뷰가 존재")
		return nil, errors.New("이미 리뷰가 존재함")
	}

	res, err := p.reviewCol.InsertOne(context.TODO(), reviewData)

	if err != nil {
		log.Error("리뷰 추가 에러", err.Error())
		return nil, err
	}

	_, err = p.reviewSaveCol.InsertOne(context.TODO(), reviewData)

	if err != nil {
		log.Error("초기 리뷰 저장 에러", err.Error())
	}
	err = p.UpdateMenuReviewStarModel(&targetOrder)

	var newReview Review
	query := bson.M{"_id": res.InsertedID}
	if err = p.reviewCol.FindOne(context.TODO(), query).Decode(&newReview); err != nil {
		log.Error("리뷰 추가 후 조회 에러", err.Error())
		return nil, err
	}
	return &newReview, err
}
