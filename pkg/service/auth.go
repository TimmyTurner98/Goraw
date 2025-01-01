package service

import (
	"crypto/sha256"
	"fmt"
)

const salt = "asdmklgfrt"

func generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password + salt))     // Добавляем соль к паролю
	return fmt.Sprintf("%x", hash.Sum(nil)) // Возвращаем хэш в виде строки
}
