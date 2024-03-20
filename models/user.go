package models

import "github.com/kamva/mgm/v3"

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel `bson:",inline"`
	GID              string   `json:"gid" bson:"gid"`
	Name             string   `json:"name" bson:"name"`
	Email            string   `json:"email" bson:"email"`
	RefreshTokens    []string `json:"refresh_tokens" bson:"refresh_tokens"`
	Banned           bool     `json:"banned" bson:"banned"`
}

func NewUser(gid string, name string, email string) *User {
	return &User{
		GID:           gid,
		Name:          name,
		Email:         email,
		RefreshTokens: []string{},
		Banned:        false,
	}
}

func (u *User) HasRefreshToken(token string) bool {
	for _, t := range u.RefreshTokens {
		if t == token {
			return true
		}
	}
	return false
}

func (u *User) RemoveRefreshToken(token string) {
	for i, t := range u.RefreshTokens {
		if t == token {
			u.RefreshTokens = append(u.RefreshTokens[:i], u.RefreshTokens[i+1:]...)
			return
		}
	}
}
