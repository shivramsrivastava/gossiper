package glib

import (
	//"log"
	"time"
)

func PerformIntersection() {

}

func ExamineFramework() {
	for {
		select {
		case <-time.After(time.Second * 1):
			PerformIntersection()
		}
	}
}
