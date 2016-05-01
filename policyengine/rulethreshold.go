package policyengine

import (
	"log"

	"../common"
)

type RuleThreshold struct {
	RecosurceLimit int
}

func (this *RuleThreshold) ApplyRule(policydecision *PolicyDecision) bool {

	log.Println("ApplyRule: ApplyRule of the interface to RuleThreshold")
	//set the DC thershold
	//no chnage the dataset
	common.ResourceThresold = this.RecosurceLimit
	log.Println("ApplyRule: RuleThreshold setting the Threshold")
	return true
}
