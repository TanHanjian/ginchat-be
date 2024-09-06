package user_service

type UserCreateDto struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Re_password string `json:"rePassword"`
}
