package db

// DB describes DB interface
type DB interface {
	AvgHeightByBoroughCode(boroughCode int) (float64, error)
	SaveData([][]interface{}) error
}
