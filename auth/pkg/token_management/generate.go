package token_management

import (
	"auth-service/auth/dto"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func CreateToken(userID int64, roles string, tenantIds string, profileID string, isVerified bool) (*dto.TokenDetails, error) {
	privateKey, _ := ioutil.ReadFile(os.Getenv("AUTH_KEYS_PATH") + "private.pem")

	td := &dto.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 480).Unix() //Expires after 8 hours
	td.AccessUUID = uuid.New().String()

	/*
		Reason for setting the refresh token to 2 days
		https://cloud.google.com/apigee/docs/api-platform/antipatterns/oauth-long-expiration
	*/
	td.RtExpires = time.Now().Add(time.Hour * 48).Unix() // Expires after 2 days
	td.RefreshUUID = uuid.New().String()

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = strconv.FormatInt(userID, 10)
	atClaims["roles"] = roles
	atClaims["tenant_ids"] = tenantIds
	atClaims["profile_id"] = profileID
	atClaims["is_verified"] = isVerified
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, atClaims)
	pk, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, err
	}
	td.AccessToken, err = at.SignedString(pk)
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = strconv.FormatInt(userID, 10)
	rtClaims["roles"] = roles
	rtClaims["tenant_ids"] = tenantIds
	rtClaims["profile_id"] = profileID
	rtClaims["is_verified"] = isVerified
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodRS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(pk)
	if err != nil {
		return nil, err
	}
	return td, err
}

func RefreshToken(redis *redis.Client, refreshToken string) (res map[string]interface{}, err error) {
	publicKey, _ := ioutil.ReadFile(os.Getenv("AUTH_KEYS_PATH") + "public.pem")

	key, _ := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		RefreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, err
		}

		deleted, delErr := DeleteKV(redis, RefreshUUID)
		if delErr != nil || deleted == 0 {
			return nil, err
		}

		userID, _ := strconv.ParseInt(claims["user_id"].(string), 10, 64)
		roles := claims["roles"].(string)
		tenantIds := claims["tenant_ids"].(string)
		profileID := claims["profile_id"].(string)
		isVerified := claims["is_verified"].(bool)

		ts, createErr := CreateToken(userID, roles, tenantIds, profileID, isVerified)
		if createErr != nil {
			return nil, err
		}

		saveErr := SaveTokenDetails(redis, userID, ts)
		if saveErr != nil {
			return nil, err
		}

		t := make(map[string]interface{})
		t["access_token"] = ts.AccessToken
		t["refresh_token"] = ts.RefreshToken
		t["expires_at"] = ts.AtExpires

		return t, nil
	}
	return nil, err
}
