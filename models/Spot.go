package models

import (
	"backend/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Spot struct {
	ID          string       `bson:"id" json:"id"`
	Name        string       `bson:"name" json:"name"`
	Description string       `bson:"description" json:"description"`
	Location    GeoJSONPoint `bson:"location" json:"location"`
	Tags        []string     `bson:"tags" json:"tags"`
	Address     string       `bson:"address" json:"address"`
	Phone       string       `bson:"phone" json:"phone"`
	Website     string       `bson:"website" json:"website"`
	Type        string       `bson:"type" json:"type"`
}

type GeoJSONPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"`
}

func (spot *Spot) Create() error {
	spot.Tags = []string{}

	_, err := db.Spots.InsertOne(db.Ctx, spot)

	return err
}

func (spot *Spot) UpdateSpot(id string, updates any) error {
	_, err := db.Spots.UpdateOne(
		db.Ctx,
		bson.M{"id": id},
		bson.M{
			"$set": updates,
		},
	)

	return err
}

func GetSpots() (spots []Spot, err error) {

	cursor, err := db.Spots.Find(db.Ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(db.Ctx, &spots); err != nil {
		return nil, err
	}

	return spots, nil
}

func GetSpotsFilter(center []float64, radius int, tags []string) (spots []Spot, err error) {
	earthRadiusMeters := 6378127.0
	radiusRadians := float64(radius) / earthRadiusMeters

	filter := bson.M{}

	// if no tags are passed,
	// it will not include them in the query
	if len(tags) == 0 {
		filter = bson.M{
			"location": bson.M{
				"$geoWithin": bson.M{
					"$centerSphere": bson.A{
						center, radiusRadians,
					},
				},
			},
		}
	} else {
		filter = bson.M{
			"location": bson.M{
				"$geoWithin": bson.M{
					"$centerSphere": bson.A{
						center, radiusRadians,
					},
				},
			},
			"tags": bson.M{
				"$in": tags,
			},
		}
	}

	cursor, err := db.Spots.Find(db.Ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(db.Ctx)

	if err = cursor.All(db.Ctx, &spots); err != nil {
		return nil, err
	}
	return spots, nil
}

func SearchSpots(input string) (spots []Spot, err error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": input,
		},
	}

	opts := options.Find().
		SetProjection(
			bson.M{
				"score": bson.M{
					"$meta": "textScore",
				},
			},
		).
		SetSort(
			bson.M{
				"score": bson.M{
					"$meta": "textScore",
				},
			},
		)

	cursor, err := db.Spots.Find(db.Ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(db.Ctx)

	if err := cursor.All(db.Ctx, &spots); err != nil {
		return nil, err
	}
	return spots, nil
}
