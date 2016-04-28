package policylib

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	"../common"
	"../consulLib"
	"github.com/hashicorp/consul/api"
)

var (
	dcDataList   map[string]*common.DC
	dcConsulList *FederaConsulClient
)

// Policy runs only when there is change in the DC info from the gosspier
// or when there is a change in the Consul KV "Fedra" store
func ListenAllDcDataChn() {

}

// ListenDcConsulList will block till the KV client connection is established
// TODO: We will read the Consul DC data and update the global pointer dcConsulList
// This has to be controlled via a mutex.
// But for now we dont lock while reading.
func ListenDcConsulList(){
	<-
}

// Its tricky to actually make a block call on the consul looking for a KV prefix
// If there is not such KV store we create one with the /Fedra with some data/no data
// we then use the modified index from that store for later block
func ListenConsulKV(string Prefix) {

	var waitIndex uint64

   // This is a blocking call since we need to watch the store for any new policy change
   // and notify processNewPloicy
	for{

		//GlobalConsulDCInfo should have valid consul client connections
	KVPairs, nextWaitIndex, err:=consulLib.GlobalConsulDCInfo.GetDataFromLocalKVStore(waitIndex)

	if err!=nil{
		log.Println("ListenConsulKV: Watch on local store failes",err)
		//TODO: We just continue for now
		continue
		}else{
			waitIndex= nextWaitIndex

			for ,_:=range KVPairs. {

			}

		}
}
}

func GetCorrectRuleType(name string) interface{

	if strings.Contains(name,"Distance"){
		return &RuleDistance{}
	}else if strings.Contains(name,"Distance"){
		return &RuleThreshold{}
	}else{
		log.Println("GetCorrectRuleType: No Rule Type found for",name)
	}
return nil

}

func processNewPolicy(data []byte) error{
	var tempPolicy Policy
	err:=json.Unmarshal(data, &tempPolicy)

	if err != nil{
		log.Println("processNewPolicy: Unable to unmarshal the processNewPolicy",err)
		return err
	}

for _,values := range tempPolicy.Rules {

	dummy:=FakeJsonRule{}
	ruleType:=GetCorrectRuleType(values.Name)
	if ruleType !=nil{
	dummy.Content = ruleType
	err:=json.Unmarshal(data, &dummy)
	if err!=nil{
		log.Println("processNewPolicy: Unable to get the content from the json")
	}
	tempPolicy.Content=dummy.Content

	}
}
}



// Create a new Data
// TODO: but this ALLDCs.List will be changing all the time
// polcy will be working with the stale data....
func init() {
	//newData := &DCData{}
	common.ALLDCs.Lck.Lock()
	newData.dcDataList = common.ALLDCs.List
	common.ALLDCs.Lck.UnLock()

    GlobalConsulMutex.Lock()
    newData.dcConsulList = GlobalConsulDCInfo
    GlobalConsulMutex.Unlock()
}

// This interface defines the behaviour of a rule
// All rule typr must implement this rule
type RuleInterface interface {
	ApplyRule([]DCData) []DCData
}

type PolicyConsul struct {
	*FederaConsulClient
}

// A Group of rules which forms a policy
//
type Policy struct {
	Rules []Rule
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

type FakeJsonRule struct {
	Content  interface{} // this represt the actual data for the rule, this can be of any type
}

// ApplyPriority sorts the rule based on priority
// A rule with priority 10 gets the max piority
//
func (this *Policy) ApplyPriority() {
	//sort the rules based on priority
	fmt.Println(this)
	sort.Sort(this)
	fmt.Println("After", this)

	//return nil
}

func CreateMockRuleArray() []Rule {

	dummy := make([]Rule, 5)
	rand.Seed(int64(time.Now().Nanosecond()))
	for i := range dummy {
		dummy[i].Priority = rand.Intn(10)
	}
	return dummy
}

func Run(prefix string) {



}

func main() {

	//simple := make([]Rule, 5)
	//s.Name = "Rule_One"
	var s Policy

	s.Rules = CreateMockRuleArray()

	file, err := os.Create("polic.json")
	if err != nil {
		log.Println(err)
	}

	x, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	file.Write(x)
	fmt.Println(string(x), err, file.Name())
	s.ApplyPriority()
	s.ApplyPriority()

}
