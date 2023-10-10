package middleware

import (
	"github.com/DarioKnezovic/api-gateway/clients"
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
	UNAUTHORIZED         = "Unauthorized"
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
		log.Printf("Checking Authorization for this route: %v", r.URL)
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

		userClient, err := clients.NewUserClient("user-service:50051")
		if err != nil {
			log.Printf("Failed to create UserClient: %v", err)
			utils.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		exists, err := userClient.CheckUserExistence(token)
		// TODO: Handle this case better
		if err != nil {
			log.Printf("Failed to check user existence: %v", err)
			utils.RespondWithError(w, http.StatusUnauthorized, "Token is not valid")
			return
		}

		if exists {
			next.ServeHTTP(w, r)
		} else {
			utils.RespondWithError(w, http.StatusUnauthorized, UNAUTHORIZED)
			return
		}
	}
}
