package users

type User struct {
	ID       uint   `json:"id,omitempty" gorm:"unique;not null"`
	UserName string `json:"user_name" gorm:"unique;not null"`
	Password string `json:"password"`
}
