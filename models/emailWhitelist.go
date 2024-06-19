package models

import "github.com/kamva/mgm/v3"

type EmailWhitelist struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel `bson:",inline"`
	Email 		  string `json:"email" bson:"email"`
}

func NewEmailWhite(gid string,  email string) *EmailWhitelist {
	return &EmailWhitelist{
		Email:         email,
	}
}
