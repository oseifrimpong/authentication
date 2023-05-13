package token_management

import (
	"auth-service/auth/dto"
	"time"

	"github.com/go-redis/redis/v7"
)

func RetrieveKV(rd *redis.Client, key string) (string, error) {
	key, err := rd.Get(key).Result()
	if err != nil {
		return "Token Doesn't exist", err
	}
	return key, nil
}

func SaveTokenDetails(rd *redis.Client, userID int64, td *dto.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := rd.Set(td.AccessUUID, userID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := rd.Set(td.RefreshUUID, userID, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func DeleteKV(rd *redis.Client, key string) (int64, error) {
	deleted, err := rd.Del(key).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
