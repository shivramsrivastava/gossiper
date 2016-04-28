package policylib

type DistanceType string

const (
	NEAREST  = "NEAREST"
	FARTHEST = "FARTHEST" //farthest
)

type RuleDistance struct {
	DistanceType string
}

func (this *RuleDistance) GetNearest(data []DCData) []DCData {

	return nil
}

func (this *RuleDistance) GetFarthest(data []DCData) []DCData {
	return nil
}

func (this *RuleDistance) ApplyRule(data []DCData) []DCData {
	switch this.Type {
	case NEAREST:
		return this.GetNearest(data)
	case FARTHEST:
		return this.GetFarthest(data)
	}

	return nil
}
