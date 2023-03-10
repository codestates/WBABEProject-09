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

// 2d1c3bec4f68afa657d5c39ec98e2b6289e1ee19 commit에서 설계를 변경
// Order에 포함되던 Review를 콜렉션 상으로 서로 분리시킴
// 다시 생각해보니 설계적으로 실수함 추후 다시 고려- TODO -
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

// star에 binding validation 적용
// - TODO - 각 struct 필수 필드에 validation 적용
type Review struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
	UserId   int                 `json:"userId" bson:"userId"`
	OrderDay string              `json:"orderDay" bson:"orderDay"`
	OrderId  int                 `json:"orderId" bson:"orderId"`
	Star     float64             `json:"star" bson:"star" binding:"gte=0,lte=5"`
	Content  string              `json:"content" bson:"content"`
	CreateAt time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt time.Time           `json:"modifyAt" bson:"modifyAt"`
}

// 타입 변동사항은 이슈 #4 참조
type Menu struct {
	Id              *primitive.ObjectID `bson:"_id,omitempty"`
	MenuId          int                 `json:"menuId" bson:"menuId"`
	Category        string              `json:"category" bson:"category"`
	Name            string              `json:"name" bson:"name"`
	Price           int                 `json:"price" bson:"price"`
	Recommend       bool                `json:"recommend" bson:"recommend"`
	Star            float64             `json:"star" bson:"star"`
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

func NewUser(user *User) {
	var dateNow = time.Now()
	user.Use = true
	user.CreateAt = dateNow
	user.ModifyAt = dateNow
}
func NewMenu(menu *Menu) {
	var dateNow = time.Now()
	menu.OrderState = 1
	menu.OrderCount = 0
	menu.Star = 0
	menu.ReorderCount = 0
	menu.Use = true
	menu.CreateAt = dateNow
	menu.ModifyAt = dateNow
}
func NewOrder(order *Order) {
	var dateNow = time.Now()
	order.State = ot.StateReceiving
	order.CreateAt = dateNow
	order.ModifyAt = dateNow
}

func NewReview(review *Review) {
	var dateNow = time.Now()
	review.CreateAt = dateNow
	review.ModifyAt = dateNow
}

// auto-increment ID를 활용하기 위해 만든 시퀀스
/*
단 조회 순간 바로 update해버리기 때문에 추후에 트랜잭션 등을 통해서 고려가 필요해보임.
조회 후 사용안한다고 해도 버그가 있는건 아니니 보류
*/
func (p *Model) GetAutoId(idType string) (int, error) {

	filter := bson.M{"_id": idType}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "seq", Value: 1},
		}},
	}
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
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "seq", Value: 1},
		}},
	}
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
			{Key: "userId", Value: orderData.UserId},
			{Key: "menu", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "menuId", Value: orderMenuData.MenuId},
				}},
			}},
			{Key: "state", Value: 7},
		}

		findOpt := options.FindOne()
		var findResult bson.M
		// 배달 완료로 별도에 orderSave 콜렉션에 저장된 과거 주문 내역을 참조
		findErr := p.orderSaveCol.FindOne(context.TODO(), matchState, findOpt).Decode(&findResult)
		if findErr == nil {
			updateTarget = append(updateTarget, bson.E{Key: "reorderCount", Value: 1})
		}
		updateTarget = append(updateTarget, bson.E{Key: "orderCount", Value: 1})
		update := bson.D{{Key: "$inc", Value: updateTarget}}
		result := p.menuCol.FindOneAndUpdate(context.TODO(), filter, update)
		err := result.Err()
		if err != nil {
			log.Error(err.Error())
		}

	}
	return nil
}

