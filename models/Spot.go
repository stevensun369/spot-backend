package models

type Spot struct {
	ID          string      `bson:"id" json:"id"`
	Name        string      `bson:"name" json:"name"`
	Description string      `bson:"description" json:"description"`
	Location    GeoLocation `bson:"location" json:"location"`
	Tags        []string    `bson:"tags" json:"tags"`
	Address     string      `bson:"address" json:"address"`
	Phone       string      `bson:"phone" json:"phone"`
	Website     string      `bson:"website" json:"website"`
	Type        string      `bson:"type" json:"type"`
}

type GeoLocation struct {
	Latitude  float64
	Longitude float64
}
