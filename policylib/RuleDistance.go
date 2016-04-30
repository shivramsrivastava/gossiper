package policylib

import "log"

type DistanceType string

const (
	NEAREST  = "NEAREST"
	FARTHEST = "FARTHEST" //farthest
)

type RuleDistance struct {
	DistanceType string
}

func (this *RuleDistance) GetNearest(policydecision *PolicyDecision) bool {
	log.Println("GetNearest: Applying the Nearest rule")

	return true
}

func (this *RuleDistance) GetFarthest(policydecision *PolicyDecision) bool {
	log.Println("GetNearest: Applying the farthest rule")
	return true
}

func (this *RuleDistance) ApplyRule(policydecision *PolicyDecision) bool {

	log.Println("ApplyRule: Applying the rule for RuleDistance")
	switch this.DistanceType {
	case NEAREST:
		return this.GetNearest(policydecision)
	case FARTHEST:
		return this.GetFarthest(policydecision)
	}

	return true
}