// 상단에 주문한 메뉴를 검사하는 기능이 존재하지만, 검증은 주문전에 해야함
// 메뉴 상태에 영향을 미치는 로직은 주문 후에 진행되야함
// 현재는 반복문이 2번 돌더라도, 직관적으로 진행
func (p *Model) CheckOrderMenuDeletedModel(orderData *Order) error {

	for _, orderMenuData := range orderData.Menu {

		filter := bson.M{"menuId": orderMenuData.MenuId, "use": true}
		findOpt := options.FindOne()
		var findResult bson.M
		// 배달 완료로 별도에 orderSave 콜렉션에 저장된 과거 주문 내역을 참조
		err := p.menuCol.FindOne(context.TODO(), filter, findOpt).Decode(&findResult)
		if err != nil {
			log.Error(err.Error())
			return err
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
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "tReview"},
					{Key: "let", Value: bson.M{
						"order_day": "$orderDay",
						"order_id":  "$orderId",
					}},
					{Key: "pipeline", Value: bson.A{bson.D{
						{Key: "$match", Value: bson.D{
							{Key: "$expr", Value: bson.D{
								{Key: "$and", Value: []interface{}{
									bson.M{"$eq": []string{"$orderDay", "$$order_day"}},
									bson.M{"$eq": []string{"$orderId", "$$order_id"}},
								}},
							}},
						}},
					}}},
					{Key: "as", Value: "orderReview"},
				}},
			},
			{
				{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$orderReview"},
				}},
			},

			{
				{Key: "$match", Value: bson.D{
					{Key: "state", Value: 7},
					{Key: "menu.menuId", Value: orderMenuData.MenuId},
				}},
			},
			{
				{Key: "$group", Value: bson.D{
					{Key: "_id", Value: orderMenuData.MenuId},
					{Key: "avgStar", Value: bson.M{"$avg": "$orderReview.star"}},
				}},
			},
		}
		cursor, _ := p.orderSaveCol.Aggregate(context.TODO(), pipeline)
		if cursor.TryNext(context.TODO()) {

			// aggregate 결과물에서 avg를 추출하는 과정
			// 더 좋은 구조가 있을꺼 같지만, 시간이 너무 소요되서 일단 동작하니 pass
			// - TODO - 컨버터 고려하기
			type avgResult struct {
				AvgStar float64 `bson:"avgStar"`
			}

			bsonResult := avgResult{}
			cursor.Decode(&bsonResult)

			updateStar := math.Round(bsonResult.AvgStar*10) / 10
			if updateStar != 0 {

				updateFilter := bson.M{"menuId": orderMenuData.MenuId}
				update := bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "star", Value: updateStar},
					}},
				}
				_, err := p.menuCol.UpdateOne(context.TODO(), updateFilter, update)
				if err != nil {
					log.Error("메뉴 star 수정 에러", err.Error())
				}
			}

		}
	}
	return nil
}

