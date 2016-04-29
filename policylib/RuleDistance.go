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

func (this *RuleDistance) GetNearest(data dcData) dcData {
	log.Println("GetNearest: Applying the Nearest rule")

	return nil
}

func (this *RuleDistance) GetFarthest(data dcData) dcData {
	log.Println("GetNearest: Applying the farthest rule")
	return nil
}

func (this *RuleDistance) ApplyRule(data dcData) dcData {

	log.Println("ApplyRule: Applying the rule for RuleDistance")
	switch this.DistanceType {
	case NEAREST:
		return this.GetNearest(data)
	case FARTHEST:
		return this.GetFarthest(data)
	}

	return nil
}
