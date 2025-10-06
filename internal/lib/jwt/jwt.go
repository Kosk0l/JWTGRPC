package jwt

import (
	"JWTGRPC/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//===================================================================================================================//

// ГЕНЕРАЦИЯ ТОКЕНОВ;
// Приходит 2 объекта классов и duration, из них берутся данные и заполняется токен
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256) // Генератор
	claims := token.Claims.(jwt.MapClaims) // Обертка - Мапа
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	// Заполнение токена
	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}