func (p *Model) InsertUserModel(userData User) (*User, error) {

	var oldUser User
	userFindFilter := bson.M{"userId": userData.UserId}
	if err := p.userCol.FindOne(context.TODO(), userFindFilter).Decode(&oldUser); err == nil {
		log.Error("이미 유저가 존재합니다.")
		return nil, errors.New("이미 존재하는 유저")
	}

	res, err := p.userCol.InsertOne(context.TODO(), userData)

	if err != nil {
		log.Error("유저 추가 에러", err.Error())
	}

	var newUser User
	query := bson.M{"_id": res.InsertedID}
	if err = p.userCol.FindOne(context.TODO(), query).Decode(&newUser); err != nil {
		log.Error("유저 추가 후 조회 에러", err.Error())
		return nil, err
	}
	return &newUser, err
}
func (p *Model) GetMenuModel(sortBy string, checkReview int) ([]primitive.M, error) {

	var newMenu []primitive.M
	// checkReview가 0이 아니면, Review 데이터를 결과에 포함시킴
	if checkReview != 0 {
		// 설계를 유지하는 대신 복잡한 Join을 사용
		// join 과정에서 배열속에 있는 key를 사용해 다른 컬렉션에 다시 join을 해야했기에 unwind를 사용
		// 문제는 unwind를 사용하면서 오더나 리뷰가 없는 메뉴들에 경우엔 문서 자체가 제외되서 표시됨
		// 즉, 리뷰가 존재하는 메뉴만 보이게됨
		// 다른 방법이 있을꺼라 생각하지만, 시간이 없으니 일단 보류
		// order에 들어가는 메뉴가 배열이기에 생기는 문제로, 설계 난이도를 낮추면 해결 가능할꺼라 판단됨
		pipeline := mongo.Pipeline{
			{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "tOrderSave"},
					{Key: "localField", Value: "menuId"},
					{Key: "foreignField", Value: "menu.menuId"},
					{Key: "pipeline", Value: bson.A{bson.D{
						{Key: "$match", Value: bson.D{
							{Key: "$expr", Value: bson.D{
								{Key: "$and", Value: []interface{}{
									bson.M{"$eq": bson.A{"$state", 7}},
								}},
							}},
						}},
					}}},
					{Key: "as", Value: "menuOrder"},
				}},
			},
			{
				{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$menuOrder"},
				}},
			},
			{
				{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "tReview"},
					{Key: "let", Value: bson.M{
						"order_day": "$menuOrder.orderDay",
						"order_id":  "$menuOrder.orderId",
					}},
					{Key: "pipeline", Value: bson.A{bson.D{
						{Key: "$match", Value: bson.D{
							{Key: "$expr", Value: bson.D{
								{Key: "$and", Value: []interface{}{
									bson.M{"$eq": []string{"$orderDay", "$$order_day"}},
									bson.M{"$eq": []string{"$orderId", "$$order_id"}},
								}},
							}},
						}},
					}}},
					{Key: "as", Value: "orderReview"},
				}},
			},
			{
				{Key: "$unwind", Value: bson.D{
					{Key: "path", Value: "$orderReview"},
				}},
			},

			{
				{Key: "$match", Value: bson.D{
					{Key: "use", Value: true},
				}},
			},
			{
				{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$_id"},
					{Key: "menuId", Value: bson.M{"$first": "$menuId"}},
					{Key: "category", Value: bson.M{"$first": "$category"}},
					{Key: "name", Value: bson.M{"$first": "$name"}},
					{Key: "price", Value: bson.M{"$first": "$price"}},
					{Key: "recommend", Value: bson.M{"$first": "$recommend"}},
					{Key: "star", Value: bson.M{"$first": "$star"}},
					{Key: "orderState", Value: bson.M{"$first": "$orderState"}},
					{Key: "orderCount", Value: bson.M{"$first": "$orderCount"}},
					{Key: "reorderCount", Value: bson.M{"$first": "$reorderCount"}},
					{Key: "createAt", Value: bson.M{"$first": "$createAt"}},
					{Key: "review", Value: bson.M{"$push": bson.M{"userId": "$orderReview.userId",
						"orderDay": "$orderReview.orderDay",
						"orderId":  "$orderReview.orderId",
						"star":     "$orderReview.star",
						"content":  "$orderReview.content",
						"createAt": "$orderReview.createAt",
					}}},
				}},
			},
			{
				{Key: "$sort", Value: bson.D{
					{Key: sortBy, Value: -1},
				}},
			},
		}
		cursor, err := p.menuCol.Aggregate(context.TODO(), pipeline)
		if err != nil {
			log.Error("메뉴 조회 에러", err.Error())
			return nil, err
		}
		if err = cursor.All(context.TODO(), &newMenu); err != nil {
			log.Error("메뉴 조회 에러", err.Error())
			return nil, err
		}
	} else {
		// 삭제된 메뉴에 경우 표시하지 않음
		filter := bson.M{"use": true}
		// switch 문으로 분기 처리를 했었지만, 잘못된 데이터가 들어와도 기본 정렬이 됨
		// sortBy에 제대로만 입력했으면 정렬은 맞춰서 되도록 SetSort 1줄로 줄임
		findOptions := options.Find().SetSort(bson.D{{Key: sortBy, Value: -1}})
		cursor, err := p.menuCol.Find(context.TODO(), filter, findOptions)
		if err != nil {
			log.Error("메뉴 조회 에러", err.Error())
			return nil, err
		}
		if err = cursor.All(context.TODO(), &newMenu); err != nil {
			log.Error("메뉴 조회 에러", err.Error())
			return nil, err
		}

	}
	return newMenu, nil
}

