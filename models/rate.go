// article_rate
package models

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type ArticleRate struct {
	Article string `json:"article"`
	Rate    int    `json:"rate"`
}

type UserRate struct {
	Id     bson.ObjectId `bson:"_id,omitempty"`
	Userid string
	Rates  []ArticleRate
}

func IterRate(f func(*UserRate)) error {
	rate := UserRate{}
	q := func(c *mgo.Collection) error {
		iter := c.Find(nil).Iter()
		for iter.Next(&rate) {
			f(&rate)
		}
		return iter.Close()
	}

	return withCollection(rateColl, nil, q)
}
