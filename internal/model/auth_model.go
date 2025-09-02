package model

type LoginForm struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember *bool  `json:"remember,omitempty"`
}

type UserInfo struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar,omitempty"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

type LoginResponse struct {
	Token     string   `json:"token"`
	UserInfo  UserInfo `json:"userInfo"`
	ExpiresIn int      `json:"expiresIn,omitempty"`
}

type ChangePasswordForm struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}
