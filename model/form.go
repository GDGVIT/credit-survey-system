package model

import (
	"time"
)

type Form struct {
	FormId          string         `bson:"_id"`
	CreditsAllotted float64        `bson:"creditsAllotted"`
	UserId          string         `bson:"userId"`
	NumResp         int64          `bson:"numResp"`
	ResponseRate    float64        `bson:"responseRate"`
	GainRate        float64        `bson:"gainRate"`
	Title           string         `bson:"title"`
	PublicDash      bool           `bson:"publicDash"`
	ShowEmail       bool           `bson:"showEmail"`
	IsAnonymous     bool           `bson:"isAnonymous"`
	Expiry          time.Time      `bson:"expiry"`
	Tags            []FormTag      `bson:"tags"`
	IsPublished     bool           `bson:"isPublished"`
	Audience        []FormAudience `bson:"audience"`
	Description     string         `bson:"description"`
}
