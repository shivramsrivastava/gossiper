package policylib

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"../common"
	"../consulLib"
)

type dcData map[string]*common.DC

var (
	dcDataList   dcData
	dcConsulList *consulLib.FederaConsulClient
)

// This interface defines the behaviour of a rule
// All rule typr must implement this rule
type RuleInterface interface {
	ApplyRule(*PolicyDecision) bool
}

type PolicyConsul struct {
	*consulLib.FederaConsulClient
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

func (this *Policy) ApplyPolicyDecisionOnRules(runleinterface RuleInterface) {

}

func (this *Policy) TakeDecision() bool {

	//create a new policy decision
	var ok bool
	var ruleInterface RuleInterface
	newDecision := NewPolicyDecision()

	newDecision.SortedDCName, ok = GetDCDataInSortedOrderByDcName()
	if ok != true {
		log.Println("TakeDecision: Unable to take a decision sort on DC's failed")
	}
	//apply priority
	ok = this.ApplyPriority()
	if ok != true {
		log.Println("TakeDecision: Apply priority failes on Policy")
	}
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
//
func GetDCDataInSortedOrderByDcName() ([]string, bool) {

	if len(dcDataList) == 0 {
		log.Println("GetDCDataInSortedOrderByDcName: the common DC maps is empty")
		return nil, false
	}

	dcDataSortedList := make([]string, len(dcDataList))

	i := 0
	for key := range dcDataList {
		dcDataSortedList[i] = key
		i++
	}

	sort.Strings(dcDataSortedList)

	return dcDataSortedList, true

}

// Policy runs only when there is change in the DC info from the gosspier
// or when there is a change in the Consul KV "Fedra" store
func ListenAllDcDataChn() {
	//<-

}

// ListenDcConsulList will block till the KV client connection is established
// TODO: We will read the Consul DC data and update the global pointer dcConsulList
// This has to be controlled via a mutex.
// But for now we dont lock while reading.
func ListenDcConsulList() {
	//<-
}

// Its tricky to actually make a block call on the consul looking for a KV prefix
// If there is not such KV store we create one with the /Fedra with some data/no data
// we then use the modified index from that store for later block
// config is passed to get the consul KV config directly
func ListenConsulKV(Prefix string, wg *sync.WaitGroup, config *common.ConsulConfig) {

	var waitIndex uint64

	// This is a blocking call since we need to watch the store for any new policy change
	// and notify processNewPloicy

	for {

		//GlobalConsulDCInfo should have valid consul client connections

		if dcConsulList == nil {
			log.Println("dcConsulList global consullist is empty")
			os.Exit(0)
		}

		KVPairs, nextWaitIndex, err := dcConsulList.GetList(Prefix, waitIndex)
		//dcConsulList.GetDataFromLocalKVStore(waitIndex)

		if err != nil {
			log.Println("ListenConsulKV: Watch on local store failes", err)
			//TODO: We just continue for now

		} else {
			log.Println("ListenConsulKV: Got the following data from consul current and new wait index", waitIndex, nextWaitIndex)
			waitIndex = nextWaitIndex
			for _, value := range KVPairs.KVPairs {

				log.Println("ListenConsulKV: key and value ", string(value.Key), string(value.Value))
				newPolicy, err := processNewPolicy(value.Key, value.Value)
				if err != nil {
					log.Println("ListenConsulKV: processNewPolicy failed ", err)
				} else {
					ok := newPolicy.TakeDecision()
					if ok != true {
						log.Println("ListenConsulKV: TakeDecision on new policy ", newPolicy, " failed")
					}
				}

			}

		}

		<-time.After(20 * time.Second)
	}
}

//RuleThreshold
func GetCorrectRuleType(name string) interface{} {

	if strings.Contains(name, "Distance") == true {
		return &RuleDistance{}
	} else if strings.Contains(name, "Threshold") == true {
		return &RuleThreshold{}
	} else {
		log.Println("GetCorrectRuleType: No Rule Type found for", name)
	}
	return nil

}

//processNewPolicy if successfuly processed the policy will be store in a gloable map for
// further processing
// TODO: We need an representaion for the action taken by the policy.

func processNewPolicy(key string, data []byte) (*Policy, error) {
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

// Create a new Data
// TODO: but this ALLDCs.List will be changing all the time
// polcy will be working with the stale data....
func init() {
	//newData := &DCData{}

	log.Println("Policy Init called")

}

func (this *Policy) Len() int {
	return len(this.Rules)
}

//
func (this *Policy) Less(i, j int) bool {
	return this.Rules[i].Priority > this.Rules[j].Priority
}

func (this *Policy) Swap(i, j int) {
	this.Rules[i], this.Rules[j] = this.Rules[j], this.Rules[i]
}

// ApplyPriority sorts the rule based on priority
// A rule with priority 10 gets the max piority
//
func (this *Policy) ApplyPriority() bool {
	//sort the rules based on priorityA

	if this.Len() == 0 {
		log.Println("ApplyPriority: The rules array is nil")
		return false
	}
	fmt.Println(this)
	sort.Sort(this)
	fmt.Println("After", this)

	return true
}

func CreateMockRuleArray() []Rule {

	dummy := make([]Rule, 5)
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := range dummy {
		dummy[i].Priority = rand.Intn(10)
	}
	return dummy
}

func getTheConsulAndDCpointers() {

	for {

		log.Println("Lopping thorugh to get the global maps")

		common.ALLDCs.Lck.Lock()
		//newData.dcDataList = common.ALLDCs.List
		dcDataList = common.ALLDCs.List
		common.ALLDCs.Lck.Unlock()

		consulLib.GlobalConsulMutex.Lock()
		dcConsulList = consulLib.GlobalConsulDCInfo
		consulLib.GlobalConsulMutex.Unlock()

		if dcDataList != nil && dcConsulList != nil {
			log.Println("getTheConsulAndDCpointers: got", dcDataList, dcConsulList)
			break
		}

		<-time.After(2 * time.Second)
	}
	return
}

func Run(prefix string, config *common.ConsulConfig) {

	localwg := &sync.WaitGroup{}
	getTheConsulAndDCpointers()
	localwg.Add(1)
	go ListenConsulKV(prefix, localwg, config)
	localwg.Wait()
}
