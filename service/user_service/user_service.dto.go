package user_service

type UserCreateDto struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Re_password string `json:"rePassword"`
}

type UserDeleteDto struct {
	User_id int `json:"userId"`
}

type UserUpdateDto struct {
	User_id  int    `json:"userId"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
