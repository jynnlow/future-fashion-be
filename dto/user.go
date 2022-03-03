package dto

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	DOB      string `json:"dob"`
	Role     string `json:"role"`
}
