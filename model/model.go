package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	conf "WBABEProject-09/config"
	log "WBABEProject-09/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	client    *mongo.Client
	userCol   *mongo.Collection
	orderCol  *mongo.Collection
	menuCol   *mongo.Collection
	reviewCol *mongo.Collection
}

type User struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
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
	User     string              `json:"user" bson:"user"`
	Menu     []OrderMenu         `json:"menu" bson:"menu"`
	Phone    string              `json:"phone" bson:"phone"`
	Address  string              `json:"address" bson:"address"`
	State    int                 `json:"state" bson:"state"`
	Review   Review              `json:"review" bson:"review"`
	CreateAt time.Time           `json:"createAt" bson:"createAt"`
	ModifyAt time.Time           `json:"modifyAt" bson:"modifyAt"`
}
type OrderMenu struct {
	id string `json:"name" bson:"name"`
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
func (p *Model) GetUserTypeByIdModel(userId string) (*User, error) {

	var user User
	objectId, _ := primitive.ObjectIDFromHex(userId)
	filter := bson.M{"_id": objectId}
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

func (p *Model) UpdateMenuModel(menuId string, menuData Menu) error {

	var oldMenu Menu
	objectId, _ := primitive.ObjectIDFromHex(menuId)
	findFilter := bson.M{"_id": objectId}
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

	updateFilter := bson.M{"_id": objectId}
	update := bson.D{{"$set", updateTarget}}
	res, err := p.menuCol.UpdateOne(context.TODO(), updateFilter, update)
	fmt.Println(res)
	if err != nil {
		log.Error("메뉴 수정 에러", err.Error())
	}

	return err
}

func (p *Model) DeleteMenuModel(menuId string) error {

	var oldMenu Menu
	objectId, _ := primitive.ObjectIDFromHex(menuId)
	findFilter := bson.M{"_id": objectId}
	if err := p.menuCol.FindOne(context.TODO(), findFilter).Decode(&oldMenu); err != nil {
		log.Error("메뉴 조회 에러", err)
		return err
	}

	if oldMenu.Use == false {
		log.Info("이미 삭제된 메뉴")
		return errors.New("이미 삭제된 메뉴")
	}

	filter := bson.M{"_id": objectId}
	delete := bson.D{{"$set", bson.D{{"use", false}}}}
	res, err := p.menuCol.UpdateOne(context.TODO(), filter, delete)
	fmt.Println(res)
	if err != nil {
		log.Error("메뉴 삭제 에러", err)
	}

	return err
}
