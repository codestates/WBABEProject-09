package model

import (
	"context"

	conf "WBABEProject-09/config"

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
