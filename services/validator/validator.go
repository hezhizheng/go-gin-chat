package validator

type User struct {
	Username  string `binding:"required,max=16,min=2"`
	Password  string `binding:"required,max=32,min=6"`
	AvatarId  string `binding:"required,numeric"`
}
