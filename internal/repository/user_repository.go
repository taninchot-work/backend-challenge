package repository

import (
	"context"
	"fmt"
	"github.com/taninchot-work/backend-challenge/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (entity.User, error)
	GetUserList(ctx context.Context, offset int, limit int) ([]entity.User, error)
	SaveUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userRepositoryImpl struct {
	mongoCollection *mongo.Collection
}

func NewUserRepository(mongoCollection *mongo.Collection) UserRepository {
	return &userRepositoryImpl{
		mongoCollection: mongoCollection,
	}
}

func (r *userRepositoryImpl) GetUserById(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting ID string '%s' to ObjectID: %v", id, err)
		return entity.User{}, fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{"_id": objectID}

	err = r.mongoCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println("Error finding user by ID:", err)
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) GetUserList(ctx context.Context, offset int, limit int) ([]entity.User, error) {
	var users []entity.User
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSkip(int64(offset))
	findOptions.SetLimit(int64(limit))

	cursor, err := r.mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Println("Error finding users:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		log.Println("Error decoding users:", err)
		return nil, err
	}
	return users, nil
}

func (r *userRepositoryImpl) SaveUser(ctx context.Context, user entity.User) (entity.User, error) {
	_, err := r.mongoCollection.InsertOne(ctx, user)
	if err != nil {
		log.Println("Error creating user:", err)
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	filter := bson.M{"email": email}

	err := r.mongoCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println("Error finding user by email:", err)
		return entity.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"updated_at": time.Now(),
	}}

	_, err := r.mongoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("Error updating user:", err)
		return entity.User{}, err
	}

	return user, nil
}

func (r *userRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting ID string '%s' to ObjectID: %v", id, err)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	filter := bson.M{"_id": objectID}

	_, err = r.mongoCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Println("Error deleting user:", err)
		return err
	}
	return nil
}
