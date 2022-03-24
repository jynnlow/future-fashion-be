package dto

type UserRequest struct {
	Username string  `json:"username"`
	Password string  `json:"password"`
	DOB      string  `json:"dob"`
	Role     string  `json:"role"`
	Chest    float32 `json:"chest"`
	Waist    float32 `json:"waist"`
	Hip      float32 `json:"hip"`
}

type ListUsersResponse struct {
	Users []*UserResponse `json:"users"`
}

type EditUserReq struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	DOB      string  `json:"dob"`
	Role     string  `json:"role"`
	Chest    float32 `json:"chest"`
	Waist    float32 `json:"waist"`
	Hip      float32 `json:"hip"`
}

type UserResponse struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	DOB      string  `json:"dob"`
	Role     string  `json:"role"`
	Chest    float32 `json:"chest"`
	Waist    float32 `json:"waist"`
	Hip      float32 `json:"hip"`
}
