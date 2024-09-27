package user_service

type UserCreateDto struct {
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required"`
	Re_password string `json:"rePassword" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,phone"`
}

type UserDeleteDto struct {
	User_id int `json:"userId" validate:"required"`
}

type UserUpdateDto struct {
	Name     string `json:"name" validate:"omitempty"`
	Password string `json:"password" validate:"omitempty"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,phone"`
}

type LoginByUserPhoneDto struct {
	Phone    string `json:"phone" validate:"required,phone"`
	Password string `json:"password" validate:"required"`
}

type LoginByUserEmailDto struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
