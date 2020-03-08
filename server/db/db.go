package db

import (
	"context"
	"log"
	"my-app2/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DB interface {
	GetTechnologies() ([]*model.Technology, error)
	// 追加
	GetTodos() ([]*model.Todo, error)
	PostTodos(string) error
	DeleteTodos(string) error
}

type MongoDB struct {
	collection *mongo.Collection
	// 追加
	todoCl *mongo.Collection
}

func NewMongo(client *mongo.Client) DB {
	tech := client.Database("tech").Collection("tech")
	// 追加
	todo := client.Database("todo").Collection("todo")
	return MongoDB{collection: tech, todoCl: todo}
}

func (m MongoDB) GetTechnologies() ([]*model.Technology, error) {
	res, err := m.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error while fetching technologies:", err.Error())
		return nil, err
	}
	var tech []*model.Technology
	err = res.All(context.TODO(), &tech)
	if err != nil {
		log.Println("Error while decoding technologies:", err.Error())
		return nil, err
	}
	return tech, nil
}

// 追加
func (m MongoDB) GetTodos() ([]*model.Todo, error) {
	res, err := m.todoCl.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("Error while fetching todos:", err.Error())
		return nil, err
	}
	var todo []*model.Todo
	err = res.All(context.TODO(), &todo)
	if err != nil {
		log.Println("Error while decoding todos:", err.Error())
		return nil, err
	}
	return todo, nil
}

// 追加
func (m MongoDB) PostTodos(text string) error {
	mdl := model.Todo{ID: primitive.NewObjectID(), Text: text}
	_, err := m.todoCl.InsertOne(context.TODO(), mdl)
	if err != nil {
		log.Println("Error while fetching todos:", err.Error())
		return err
	}
	return nil
}

// 追加
func (m MongoDB) DeleteTodos(id string) error {
	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err := m.todoCl.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		log.Println("Error while fetching todos:", err.Error())
		return err
	}
	return nil
}
