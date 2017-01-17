package models

const (
	UserRoleOwner  = 0
	UserRoleAdmin  = 1
	UserRoleEditor = 2
	UserRoleWriter = 3
)

type User struct {
	FullName string `json:"fullName" bson:"fullName"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`
	Salt     string `json:"-" bson:"salt"`
	UID      string `json:"uid" bson:"uid"`
	Role     int8   `json:"role" bson:"role"`
}
