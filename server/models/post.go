package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post is mongodb post structure
type Post struct {
	ID         string
	Title      string
	Desc       string
	Body       string
	CreateDate string
	UpdateDate string
}

// CreatePost creates post
func CreatePost(post *Post) (string, error) {
	result, err := collectionPost.InsertOne(context.Background(), post)
	if err != nil {
		return "", err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		return oid.Hex(), nil
	}
	// 이부분 에러처리 어떻게 할지 생각해보자.
	return "", nil
}

// FindAllPost find my entire posts
func FindAllPost() ([]Post, error) {
	cur, err := collectionPost.Find(context.Background(), bson.D{{}})
	if err != nil {
		return nil, err
	}

	var result []Post
	for cur.Next(context.Background()) {
		elem := &bson.D{}

		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		m := elem.Map()
		temp := Post{
			ID:    m["_id"].(primitive.ObjectID).Hex(),
			Title: m["title"].(string),
		}
		result = append(result, temp)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	cur.Close(context.Background())

	return result, nil
}

// FindPostByObjectID find post by _id (objectid)
func FindPostByObjectID(id string) (Post, error) {
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Post{}, err
	}
	filter := bson.M{"_id": docID}
	var result Post

	err = collectionPost.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
