package service

import (
	"errors"
	"time"

	"template-backend/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret123")

// Login 模拟登录，用户名 admin 密码 123456
func Login(username, password string) (*model.LoginResponse, error) {
	if username != "admin" || password != "123456" {
		return nil, nil
	}
	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString(jwtSecret)

	user := model.UserInfo{
		ID:          1,
		Username:    "admin",
		Email:       "admin@example.com",
		Nickname:    "管理员",
		Roles:       []string{"admin"},
		Permissions: []string{"*:*"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	return &model.LoginResponse{
		Token:     tokenStr,
		UserInfo:  user,
		ExpiresIn: 3600,
	}, nil
}

func GetUserInfo(username string) model.UserInfo {
	return model.UserInfo{
		ID:          1,
		Username:    username,
		Email:       username + "@example.com",
		Nickname:    "管理员",
		Roles:       []string{"admin"},
		Permissions: []string{"*:*", "user:view"},
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
}

// RefreshToken 解析旧 token 并返回新 token（简单策略：只要旧 token 可解析即发新 token）
func RefreshToken(oldToken string) (string, error) {
	tok, err := jwt.Parse(oldToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !tok.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	username, _ := claims["username"].(string)
	newTok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	newTokenStr, _ := newTok.SignedString(jwtSecret)
	return newTokenStr, nil
}

// ChangePassword 模拟函数：只接受旧密码等于 "123456"
func ChangePassword(username, oldPassword, newPassword string) bool {
	if oldPassword != "123456" {
		return false
	}
	// 模拟更新：实际应写 DB / bcrypt 等
	return true
}
