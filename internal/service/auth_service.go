// internal/service/auth_service.go
package service

import (
	"errors"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret123")

type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthService {
	return &AuthService{userRepo: userRepo, roleRepo: roleRepo}
}

// Login 用户登录验证
func (s *AuthService) Login(username, password string) (*model.LoginResponse, error) {
	// 根据用户名查找用户
	users, _, err := s.userRepo.GetList(1, 1, map[string]interface{}{"username": username})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("用户不存在")
	}

	user := users[0]

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // token 有效期24小时
	})
	tokenStr, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, err
	}

	// 获取用户角色
	roleNames, roles, err := s.getUserRoles(user.ID)
	if err != nil {
		return nil, err
	}

	// 构建用户信息
	userInfo := model.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Roles:       roleNames,
		Permissions: s.getUserPermissions(roles), // 根据角色获取权限
		CreatedAt:   user.CreatedAt.Format(time.DateTime),
		UpdatedAt:   user.UpdatedAt.Format(time.DateTime),
	}

	return &model.LoginResponse{
		Token:     tokenStr,
		UserInfo:  userInfo,
		ExpiresIn: 24 * 3600, // 24小时
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID uint) (model.UserInfo, error) {
	// 根据ID获取用户
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return model.UserInfo{}, err
	}

	// 获取用户角色
	roleNames, roles, err := s.getUserRoles(user.ID)
	if err != nil {
		return model.UserInfo{}, err
	}

	return model.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Nickname:    user.Nickname,
		Roles:       roleNames,
		Permissions: s.getUserPermissions(roles),
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   user.UpdatedAt.Format(time.RFC3339),
	}, nil
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(oldToken string) (string, error) {
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

	userID, ok := claims["userID"].(float64)
	if !ok {
		return "", errors.New("invalid user ID in token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("invalid username in token")
	}

	// 验证用户是否存在
	_, err = s.userRepo.GetByID(uint(userID))
	if err != nil {
		return "", errors.New("user not found")
	}

	// 生成新token
	newTok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   uint(userID),
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	newTokenStr, err := newTok.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return newTokenStr, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.userRepo.Update(user)
}

// getUserRoles 获取用户角色
func (s *AuthService) getUserRoles(userID uint) ([]string, []model.Role, error) {
	roles, err := s.userRepo.GetUserRoles(userID)
	if err != nil {
		return nil, nil, err
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.RoleName)
	}
	return roleNames, roles, nil
}

// getUserPermissions 根据角色获取权限
func (s *AuthService) getUserPermissions(roles []model.Role) []string {
	permissions := make([]string, 0)

	for _, role := range roles {
		getPermissions, err := s.roleRepo.GetPermissions(role.ID)
		if err != nil {
			continue
		}
		for _, permission := range getPermissions {
			permissions = append(permissions, permission.PermissionCode)
		}
	}

	return permissions
}