func (p *Model) GetMenuDetailModel(meniId int) ([]bson.M, error) {

	var bsonResult []bson.M

	pipeline := mongo.Pipeline{
		{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "tOrderSave"},
				{Key: "let",
					Value: bson.M{"order_day": "$orderDay", "order_id": "$orderId"},
				},
				{Key: "pipeline", Value: bson.A{bson.D{
					{Key: "$match", Value: bson.D{
						{Key: "$expr", Value: bson.D{
							{Key: "$and", Value: []interface{}{
								bson.M{"$eq": []string{"$orderDay", "$$order_day"}},
								bson.M{"$eq": []string{"$orderId", "$$order_id"}},
								bson.M{"$eq": bson.A{"$state", 7}},
							}},
						}},
					}},
				}}},
				{Key: "as", Value: "orderReview"},
			}},
		},
		{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$orderReview"}}}},

		{{Key: "$match", Value: bson.D{{Key: "orderReview.menu.menuId", Value: meniId}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "userId", Value: 1},
			{Key: "star", Value: 1},
			{Key: "content", Value: 1},
		},
		}},
	}
	cursor, err := p.reviewCol.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Error("메뉴 상세 조회 에러: ", err.Error())
		return nil, err
	}

	if err = cursor.All(context.TODO(), &bsonResult); err != nil {
		log.Error("메뉴 상세 조회 에러: ", err.Error())
		return nil, err
	}
	return bsonResult, nil
}
func (p *Model) InsertMenuModel(menuData Menu) (*Menu, error) {

	var oldMenu Menu
	menuFindFilter := bson.M{"category": menuData.Category, "name": menuData.Name, "use": true}
	if err := p.menuCol.FindOne(context.TODO(), menuFindFilter).Decode(&oldMenu); err == nil {
		log.Error("이미 메뉴가 존재합니다.")
		return nil, errors.New("이미 존재하는 메뉴")
	}

	// 재발급 대신 일단 중복만 체크하고 pass
	// 시간이 있다면 중복된 ID 대신 재발급 로직 추가 - TODO -
	menuId, err := p.GetAutoId("menuId")
	menuIdFilter := bson.M{"menuId": menuId}
	if err := p.menuCol.FindOne(context.TODO(), menuIdFilter).Decode(&oldMenu); err == nil {
		log.Error("메뉴ID가 중복됬습니다.")
		return nil, errors.New("메뉴ID 중복")
	}

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
	findFilter := bson.M{"menuId": menuId, "use": true}
	if err := p.menuCol.FindOne(context.TODO(), findFilter).Decode(&oldMenu); err != nil {
		log.Error("메뉴 조회 에러", err.Error())
		return err
	}

	// 업데이트 대상을 변경사항으로 구분해 append 하는 방식
	// 좋은 구성인지는 잘 모르겠음
	updateTarget := bson.D{}
	switch {
	case menuData.Category != oldMenu.Category:
		updateTarget = append(updateTarget, bson.E{Key: "category", Value: menuData.Category})
		fallthrough
	case menuData.Name != oldMenu.Name:
		updateTarget = append(updateTarget, bson.E{Key: "name", Value: menuData.Name})
		fallthrough
	case menuData.Price != oldMenu.Price:
		updateTarget = append(updateTarget, bson.E{Key: "price", Value: menuData.Price})
		fallthrough
	case menuData.Recommend != oldMenu.Recommend:
		updateTarget = append(updateTarget, bson.E{Key: "recommend", Value: menuData.Recommend})
		fallthrough
	case menuData.OrderState != oldMenu.OrderState:
		updateTarget = append(updateTarget, bson.E{Key: "orderState", Value: menuData.OrderState})
		fallthrough
	case menuData.OrderDailyLimit != oldMenu.OrderDailyLimit:
		updateTarget = append(updateTarget, bson.E{Key: "orderDailyLimit", Value: menuData.OrderDailyLimit})
		fallthrough
	default:
		updateTarget = append(updateTarget, bson.E{Key: "modifyAt", Value: time.Now()})
	}

	updateFilter := bson.M{"menuId": menuId}
	update := bson.D{{Key: "$set", Value: updateTarget}}
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
	delete := bson.D{{Key: "$set", Value: bson.D{{Key: "use", Value: false}}}}
	_, err := p.menuCol.UpdateOne(context.TODO(), filter, delete)

	if err != nil {
		log.Error("메뉴 삭제 에러", err)
	}

	return err
}

