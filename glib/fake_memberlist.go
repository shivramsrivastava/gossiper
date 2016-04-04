package glib

import (
	"log"

	ml "github.com/hashicorp/memberlist"
)

/*

type ML_interface interface {
	Create(conf *ml.Config) (*ml.Memberlist, error)
	Join(existing []string) (int, error)
	Leave(timeout time.Duration) error
	Members() []*ml.Node
}
*/

type Fake_memberlist struct {
}

func (F_ml *Fake_memberlist) Create(conf ml.Config) (ml.Memberlist, error) {

	if conf.Logger == nil {

	}

}
