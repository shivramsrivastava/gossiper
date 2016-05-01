package policyengine

// This interface defines the behaviour of a rule
import (
	"encoding/json"
	"log"
	"strings"
	"time"

)


// All rule typr must implement this rule
type RuleInterface interface {
	ApplyRule(*PolicyDecision) bool
}

// A Group of rules which forms a policy
//
type Policy struct {
	Name  string //name of the policy which will be the key in our case
	Rules []Rule
}

//Rule are applied based on the rule priority
//it represents one single rule
// Note: If all the rules have the same priority picks the one which our sort function returns.
type Rule struct {
	Name string //this should hold the name of the datatype/struct which will be used to
	//instiantiate the type

	Priority int         // it should be between 1-10 this 10 being highest
	Scope    string      // global scope/local scope
	Content  interface{} // this represt the actual data for the rule, this can be of any type
}

type FakeJsonRuleContent struct {
	Content interface{} // this represt the actual data for the rule, this can be of any type
}

type PolicyDecision struct {
	//after applying all the policy
	SortedDCName     []string  //will have all the sorted dc index 0 is the one which need to be get offers rest suppressed
	LastDecisionTime time.Time // last time this policy decision was taken
	SortValue        []float64
	// can be later use to refresh the same
}

func NewPolicyDecision() *PolicyDecision {
	return &PolicyDecision{}
}

func (this *PE) ProcessNewPolicy(key string, data []byte) (*Policy, error) {
	tempPolicy := &Policy{}
	err := json.Unmarshal(data, &tempPolicy)

	if err != nil {
		log.Println("processNewPolicy: Unable to unmarshal the processNewPolicy", err)
		return nil, err
	}

	for index, values := range tempPolicy.Rules {

		dummy := FakeJsonRuleContent{}
		ruleType := GetCorrectRuleType(values.Name)
		if ruleType != nil {
			dummy.Content = ruleType
			err := json.Unmarshal(data, &dummy)
			if err != nil {
				log.Println("processNewPolicy: Unable to get the content from the json")
				return nil, err
			}
			tempPolicy.Rules[index].Content = dummy.Content
			log.Println("processNewPolicy: info", tempPolicy)

		}
	}

	return tempPolicy, nil
}

//RuleThreshold
func GetCorrectRuleType(name string) interface{} {

	if strings.Contains(name, "MinMax") == true {
		return &RuleMinMax{}
	} else if strings.Contains(name, "Threshold") == true {
		return &RuleThreshold{}
	} else {
		log.Println("GetCorrectRuleType: No Rule Type found for", name)
	}
	return nil

}

func (this *Policy) TakeDecision() bool {

	//create a new policy decision
	var ok bool
	var ruleInterface RuleInterface
	newDecision := NewPolicyDecision()

	newDecision.SortedDCName, ok = GetValidDCsInfo()
	if ok != true {
		log.Println("TakeDecision: Unable to take a decision sort on DC's failed")
	}
	//apply priority
	/*ok = this.ApplyPriority()
	if ok != true {
		log.Println("TakeDecision: Apply priority failes on Policy")
	}*/
	for _, rule := range this.Rules {
		ruleInterface, ok = rule.Content.(RuleInterface)
		if ok != true {
			log.Println("TakeDecision: unable to get the Rule failed to convert to RuleInterface")
			return false
		}
		//this.ApplyPolicyDecisionOnRules(rule.Content)
		ok := ruleInterface.ApplyRule(newDecision)
		if ok != true {
			log.Println("TakeDecision: Applying rule failes", rule.Name, newDecision.SortedDCName)
		}

	}
	//newDecision.SortedDCName[0]//will be the curent one to get offers

	return true
}

//GetDCDataInSortedOrderByDcName
// this will sort the dc names in
// NewPolicyDecision
func GetValidDCsInfo() ([]string, bool) {

	if len(dcDataList) == 0 {
		log.Println("GetValidDCsInfo: the common DC maps is empty")
		return nil, false
	}

	dcDataSortedList := []string{}
	for key, val := range dcDataList {
		if val.OutOfResource == false {
			dcDataSortedList = append(dcDataSortedList, key)
		}
	}
	return dcDataSortedList, true

}
