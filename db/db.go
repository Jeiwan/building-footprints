package db

// DB describes DB interface
type DB interface {
	SaveData([][]interface{}) error
}