func (p *Model) GetInOrderModel(userId int, userType int) (*[]bson.M, error) {

	var orderList []bson.M
	filter := bson.M{}
	if userType == 2 {
		filter = bson.M{"userId": userId}
	}

	findOptions := options.Find().SetSort(
		bson.D{{Key: "state", Value: 1}},
	).SetProjection(
		bson.D{{Key: "_id", Value: 0}},
	)
	cursor, err := p.orderCol.Find(context.TODO(), filter, findOptions)
	if err == mongo.ErrNoDocuments {
		log.Error("조회 결과 없음", err.Error())
		return nil, nil
	} else if err != nil {
		log.Error("오더 조회 에러", err.Error())
		return nil, err
	}
	if err = cursor.All(context.TODO(), &orderList); err != nil {
		log.Error("오더 조회 에러", err.Error())
		return nil, err
	}

	return &orderList, err
}

func (p *Model) GetDoneOrderModel(userId int, userType int) (*[]bson.M, error) {

	var orderList []bson.M
	filter := bson.D{bson.E{Key: "state", Value: ot.StateDelivered}}
	if userType == 2 {
		filter = append(filter, bson.E{Key: "userId", Value: userId})
	}

	findOptions := options.Find().SetSort(
		bson.D{{Key: "createAt", Value: -1}},
	).SetProjection(
		bson.D{{Key: "_id", Value: 0}},
	)
	cursor, err := p.orderSaveCol.Find(context.TODO(), filter, findOptions)
	if err == mongo.ErrNoDocuments {
		log.Error("조회 결과 없음", err.Error())
		return nil, nil
	} else if err != nil {
		log.Error("오더 조회 에러", err.Error())
		return nil, err
	}
	if err = cursor.All(context.TODO(), &orderList); err != nil {
		log.Error("오더 조회 에러", err.Error())
		return nil, err
	}

	return &orderList, err
}
func (p *Model) InsertOrderModel(orderData Order) (*Order, error) {

	err := p.CheckOrderMenuDeletedModel(&orderData)
	if err != nil {
		log.Error("오더 추가 에러, 메뉴 조회 실패", err.Error())
		return &orderData, err
	}

	now := time.Now()
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

	// 오더에 삭제된 메뉴 검증
	err := p.CheckOrderMenuDeletedModel(&orderData)
	if err != nil {
		log.Error("오더 수정 에러, 메뉴 조회 실패", err.Error())
		return err
	}

	var oldOrder Order
	findFilter := bson.M{"userId": orderData.UserId, "orderDay": orderData.OrderDay, "orderId": orderData.OrderId}
	if err := p.orderCol.FindOne(context.TODO(), findFilter).Decode(&oldOrder); err != nil {
		// 실패 경우 완료된 경우일 수 있으니 orderSave에서 다시 조회
		findSaveFilter := bson.M{"userId": orderData.UserId, "orderDay": orderData.OrderDay, "orderId": orderData.OrderId, "state": 7}
		if err := p.orderSaveCol.FindOne(context.TODO(), findSaveFilter).Decode(&oldOrder); err != nil {
			log.Error("오더 조회 에러", err.Error())
			return err
		}
	}

	// - TODO - 리턴된 비교map을 활용해 취소된 메뉴와 추가된 메뉴에 대한 menu count 증가 로직이 필요
	// 시간 배분을 고려해 나중에 작업
	menuMap, compareResult := CompareMenu(oldOrder.Menu, orderData.Menu)

	// 배달중 이상의 상태에서는 오더 추가가 불가능
	if compareResult == 2 && oldOrder.State >= ot.StateInDelivery {
		// 추가불가능 경우 리턴된 비교map을 활용해 기존 order외 추가건은 신규 오더로 전환
		log.Error("오더 추가 불가 신규오더로 변경")

		orderMenu := []OrderMenu{}
		for _, menu := range orderData.Menu {
			// 기존 메뉴 검사에서 추가 메뉴인 경우에만 신규 오더로 넣음
			if menuMap[menu.MenuId] == "new" {
				orderMenu = append(orderMenu, menu)
			}
		}

		now := time.Now()
		day := now.Format("2006-01-02")
		orderId, _ := p.GetOrderId(day)

		// 기존 오더 정보에서 일부 정보를 그대로 가져와 필요 부분만 수정해 신규로 insert
		orderData.OrderDay = day
		orderData.OrderId = orderId
		orderData.Menu = orderMenu
		orderData.State = ot.StateReceiving
		orderData.CreateAt = now
		orderData.ModifyAt = now

		_, err := p.orderCol.InsertOne(context.TODO(), orderData)

		if err != nil {
			log.Error("오더 추가 에러", err.Error())
			return err
		}

		_, err = p.orderSaveCol.InsertOne(context.TODO(), orderData)

		if err != nil {
			log.Error("초기 오더 저장 에러", err.Error())
		}
		err = p.CheckOrderMenuModel(&orderData)

		if err != nil {
			log.Error("오더 메뉴 정산 에러", err.Error())
		}

		// 이후 수정 로직은 사용하지 않음으로 여기서 return
		return nil

	} else if compareResult == 1 && oldOrder.State >= ot.StateCooking {
		// 조리중 이상의 상태에서는 오더 변경이 불가능
		log.Error("오더 상태 변경 에러")
		errorMsg := fmt.Sprintf("오더를 변경할 수 없습니다. 현재상태: %s", ot.GetOrderStateText(orderData.State))
		return errors.New(errorMsg)

	}

	updateTarget := bson.D{}
	switch {
	case orderData.Phone != oldOrder.Phone:
		updateTarget = append(updateTarget, bson.E{Key: "phone", Value: orderData.Phone})
		fallthrough
	case orderData.Address != oldOrder.Address:
		updateTarget = append(updateTarget, bson.E{Key: "address", Value: orderData.Address})
		fallthrough
	case compareResult != 0:
		updateTarget = append(updateTarget, bson.E{Key: "menu", Value: orderData.Menu})
		fallthrough
	default:
		updateTarget = append(updateTarget, bson.E{Key: "modifyAt", Value: time.Now()})
	}

	updateFilter := bson.M{"userId": orderData.UserId, "orderDay": orderData.OrderDay, "orderId": orderData.OrderId}
	update := bson.D{{Key: "$set", Value: updateTarget}}
	_, err = p.orderCol.UpdateOne(context.TODO(), updateFilter, update)
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

		oldOrder.Id = nil
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
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "state", Value: orderData.State},
		}},
	}
	_, err := p.orderCol.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Error("오더 상태 수정 에러", err.Error())
	}
	return nil
}

