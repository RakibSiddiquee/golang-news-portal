package models

import "github.com/upper/db/v4"

type Models struct {
}

func New(db db.Session) Models {
	return Models{}
}
