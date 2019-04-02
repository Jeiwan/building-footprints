package db

import (
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const footprintsColl = "footprints"
const batchSize = 1000

// Mongo implements MongoDB
type Mongo struct {
	c  *mongo.Client
	db string
}

// NewMongo returns a Mongo
func NewMongo(URL, dbName string) (DB, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(URL))
	if err != nil {
		return nil, err
	}

	if err := c.Connect(context.Background()); err != nil {
		return nil, err
	}

	m := &Mongo{
		c:  c,
		db: dbName,
	}

	if err := m.createIndexes(); err != nil {
		return nil, fmt.Errorf("createIndexes: %s", err.Error())
	}

	return m, nil
}

func (m Mongo) createIndexes() error {
	if _, err := m.c.Database(m.db).Collection(footprintsColl).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		return err
	}

	if _, err := m.c.Database(m.db).Collection(footprintsColl).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{"borough_code", 1}},
			Options: options.Index(),
		},
	); err != nil {
		return err
	}

	return nil
}

// AvgHeightByBoroughCode returns average height by borough code
func (m Mongo) AvgHeightByBoroughCode(boroughCode int) (float64, error) {
	ctx := context.Background()
	pipe := mongo.Pipeline{
		{{"$match", bson.M{"borough_code": boroughCode}}},
		{{"$group", bson.M{
			"_id":        "$borough_code",
			"avg_height": bson.M{"$avg": "$heightroof"},
		}}},
	}

	cur, err := m.c.Database(m.db).Collection(footprintsColl).Aggregate(
		ctx,
		pipe,
	)
	if err != nil {
		return -1.0, err
	}

	var result struct {
		ID        int     `bson:"_id"`
		AvgHeight float64 `bson:"avg_height"`
	}
	if !cur.Next(ctx) {
		return -1.0, errors.Errorf("failed to calculate avg height: no results")
	}

	if err := cur.Decode(&result); err != nil {
		return -1.0, err
	}

	return result.AvgHeight, nil
}

// SaveData saves data to Mongo
func (m Mongo) SaveData(rows [][]interface{}) error {
	ctx := context.Background()

	var wm []mongo.WriteModel
	var row []interface{}

	batches := int(math.Ceil(float64(len(rows)) / float64(batchSize)))
	coll := m.c.Database(m.db).Collection(footprintsColl)

Loop:
	for i := 0; i < batches; i++ {
		for j := 0; j < batchSize; j++ {
			rowIdx := (i * batchSize) + j
			if rowIdx >= len(rows) {
				break Loop
			}
			row = rows[rowIdx]

			m := mongo.NewUpdateOneModel().
				SetFilter(bson.M{"id": row[0]}).
				SetUpdate(bson.M{"$set": bson.M{
					"id":           row[0],
					"borough_code": row[1],
					"heightroof":   row[2],
				}}).
				SetUpsert(true)

			wm = append(wm, m)
		}

		_, err := coll.BulkWrite(ctx, wm)
		if err != nil {
			return err
		}

		wm = nil
	}

	return nil
}