func (p *Model) GetReviewModel(userId int, sortBy string) ([]Review, error) {

	var reviewList []Review

	filter := bson.M{"userId": userId}
	findOptions := options.Find().SetSort(bson.D{{Key: sortBy, Value: -1}})
	cursor, err := p.reviewCol.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Error("리뷰 조회 에러", err.Error())
		return nil, err
	}
	if err = cursor.All(context.TODO(), &reviewList); err != nil {
		log.Error("리뷰 조회 에러", err.Error())
		return nil, err
	}
	return reviewList, err
}

func (p *Model) InsertReviewModel(reviewData Review) (*Review, error) {

	var targetOrder Order
	orderFindFilter := bson.M{
		"userId":   reviewData.UserId,
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
		"state":    7,
	}
	if err := p.orderSaveCol.FindOne(context.TODO(), orderFindFilter).Decode(&targetOrder); err != nil {
		log.Error("리뷰 타겟 오더 조회 에러", err.Error())
		return nil, err
	}
	var oldReview Review
	reviewFindFilter := bson.M{
		"userId":   reviewData.UserId,
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
	}
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
	// 리뷰 추가 후 메뉴에 평점 반영
	// 중요 동작이 아니고 내부적으로 실패시 Log를 남겨서 err에 대한 별도 처리는 없음
	err = p.UpdateMenuReviewStarModel(&targetOrder)

	var newReview Review
	query := bson.M{"_id": res.InsertedID}
	if err = p.reviewCol.FindOne(context.TODO(), query).Decode(&newReview); err != nil {
		log.Error("리뷰 추가 후 조회 에러", err.Error())
		return nil, err
	}
	return &newReview, err
}

