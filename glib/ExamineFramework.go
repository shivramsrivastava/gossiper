package glib

import (
	"log"
	"time"

	"../common"
)

func PerformIntersection() {

	log.Printf("PerformIntersection() Sent to TCP cannel")
  FrmWrkLck.Lock()
	for mk, mv := range AllFrameworks {
		for n, _ := range mv {
			isCommon := true
			for sk, sv := range AllFrameworks {
				if mk != sk {
					if _, isvalid := sv[n]; !isvalid {
						isCommon = false
					}
				}
			}
			if isCommon {
				CommonFramework[n] = false
			}
		}
	}

  FrmWrkLck.Unlock()

	common.ToAnon.Lck.Lock()
	for k, v := range CommonFramework {
		common.ToAnon.M[k] = v

	}
	common.ToAnon.Lck.Unlock()
	log.Printf("PerformIntersection() Sent to TCP cannel")
	common.ToAnon.Ch <- true
}

func ExamineFramework() {
	go func() {
		for {
			<-time.After(time.Second * 10)
			log.Printf("Dump AllFramework %v", AllFrameworks)
			log.Printf("Dump CommonFramework %v", CommonFramework)
		}
	}()

	for {
		select {
		case <-time.After(time.Second * 1):
			log.Printf("Performing intersection")
			PerformIntersection()

		case <-time.After(time.Second * 10):
			log.Printf("Dump AllFramework %v", AllFrameworks)
		}
	}

}
