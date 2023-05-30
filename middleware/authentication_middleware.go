package middleware

import (
	"fmt"
	"github.com/DarioKnezovic/api-gateway/config"
	"github.com/DarioKnezovic/api-gateway/utils"
	"log"
	"net/http"
	"strings"
)

const (
	AUTHORIZATION_HEADER = "Authorization"
	BEARER               = "bearer"
	AUTH_TOKEN_MISSING   = "Authorization token is missing"
	AUTH_TOKEN_NOT_VALID = "Authorization token is not valid"
)

func extractTokenFromBearerHeader(authorizationHeader string) (string, bool) {
	parts := strings.Split(authorizationHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != BEARER {
		return "", false
	}

	return parts[1], true
}

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := config.LoadConfig()

		authorizationHeader := r.Header.Get(AUTHORIZATION_HEADER)
		if authorizationHeader == "" {
			log.Println("Auth token is missing")
			utils.RespondWithError(w, http.StatusForbidden, AUTH_TOKEN_MISSING)
			return
		}

		token, valid := extractTokenFromBearerHeader(authorizationHeader)
		if !valid {
			log.Println("Auth token is not valid")
			utils.RespondWithError(w, http.StatusForbidden, AUTH_TOKEN_NOT_VALID)
			return
		}

		claims, err := utils.VerifyJWT(token, []byte(cfg.JWTSecretKey))
		if err != nil {
			log.Println(err)
			utils.RespondWithError(w, http.StatusForbidden, AUTH_TOKEN_NOT_VALID)
			return
		}
		fmt.Println(claims)
		// TODO: Implement gRPC

		// Call the next handler
		next.ServeHTTP(w, r)
	}
}
