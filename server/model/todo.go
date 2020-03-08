package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// 追加
type Todo struct {
	ID   primitive.ObjectID `json:"id" bson:"_id"`
	Text string             `json:"text" bson:"text"`
}
