package data

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisData struct {
	client *redis.Client
}

func NewRedisData(opt redis.Options) DataStore {
	return &RedisData{
		client: redis.NewClient(&opt),
	}
}

func (r *RedisData) Ping() (string, error) {
	pong, err := r.client.Ping(context.Background()).Result()
	if err != redis.Nil {
		return pong, err
	}
	return pong, nil
}

func (r *RedisData) Find(key string) (map[string]string, error) {
	k := "shortner:" + key
	data, err := r.client.HGetAll(context.Background(), k).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("key not found")
	}

	return data, nil

}
func (r *RedisData) Delete(key string) (bool, error) {
	key = "shortner:" + key
	data, err := r.client.Del(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	if data == 0 {
		return false, nil
	}
	return true, nil
}

func (r *RedisData) Store(s map[string]string) error {
	key := "shortner:" + s["short"]
	_, err := r.client.HSet(context.Background(), key, s).Result()
	if err != redis.Nil {
		return err
	}
	return nil
}

func (r *RedisData) SearchUrl(url string) ([]map[string]string, error) {
	var output []map[string]string
	listKey := r.client.Scan(context.Background(), 0, "shortner:*", 1).Iterator()
	for listKey.Next(context.Background()) {
		key := listKey.Val()
		data, err := r.client.HGetAll(context.Background(), key).Result()
		if err != nil {
			continue
		}
		if data["url"] == url {
			output = append(output, data)
		}
	}
	if err := listKey.Err(); err != nil {
		return output, err
	}

	return output, nil
}
