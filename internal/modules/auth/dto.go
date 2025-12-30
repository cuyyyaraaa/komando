package auth

type LoginRequest struct {
	NIP      string `json:"nip" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDTO struct {
	UserID     int64  `json:"user_id"`
	NIP        string `json:"nip"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	RoleCode   string `json:"role_code"`
	RegionalID int64  `json:"regional_id"`
	IsActive   bool   `json:"is_active"`
}

type LoginResponse struct {
	AccessToken string  `json:"access_token"`
	User        UserDTO `json:"user"`
}
