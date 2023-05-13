package middleware

import (
	"auth-service/auth/dto"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	publicKey, err := ioutil.ReadFile(os.Getenv("AUTH_KEYS_PATH") + "public.pem")
	if err != nil {
		return nil, errors.New("failed to read public key")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, errors.New("failed to encode public key")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// TokenValid checks if the token is valid
func tokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request) (*dto.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		roles, ok := claims["roles"].(string)
		if !ok {
			return nil, err
		}
		tenantsIds, ok := claims["tenant_ids"].(string)
		if !ok {
			return nil, err
		}
		profileID, ok := claims["profile_id"].(string)
		if !ok {
			return nil, err
		}
		return &dto.AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
			Role:       roles,
			TenantsIds: tenantsIds,
			ProfileID:  profileID,
		}, nil
	}
	return nil, err
}

//TokenAuthMiddleware to secure the routes
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tokenValid(c.Request)
		if err != nil {
			response := &dto.APIResponse{
				StatusCode: 4001,
				Message:    err.Error(),
			}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		c.Next()
	}
}
