package policylib

import (
	"log"
	"sort"

	"../common"
)

type RuleThreshold struct {
	RecosurceLimit float64
}

//this function can be put int the common code
//We dont need to match local only but all dcs
// dcName from policy decision
func GetMatchingConsulDCInfo(dcName string) (*common.DC, bool) {

	//dcDataList cant be null here
	localDCdata, ok := dcDataList[dcName]
	if true != ok {
		log.Println("GetMatchingConsulDCInfo: there is not matching DC in common DC map")
		return nil, false
	}
	return localDCdata, true
}

func (this *RuleThreshold) GetnewDCarrangment(policydecision *PolicyDecision) bool {

	for index, gossiperDCName := range policydecision.SortedDCName {

		gossiperDcInfo, ok := GetMatchingConsulDCInfo(gossiperDCName)

		if ok != true {
			log.Println("IsCurrentDCWithInPolicyThershold: failes to get the dc info ")
			continue
		}

		var cpuPercentage, memPercentage, diskPercentage float64

		if gossiperDcInfo.CPU == 0 {
			log.Println("IsCurrentDCWithInPolicyThershold: Cannot apply policy since the total CPU in DC is nil")
			return false
		} else if gossiperDcInfo.MEM == 0 {
			log.Println("IsCurrentDCWithInPolicyThershold: Cannot apply policy since the total MEM in DC is nil")
			return false

		} else if gossiperDcInfo.DISK == 0 {
			log.Println("IsCurrentDCWithInPolicyThershold: Cannot apply policy since the total CPU in DISK is nil")
			return false

		}
		cpuPercentage = ((gossiperDcInfo.Ucpu / gossiperDcInfo.CPU) * 100)
		memPercentage = ((gossiperDcInfo.Umem / gossiperDcInfo.MEM) * 100)
		diskPercentage = ((gossiperDcInfo.Udisk / gossiperDcInfo.DISK) * 100)

		//TODO: not the right place to decide
		/*if cpuPercentage >= this.RecosurceLimit || memPercentage >= this.RecosurceLimit || diskPercentage >= this.RecosurceLimit {
			log.Println("IsCurrentDCWithInPolicyThershold: DC burst ", cpuPercentage, memPercentage, diskPercentage)
			return false
		}*/
		policydecision.SortValue[index] = (cpuPercentage + memPercentage + diskPercentage/3)
	}

	sort.Sort(policydecision)
	log.Println("")
	return true
}

//sort interface based on the SortValues sort the SortedDCName also
func (this *PolicyDecision) Len() int {
	return len(this.SortValue)
}

func (this *PolicyDecision) Swap(i, j int) {
	this.SortValue[i], this.SortValue[j] = this.SortValue[j], this.SortValue[i]
	this.SortedDCName[i], this.SortedDCName[j] = this.SortedDCName[j], this.SortedDCName[i]

}

//sorts in ascending order
func (this *PolicyDecision) Less(i, j int) bool {

	return this.SortValue[i] < this.SortValue[j]

}

func (this *RuleThreshold) ApplyRule(policydecision *PolicyDecision) bool {

	log.Println("ApplyRule: ApplyRule of the interface to RuleThreshold")

	return this.GetnewDCarrangment(policydecision)
}