func (p *Model) UpdateReviewModel(reviewData Review) error {

	var targetOrder Order
	findOrderFilter := bson.M{
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
		"userId":   reviewData.UserId,
		"state":    7,
	}
	if err := p.orderSaveCol.FindOne(context.TODO(), findOrderFilter).Decode(&targetOrder); err != nil {
		log.Error("오더 조회 에러", err.Error())
		return err
	}

	var oldReview Review
	findReviewFilter := bson.M{
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
		"userId":   reviewData.UserId,
	}
	if err := p.reviewCol.FindOne(context.TODO(), findReviewFilter).Decode(&oldReview); err != nil {
		log.Error("리뷰 조회 에러", err.Error())
		return err
	}
	oldReview.Id = nil
	if _, err := p.reviewSaveCol.InsertOne(context.TODO(), oldReview); err != nil {
		log.Error("리뷰 백업 에러", err.Error())
		return err
	}

	updateTarget := bson.D{}

	switch {
	case reviewData.Star != oldReview.Star:
		updateTarget = append(updateTarget, bson.E{Key: "star", Value: reviewData.Star})
		fallthrough
	case reviewData.Content != oldReview.Content:
		updateTarget = append(updateTarget, bson.E{Key: "content", Value: reviewData.Content})
		fallthrough
	default:
		updateTarget = append(updateTarget, bson.E{Key: "modifyAt", Value: time.Now()})
	}

	updateFilter := bson.M{"orderDay": reviewData.OrderDay, "orderId": reviewData.OrderId, "userId": reviewData.UserId}
	update := bson.D{{Key: "$set", Value: updateTarget}}
	_, err := p.reviewCol.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		log.Error("리뷰 상태 수정 에러", err.Error())
	}

	// 업데이트 후 메뉴에 평점 반영
	err = p.UpdateMenuReviewStarModel(&targetOrder)
	return nil
}

// 리뷰 삭제를 위해 추가했지만, 수정과 기능이 크게 차이가 없음
// model에서 flag를 받아 수정인지 삭제인지 분기를 타는 것도 방법일듯 함(일단 보류)
func (p *Model) DeleteReviewModel(reviewData Review) error {

	var targetOrder Order
	findOrderFilter := bson.M{
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
		"userId":   reviewData.UserId,
		"state":    7,
	}
	if err := p.orderSaveCol.FindOne(context.TODO(), findOrderFilter).Decode(&targetOrder); err != nil {
		log.Error("오더 조회 에러", err.Error())
		return err
	}

	var oldReview Review
	findReviewFilter := bson.M{
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
		"userId":   reviewData.UserId,
	}
	if err := p.reviewCol.FindOne(context.TODO(), findReviewFilter).Decode(&oldReview); err != nil {
		log.Error("리뷰 조회 에러", err.Error())
		return err
	}
	oldReview.Id = nil
	if _, err := p.reviewSaveCol.InsertOne(context.TODO(), oldReview); err != nil {
		log.Error("리뷰 백업 에러", err.Error())
		return err
	}

	deleteFilter := bson.M{
		"orderDay": reviewData.OrderDay,
		"orderId":  reviewData.OrderId,
		"userId":   reviewData.UserId,
	}
	_, err := p.reviewCol.DeleteOne(context.TODO(), deleteFilter)
	if err != nil {
		log.Error("리뷰 상태 삭제 에러", err.Error())
	}

	// 삭제 후 메뉴에 평점 반영
	err = p.UpdateMenuReviewStarModel(&targetOrder)
	return nil
}
