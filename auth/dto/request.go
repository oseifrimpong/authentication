package dto

// APIResponse structure
type APIResponse struct {
	StatusCode int64                  `json:"status_code,string"`
	Message    string                 `json:"message"`
	Data       map[string]interface{} `json:"data"`
}

type Payload struct {
	AuthRequestData
	// Phone
	// Wechat
}

// Wechat payload
type Wechat struct {
	Code          string `json:"code,omitempty"`
	EncryptedData string `json:"encrypted_data,omitempty"`
	IV            string `json:"iv,omitempty"`
}

// AuthRequestData
type AuthRequestData struct {
	Email     string `json:"email" binding:"required,email"`
	UserName  string `json:"username"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password" binding:"required"`
}

// Phone
type Phone struct {
	PhoneCode int `json:"phone_code,string"`
	Number    int `json:"number,string"`
}

type APIRequest struct {
	Type     string  `json:"type" binding:"required"`
	UserType string  `json:"user_type,omitempty"`
	Payload  Payload `json:"payload" binding:"required,dive"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
	*VerifyEmailRequest
}

type UserServiceAPIResponse struct {
	ID          string   `json:"id,omitempty"`
	Email       string   `json:"email,omitempty"`
	ProfileID   string   `json:"profile_id,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	TenantIds   []string `json:"tenantIds,omitempty"`
	CreatedAt   int      `json:"createdAt,omitempty"`
}

type UserServiceAPIRequest struct {
	User *UserPayload `json:"user"`
}

type UserPayload struct {
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UserName  string `json:"username,omitempty"`
	UserID    string `json:"user_id,omitempty"`
	UserType  string `json:"user_type,omitempty"`
}

type UserServiceAPIResponseLogin struct {
	Data  UserServiceAPIResponse `json:"data"`
	Total int                    `json:"total,omitempty"`
	Page  int                    `json:"page,omitempty"`
	Size  int                    `json:"size,omitempty"`
	Where interface{}            `json:"where,omitempty"`
	Sort  struct {
		CreatedAt string `json:"createdAt"`
	} `json:"sort,omitempty"`
}

type AuthResponse struct {
	Message string                  `json:"message"`
	Data    *UserServiceAPIResponse `json:"data"`
}

type VerifyEmailRequest struct {
	CTX   string `json"ctx" binding:"required"`
	Email string `json:"email" binding:"required" `
}

type RequestPasswordChange struct {
	Email string `json:"email" binding:"required,email"`
}
