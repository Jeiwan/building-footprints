package db

import (
	"context"
	"fmt"
	"math"

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
	_, err := m.c.Database(m.db).Collection(footprintsColl).Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bson.D{{"id", 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}

// SaveData saves data to Mongo
func (m Mongo) SaveData(rows [][]interface{}) error {
	ctx := context.Background()

	var wm []mongo.WriteModel
	var row []interface{}

	batches := int(math.Ceil(float64(len(rows)) / float64(batchSize)))
	coll := m.c.Database(m.db).Collection(footprintsColl)

	for i := 0; i < batches; i++ {
		for j := 0; j < batchSize; j++ {
			rowIdx := (i * batchSize) + j
			if rowIdx >= len(rows) {
				break
			}
			row = rows[rowIdx]

			m := mongo.NewUpdateOneModel().
				SetFilter(bson.M{"id": row[0]}).
				SetUpdate(bson.M{"$set": bson.M{
					"id":         row[0],
					"bin":        row[1],
					"heightroof": row[2],
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
