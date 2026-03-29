package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const jwtSecret = "todo-scheduler-secret"

type claims struct {
	PassHash string `json:"pass_hash"`
	jwt.RegisteredClaims
}

// signinHandler godoc
// @Summary     Авторизация
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       credentials body object{password=string} true "Пароль"
// @Success     200 {object} map[string]string "JWT-токен"
// @Failure     400 {object} map[string]string "неверный пароль"
// @Router      /api/signin [post]
func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if req.Password != pass {
		writeError(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		PassHash: hashPassword(pass),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	})

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		writeError(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"token": tokenStr}, http.StatusOK)
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "authentification required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неверный метод подписи")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "authentification required", http.StatusUnauthorized)
			return
		}

		c, ok := token.Claims.(*claims)
		if !ok || c.PassHash != hashPassword(pass) {
			http.Error(w, "authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func hashPassword(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}
