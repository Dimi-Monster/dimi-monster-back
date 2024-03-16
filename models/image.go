package models

import "github.com/kamva/mgm/v3"

type Image struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel `bson:",inline"`
	OwnerID          string   `json:"owner_id" bson:"owner_id"`
	Description      string   `json:"description" bson:"description"`
	Location         string   `json:"location" bson:"location"`
	LikedBy          []string `json:"liked_by" bson:"liked_by"`
}

func NewImage(ownerID string, description string, location string) *Image {
	return &Image{
		OwnerID:     ownerID,
		Description: description,
		Location:    location,
		LikedBy:     []string{},
	}
}

func (i *Image) AddLike(userID string) {
	i.LikedBy = append(i.LikedBy, userID)
}

func (i *Image) DeleteLike(userID string) {
	for j, user := range i.LikedBy {
		if user == userID {
			i.LikedBy = append(i.LikedBy[:j], i.LikedBy[j+1:]...)
			return
		}
	}
}

func (i *Image) HasLike(userID string) bool {
	for _, user := range i.LikedBy {
		if user == userID {
			return true
		}
	}
	return false
}
