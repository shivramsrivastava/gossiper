package glib

//Create a fake library for gossiper to be used in our test programs

import (
	"errors"
	"time"

	ml "github.com/hashicorp/memberlist"
)

type Fake_memberlist struct {
}

func (F_ml *Fake_memberlist) Create(conf *ml.Config) (*ml.Memberlist, error) {

	if conf.Logger == nil {

		return nil, errors.New("Fake MemberList Create Error")

	}

	return &ml.Memberlist{}, nil
}

func (F_ml *Fake_memberlist) Join(existing []string) (int, error) {

	if len(existing) == 0 {
		/* Join no one */
		return 0, errors.New("Fake Memberlist  Join Error")
	}

	/* Joined every one */
	return len(existing), nil
}

func (F_ml *Fake_memberlist) Leave(tiemout time.Duration) error {

	return nil
}

func (F_ml *Fake_memberlist) Members() *ml.Node {
	return nil
}
