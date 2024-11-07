package protocols

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserStruct struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Status      string             `json:"status" bson:"status"`
	Password    string             `json:"password" bson:"password"`
	Language    string             `json:"language" bson:"language"`
	Permissions UserPermissions    `json:"permissions" bson:"permissions"`
}

type AuthUserStruct struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Status      string             `json:"status" bson:"status"`
	Language    string             `json:"language" bson:"language"`
	Permissions UserPermissions    `json:"permissions" bson:"permissions"`
	Token       string             `json:"token"`
}

type UserPermissions struct {
	CreateRollout bool `json:"create_rollout" bson:"create_rollout"`
	UpdateRelease bool `json:"update_release" bson:"update_release"`
	RemoveRollout bool `json:"remove_rollout" bson:"remove_rollout"`
	Admin         bool `json:"admin" bson:"admin"`
}

type AuthUser struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Mongo interface {
	FindAllUsers() ([]UserStruct, error)
	FindUserByEmail(email string) (UserStruct, error)
	FindUserById(id string) (UserStruct, error)
	CreateUser(UserStruct UserStruct) error
	UpdateUser(id string, UserStruct UserStruct) (UserPermissions, error)
	DeleteUser(id string) error
}
