package policylib

type RuleThreshold struct {
	RecosurceLimit int
}

func (this *RuleThreshold) GetThersholds(data []DCData) []DCData {
	for _, dcData := range data {
		//vall.ALLDCs.List[common]
		for key, value := range dcData.data {

		}

	}
}

func (this *RuleThreshold) ApplyRule(data []DCData) []DCData {

	return this.GetThersholds(data)
}
