package model

import (
	"context"
	"errors"
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
	menuCol       *mongo.Collection
	reviewCol     *mongo.Collection
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
	Review   Review              `json:"review" bson:"review"`
	CreateAt time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt time.Time           `json:"modifyAt" bson:"modifyAt"`
}
type OrderMenu struct {
	MenuId int    `json:"menuId" bson:"menuId"`
	Name   string `json:"name" bson:"name"`
}

type Review struct {
	Star     float32   `json:"star" bson:"star"`
	Content  string    `json:"content" bson:"content"`
	Use      bool      `json:"use" bson:"use"`
	CreateAt time.Time `json:"createAt" bson:"createAt"`
	ModifyAt time.Time `json:"modifyAt" bson:"modifyAt"`
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
		r.menuCol = db.Collection(cfg.DB.MenuCollection)
		r.reviewCol = db.Collection(cfg.DB.ReviewCollection)
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
	}

	var newOrder Order
	query := bson.M{"_id": res.InsertedID}
	if err = p.orderCol.FindOne(context.TODO(), query).Decode(&newOrder); err != nil {
		log.Error("오더 추가 후 조회 에러", err.Error())
		return nil, err
	}
	return &newOrder, err
}
