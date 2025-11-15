package repository

import (
	"context"
	"strings"
	"time"
	"userService/application/domain"
	"userService/infrastructure/persistence/mongodb/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		Collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Create(user *domain.User) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the next ID
	nextID, err := r.getNextID(ctx)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to generate ID", Err: err}
	}

	user.ID = nextID
	mongoUser := model.FromUser(user)

	_, err = r.Collection.InsertOne(ctx, mongoUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, &domain.AlreadyExistsError{Name: user.Email}
		}
		return nil, &domain.InternalError{Msg: "failed to persist user", Err: err}
	}

	return user, nil
}

func (r *MongoUserRepository) GetByID(id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var mongoUser model.MongoUser
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&mongoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &domain.NotFoundError{ID: id}
		}
		return nil, &domain.InternalError{Msg: "failed to retrieve user", Err: err}
	}

	return mongoUser.ToDomain(), nil
}

func (r *MongoUserRepository) Update(id int, updates map[string]interface{}) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert field names to MongoDB field names
	mongoUpdates := bson.M{}
	for key, value := range updates {
		switch key {
		case "email":
			mongoUpdates["email"] = value
		case "first_name":
			mongoUpdates["first_name"] = value
		case "last_name":
			mongoUpdates["last_name"] = value
		case "social_media_links":
			mongoUpdates["social_media_links"] = value
		case "ticket_list":
			mongoUpdates["ticket_list"] = value
		}
	}

	if len(mongoUpdates) == 0 {
		return r.GetByID(id)
	}

	result := r.Collection.FindOneAndUpdate(
		ctx,
		bson.M{"id": id},
		bson.M{"$set": mongoUpdates},
	)

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, &domain.NotFoundError{ID: id}
		}
		if mongo.IsDuplicateKeyError(result.Err()) {
			return nil, &domain.AlreadyExistsError{Name: "email"}
		}
		return nil, &domain.InternalError{Msg: "failed to update user", Err: result.Err()}
	}

	return r.GetByID(id)
}

func (r *MongoUserRepository) Delete(id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First, get the user to return it
	user, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Then delete it
	result, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to delete user", Err: err}
	}

	if result.DeletedCount == 0 {
		return nil, &domain.NotFoundError{ID: id}
	}

	return user, nil
}

// getNextID generates the next sequential ID for users
func (r *MongoUserRepository) getNextID(ctx context.Context) (int, error) {
	// Find the document with the highest ID
	opts := options.Find().SetSort(bson.D{{Key: "id", Value: -1}}).SetLimit(1)
	cursor, err := r.Collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var users []model.MongoUser
	if err = cursor.All(ctx, &users); err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 1, nil
	}

	return users[0].ID + 1, nil
}

// CreateIndexes creates necessary indexes for the users collection
func (r *MongoUserRepository) CreateIndexes(ctx context.Context) error {
	// Create unique index on email
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	// Create unique index on id
	indexModel = mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = r.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

	return nil
}
