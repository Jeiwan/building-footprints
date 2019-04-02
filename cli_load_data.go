package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jeiwan/building-footprints/db"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type metadataColumn struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Flags    []string `json:"flags"`
}

type metadataView struct {
	Columns []metadataColumn `json:"columns"`
}

type metadata struct {
	View metadataView `json:"view"`
}

type dataFileStructure struct {
	Metadata metadata        `json:"meta"`
	Data     [][]interface{} `json:"data"`
}

func cliLoadData(c *cli.Context) error {
	dataFile, err := os.Open(c.String("data-file"))
	if err != nil {
		return err
	}

	var data dataFileStructure
	if err := json.NewDecoder(dataFile).Decode(&data); err != nil {
		return err
	}

	binColumnIdx := -1
	heightColumnIdx := -1
	idColumnIdx := -1
	skippedColumns := 0
	for _, c := range data.Metadata.View.Columns {
		if len(c.Flags) != 0 && c.Flags[0] == "hidden" {
			skippedColumns++
		}

		if c.Name == "BIN" {
			binColumnIdx = c.Position + skippedColumns - 1
		}

		if c.Name == "HEIGHTROOF" {
			heightColumnIdx = c.Position + skippedColumns - 1
		}

		if c.Name == "DOITT_ID" {
			idColumnIdx = c.Position + skippedColumns - 1
		}
	}
	if binColumnIdx == -1 {
		return fmt.Errorf("BIN column not found")
	}

	if heightColumnIdx == -1 {
		return fmt.Errorf("HEIGHTROOF column not found")
	}

	if idColumnIdx == -1 {
		return fmt.Errorf("DOITT_ID column not found")
	}

	for i, building := range data.Data {
		bin, ok := building[binColumnIdx].(string)
		if !ok || len(bin) != 7 {
			logrus.Errorf("wrong BIG value: %v", building[binColumnIdx])
			continue
		}

		boroughCode := string(bin[0])

		trimmedData := []interface{}{
			building[idColumnIdx],
			boroughCode,
			building[heightColumnIdx],
		}
		data.Data[i] = trimmedData
	}

	db, err := db.NewMongo(c.String("mongo-url"), c.String("mongo-db-name"))
	if err != nil {
		return err
	}

	return db.SaveData(data.Data)
}
