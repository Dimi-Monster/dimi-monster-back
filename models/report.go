package models

import (
	"github.com/kamva/mgm/v3"
	"pentag.kr/dimimonster/utils/random"
)

type Report struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model.
	mgm.DefaultModel `bson:",inline"`
	TargetImageID    string `json:"target_image_id" bson:"target_image_id"`
	TargetOwnerID    string `json:"target_owner_id" bson:"target_owner_id"`
	ReporterID       string `json:"reporter_id" bson:"reporter_id"`
	Category 	     string `json:"category" bson:"category"`
	Reason           string `json:"reason" bson:"reason"`
	ProcessReason    string `json:"process_reason" bson:"process_reason"`
	Status 		     string `json:"status" bson:"status"`
	Secret           string `json:"secret" bson:"secret"`
}

func NewReport(targetImageID string, reporterID string, category string, reason string) (*Report, error) {
	foundImage := &Image{}
	err := mgm.Coll(foundImage).FindByID(targetImageID, foundImage)
	if err != nil {
		return nil, err
	}
	return &Report{
		TargetImageID: targetImageID,
		TargetOwnerID: foundImage.OwnerID,
		ReporterID:    reporterID,
		Category:      category,
		Reason:        reason,
		Status:  	   "pending",
		Secret:        random.RandString(32),
	}, nil
}
