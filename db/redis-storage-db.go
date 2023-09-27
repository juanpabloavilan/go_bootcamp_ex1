package db

import (
	"context"
	"encoding/json"
	"errors"
	"example/bootcamp_ex1/entities"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrConnectionFailed   = errors.New("couldn't connect to database")
	ErrConsultingRecords  = errors.New("error consulting records")
	ErrUnmarshalingRecord = errors.New("error unmarshaling record")
	ErrMarshalingRecord   = errors.New("error unmarshaling record")
)

type redisStorage struct {
	client *redis.Client
	prefix string
}

func NewRedisStorage() *redisStorage {

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()

	if err != nil {
		slog.Error(ErrConnectionFailed.Error(), err)
		panic(err)
	}
	slog.Info("Connection succesful with redis")

	return &redisStorage{client: client, prefix: "user:"}
}

func (r *redisStorage) Get(id uuid.UUID) (entities.User, error) {
	return r.getValueCache(id.String())
}

func (r *redisStorage) GetAll() ([]entities.User, error) {
	return r.getAllValuesCache()
}

func (r *redisStorage) Create(user entities.User) (uuid.UUID, error) {
	id := user.Id
	err := r.setValueCache(id.String(), user)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *redisStorage) Update(id uuid.UUID, user entities.User) (entities.User, error) {
	//If user not exists return error
	_, err := r.Get(id)
	if err != nil {
		return entities.User{}, err
	}
	//Updating new record
	err = r.setValueCache(id.String(), user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil

}

func (r *redisStorage) Delete(id uuid.UUID) (uuid.UUID, error) {
	//If user not exists return error
	_, err := r.Get(id)
	if err != nil {
		return uuid.Nil, err
	}
	// Delete user
	err = r.deleteValueCache(id.String())
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil

}

func (r *redisStorage) setValueCache(key string, user entities.User) error {
	ctx := context.Background()
	serialized, err := marshalUser(user)
	if err != nil {
		return err
	}
	key = r.prefix + key
	err = r.client.Set(ctx, key, serialized, 0).Err()
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (r *redisStorage) getValueCache(key string) (entities.User, error) {
	ctx := context.Background()
	key = r.prefix + key
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		slog.Error(err.Error())
		return entities.User{}, ErrUserNotFound
	}

	return unmarshalUser(value)

}

func (r *redisStorage) getAllValuesCache() ([]entities.User, error) {
	// Getting all the keys of the existing user records
	ctx := context.Background()
	prefix := r.prefix + "*"
	iter := r.client.Scan(ctx, 0, prefix, 0).Iterator()
	users := make([]entities.User, 0)

	keys := make([]string, 0)

	for iter.Next(ctx) {
		value := iter.Val()
		keys = append(keys, value)
	}
	if err := iter.Err(); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	values, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		slog.Error(err.Error())
		return []entities.User{}, ErrConsultingRecords
	}

	for _, val := range values {
		jsonValue := fmt.Sprint(val)
		user, err := unmarshalUser(jsonValue)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *redisStorage) deleteValueCache(key string) error {
	key = r.prefix + key
	return r.client.Del(context.Background(), key).Err()

}

func marshalUser(user entities.User) ([]byte, error) {
	//Marshaling user
	serialized, err := json.Marshal(user)
	if err != nil {
		slog.Error(err.Error())
		return nil, ErrMarshalingRecord
	}

	return serialized, nil

}

func unmarshalUser(jsonValue string) (entities.User, error) {
	//Unmarshaling json user values
	var deserializedUser entities.User
	err := json.Unmarshal([]byte(jsonValue), &deserializedUser)
	if err != nil {
		slog.Error(err.Error())
		return entities.User{}, ErrUnmarshalingRecord
	}

	return deserializedUser, nil
}
