package policylib

import "log"

type RuleThreshold struct {
	RecosurceLimit int
}

func (this *RuleThreshold) GetThersholds(data dcData) dcData {

	log.Println("GetThersholds: Applying the thershold rule")
	for range data {
		//vall.ALLDCs.List[common]

	}
	//we need to return new map
	return nil
}

func (this *RuleThreshold) ApplyRule(data dcData) dcData {

	log.Println("ApplyRule: ApplyRule of the interface to RuleThreshold")

	return this.GetThersholds(data)
}
