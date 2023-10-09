package db

import (
	"context"
	"encoding/json"
	"example/bootcamp_ex1/entities"
	"fmt"
	"log/slog"
	"os"
	"reflect"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type redisStorage[T entities.StorageObject] struct {
	client *redis.Client
	prefix string
}

func NewRedisStorage[T entities.StorageObject]() *redisStorage[T] {
	redisStorage := new(redisStorage[T])
	// Creating and assigning client
	redisStorage.client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST"),
	})
	// Assigning prefix to search in redis. it has the form of "entityType:id" "user:b6cfb84-4831-429e-a61b-4d28b154fb8c"
	redisStorage.prefix = reflect.TypeOf(new(T)).String() + ":"

	// Verifying Connection
	_, err := redisStorage.client.Ping(context.Background()).Result()

	if err != nil {
		err = StorageError{
			Code:        RedisConnectionFailed,
			Description: fmt.Sprintf("error connecting with database %q", err),
		}
		panic(err)
	}
	slog.Info("Connection succesful with redis")

	// Returning instance
	return redisStorage
}

func (r *redisStorage[T]) Get(id uuid.UUID) (T, error) {
	return r.getValueCache(id.String())
}

func (r *redisStorage[T]) GetAll() ([]T, error) {
	return r.getAllValuesCache()
}

func (r *redisStorage[T]) Create(thing T) (uuid.UUID, error) {
	id := thing.GetId()
	err := r.setValueCache(id.String(), thing)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *redisStorage[T]) Update(id uuid.UUID, thing T) (T, error) {
	//If thing not exists return error
	_, err := r.Get(id)
	var zeroValue T
	if err != nil {
		return zeroValue, err
	}
	//Updating new record
	err = r.setValueCache(id.String(), thing)
	if err != nil {
		return zeroValue, err
	}
	return thing, nil

}

func (r *redisStorage[T]) Delete(id uuid.UUID) (uuid.UUID, error) {
	//If thing not exists return error
	_, err := r.Get(id)
	if err != nil {
		return uuid.Nil, err
	}
	// Delete thing
	key := r.prefix + id.String()
	err = r.client.Del(context.Background(), key).Err()
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil

}

func (r *redisStorage[T]) setValueCache(key string, thing T) error {
	ctx := context.Background()
	serialized, err := json.Marshal(thing)
	if err != nil {
		return StorageError{
			Code:        ErrMarshalingEntity,
			Description: fmt.Sprintf("error unmarshaling entity %q", err),
		}
	}
	key = r.prefix + key
	err = r.client.Set(ctx, key, string(serialized), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *redisStorage[T]) getValueCache(key string) (T, error) {
	// Try to get the value
	ctx := context.Background()
	key = r.prefix + key
	value, err := r.client.Get(ctx, key).Result()

	// Handle error
	var zeroValue T
	if err != nil {
		return zeroValue, StorageError{
			Code:        ErrEntityNotFound,
			Description: fmt.Sprintf("cannot find entity with this id %q", key),
		}
	}
	// Try to deserialized
	deserialized := new(T)
	err = json.Unmarshal([]byte(value), deserialized)
	//Handle deserialized error
	if err != nil {
		return zeroValue, StorageError{
			Code:        ErrUnmarshalingEntity,
			Description: fmt.Sprintf("cannot umarshal entity: %q", err),
		}
	}

	return *deserialized, nil

}

func (r *redisStorage[T]) getAllValuesCache() ([]T, error) {
	// Getting all the keys of the existing thing records
	ctx := context.Background()
	prefix := r.prefix + "*"
	iter := r.client.Scan(ctx, 0, prefix, 0).Iterator()
	things := make([]T, 0)
	keys := make([]string, 0)

	for iter.Next(ctx) {
		value := iter.Val()
		keys = append(keys, value)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return things, nil
	}
	values, err := r.client.MGet(ctx, keys...).Result()

	if err != nil {
		return nil, StorageError{
			Code:        ErrGettingRecords,
			Description: fmt.Sprintf("error getting all records %q", err),
		}
	}

	for _, val := range values {
		jsonValue := fmt.Sprint(val)
		currentThing := new(T)
		err := json.Unmarshal([]byte(jsonValue), currentThing)
		if err != nil {
			return nil, StorageError{
				Code:        ErrUnmarshalingEntity,
				Description: fmt.Sprintf("error unmarshaling entity %q", err),
			}
		}

		things = append(things, *currentThing)
	}

	return things, nil
}
