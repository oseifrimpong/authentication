package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-redis/redis/v7"
	"golang.org/x/crypto/bcrypt"
)

func GenerateSnowflake() (snowflake.ID, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return 0, err
	}

	id := node.Generate()
	return id, nil
}

func ComparePasswords(hashedPwd string, plainPwd string) (bool, error) {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false, err
	}
	return true, err
}

func SecurePass(pass string) string {
	p := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(p, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func GenerateNonce(redisClient *redis.Client, userID string) (string, error) {
	nonceBytes := make([]byte, 64)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("could not generate nonce")
	}

	TTL := time.Now().Add(time.Minute * 5)
	now := time.Now()

	errAccess := redisClient.Set(base64.URLEncoding.EncodeToString(nonceBytes), userID, TTL.Sub(now)).Err()
	if errAccess != nil {
		return "", errAccess
	}

	return base64.URLEncoding.EncodeToString(nonceBytes), nil
}
