package policyengine

import (
	"log"
	"sync"
	"time"

	"../common"
	"../consulib"
)

type dcData map[string]*common.DC

var (
	dcDataList dcData
)

type PE struct {
	*consulib.ConsulHandle
	policy           *Policy
	Current_DS_Index uint64

	Lck sync.Mutex
}

func NewPE(config *common.ConsulConfig) *PE {
	newPE := &PE{Current_DS_Index: 0}

	var ok bool
	newPE.ConsulHandle, ok = consulib.NewConsulHandle(config)
	if !ok {
		log.Println("NewPe: Error in creating a consulib Handle")
		return nil
	}

	return newPE

}

// ApplyNewPolicy
// Only
//
func (this *PE) ApplyNewPolicy() {

	for {
		//<-common.
		//Apply the this.Policy and take a decision

		<-common.TriggerPolicyCh

		this.Lck.Lock()
		ok := this.policy.TakeDecision()
		if ok != true {
			log.Println("ApplyNewPolicy: TakeDecision on new policy ", this.policy.Name, " failed")
		}
		this.Lck.Unlock()

	}
}

// populate policy
// read and update policy
// if leader repllicate
//
func (this *PE) UpdatePolicyFromDS(config *common.ConsulConfig) {
	for {

		data, resultingIndex, ok := this.WatchStore(this.Current_DS_Index)

		if ok && (this.Current_DS_Index < resultingIndex) {

			this.Lck.Lock()

			if config.IsLeader == true {
				err := this.ReplicateStore(data)
				if err != nil {
					log.Fatalln("UpdatePolicyFromDS: Data replication failed", err)
				}
			}
			//set the new ModifiedIndex
			this.Current_DS_Index = resultingIndex
			for _, value := range data.KVPairs { //only one policy will be passed since we store only one currently in our PE
				log.Println("UpdatePolicyFromDS: key and value ", string(value.Key), string(value.Value))
				newpolicy, err := this.ProcessNewPolicy(value.Key, value.Value)
				if err != nil {
					log.Println("UpdatePolicyFromDS: ProcessNewPolicy failes to ", err)
				}
				this.policy = newpolicy
			}
			this.Lck.Unlock()
		}
		time.Sleep(5 * time.Second)
	}

}

func getTheConsulAndDCpointers() {

	for {

		log.Println("getTheConsulAndDCpointers: Lopping thorugh to get the global maps")

		common.ALLDCs.Lck.Lock()
		dcDataList = common.ALLDCs.List
		common.ALLDCs.Lck.Unlock()

		if dcDataList != nil {
			log.Println("getTheConsulAndDCpointers: Got common.ALLDCs ", dcDataList)
			break
		}

		<-time.After(2 * time.Second)
	}
	log.Println("getTheConsulAndDCpointers: ")
	return
}

// BootStrapPolicy is used to booststap policy
//
func (this *PE) BootStrapPolicy(config *common.ConsulConfig) {

	data, resultingIndex, ok := this.WatchStore(this.Current_DS_Index)

	if ok && (this.Current_DS_Index < resultingIndex) {

		if config.IsLeader == true {
			err := this.ReplicateStore(data)
			if err != nil {
				log.Fatalln("BootStrapPolicy: Data replication failed", err)
			}
		}
		//set the new ModifiedIndex
		this.Current_DS_Index = resultingIndex
		for _, value := range data.KVPairs {
			log.Println("BootStrapPolicy: key and value ", string(value.Key), string(value.Value))
			newpolicy, err := this.ProcessNewPolicy(value.Key, value.Value)
			if err != nil {
				log.Println("BootStrapPolicy: ProcessNewPolicy failes to ", err)
			}
			this.policy = newpolicy
			ok := this.policy.TakeDecision()
			if ok != true {
				log.Println("BootStrapPolicy: TakeDecision on new policy ", this.policy.Name, " failed")
			}
		}

	}

}

//Entry point for the policy engine
func Run(config *common.ConsulConfig) {

	getTheConsulAndDCpointers()

	pe := NewPE(config)

	pe.BootStrapPolicy(config)

	go pe.UpdatePolicyFromDS(config)
	go pe.ApplyNewPolicy()

}
