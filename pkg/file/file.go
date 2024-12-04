package file

import "time"

type File struct {
	ID string `bson:"_id"`
	Name   string `bson:"filename" json:"filename"`
	Length int64  `bson:"length" json:"length"`
	Date   time.Time `bson:"uploadDate" json:"date"`
}