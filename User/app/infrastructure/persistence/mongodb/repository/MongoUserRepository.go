package repository

import (
	"context"
	"strings"
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

func (r *MongoUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.ID < 1 {
		return nil, &domain.ValidationError{Field: "id", Reason: "user ID must be provided"}
	}

	mongoUser := model.FromUser(user)

	_, err := r.Collection.InsertOne(ctx, mongoUser)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, &domain.AlreadyExistsError{Name: user.Email}
		}
		return nil, &domain.InternalError{Msg: "failed to persist user", Err: err}
	}

	return user, nil
}

func (r *MongoUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
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

func (r *MongoUserRepository) Update(ctx context.Context, id int, updates map[string]interface{}) (*domain.User, error) {

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
		case "first_name_private":
			mongoUpdates["first_name_private"] = value
		case "last_name_private":
			mongoUpdates["last_name_private"] = value
		case "ticket_list":
			if domainTickets, ok := value.([]domain.Ticket); ok {
				mongoTickets := make([]model.MongoTicket, len(domainTickets))
				for i, dt := range domainTickets {
					mongoTickets[i] = model.MongoTicket{
						PacketID: dt.PacketID,
						EventID:  dt.EventID,
						Code:     dt.Code,
					}
				}
				mongoUpdates["ticket_list"] = mongoTickets
			} else {
				mongoUpdates["ticket_list"] = value
			}
		}
	}

	if len(mongoUpdates) == 0 {
		return r.GetByID(ctx, id)
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

	return r.GetByID(ctx, id)
}

func (r *MongoUserRepository) Delete(ctx context.Context, id int) (*domain.User, error) {

	user, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	result, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to delete user", Err: err}
	}

	if result.DeletedCount == 0 {
		return nil, &domain.NotFoundError{ID: id}
	}

	return user, nil
}

func (r *MongoUserRepository) getNextID(ctx context.Context) (int, error) {

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

func (r *MongoUserRepository) CreateIndexes(ctx context.Context) error {

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.Collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}

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

func (r *MongoUserRepository) GetUsersByEventID(ctx context.Context, eventID int) ([]*domain.User, error) {
	filter := bson.M{
		"ticket_list": bson.M{
			"$elemMatch": bson.M{
				"event_id": eventID,
			},
		},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to query users by event ID", Err: err}
	}
	defer cursor.Close(ctx)

	var mongoUsers []model.MongoUser
	if err = cursor.All(ctx, &mongoUsers); err != nil {
		return nil, &domain.InternalError{Msg: "failed to decode users", Err: err}
	}

	users := make([]*domain.User, 0, len(mongoUsers))
	for _, mu := range mongoUsers {
		users = append(users, mu.ToDomain())
	}

	return users, nil
}

func (r *MongoUserRepository) GetUsersByPacketID(ctx context.Context, packetID int) ([]*domain.User, error) {
	filter := bson.M{
		"ticket_list": bson.M{
			"$elemMatch": bson.M{
				"packet_id": packetID,
			},
		},
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, &domain.InternalError{Msg: "failed to query users by packet ID", Err: err}
	}
	defer cursor.Close(ctx)

	var mongoUsers []model.MongoUser
	if err = cursor.All(ctx, &mongoUsers); err != nil {
		return nil, &domain.InternalError{Msg: "failed to decode users", Err: err}
	}

	users := make([]*domain.User, 0, len(mongoUsers))
	for _, mu := range mongoUsers {
		users = append(users, mu.ToDomain())
	}

	return users, nil
}
