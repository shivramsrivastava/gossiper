package consulLib

import (
	"log"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

// This be called when ever there is a new Consul client is created or updated
// TODO: We need to handle the update case properly
func (this *FederaConsulClient) UpdateAndSignalGlobalConsulInfo() {
	GlobalConsulMutex.Lock()
	GlobalConsulDCInfo = this
	GlobalConsulMutex.Unlock()
	// TODO: for now lets not signal
	// We use use as is, later we might have to do something like this
	// GlobalConsulDCChannel <- true
}

// This function will be only called be the consul Leader gossiper.
// Poll's the local store and updaet the other DC's
// StorePreFix Federa
// TODO:blockFetch to use cas for fetching only when a new change is avaliable in the KV store
// This function will be called only by Leader of the consul replicate
func (this *FederaConsulClient) PollAndUpdateKV(blockFetch bool) {

	log.Println("[INFO] PollAndUpdateKV: KV store dosent exist")
	var waitIndex uint64
	for {

		//read from local store
		log.Println("[INFO] PollAndUpdateKV: Data received", waitIndex)
		data, res, err := this.GetList(Prefix, waitIndex)
		if err == nil && data != nil {
			log.Println("[INFO] PollAndUpdateKV: Data received", waitIndex)
			if blockFetch == false {
				//This when set to true we need to block till the data changes
				waitIndex = res
			}
			this.UpdateKVAcrossDcs(data)
		} else {
			log.Println("[INFO] PollAndUpdateKV: KV store dosent exist")
		}
		<-time.After(5 * time.Second)
	}
}

func (this *FederaConsulClient) UpdateKVAcrossDcs(data *KVData) {
	//TODO: expensive lock coz of nested loop
	localWG := new(sync.WaitGroup)
	this.DClist.Lock()
	for dcName := range this.DClist.DCClientConnection {
		//update all the DC's KV store in seperate goroutines
		dc := this.DClist.DCClientConnection[dcName]
		localWG.Add(1)
		go func(wg *sync.WaitGroup, dc *FederaConsulClient) {
			defer wg.Done()
			for _, kvpair := range data.KVPairs {
				dc.PutData(kvpair.Key, kvpair.Value)
			}

		}(localWG, dc)
		time.Sleep(5 * time.Second)
	}
	localWG.Wait()
	this.DClist.Unlock()
}

//GetDataFromLocalKVStore this will be used by the non-leader gossiper module
func (this *FederaConsulClient) GetDataFromLocalKVStore(waitIndex uint64) (api.KVPairs, uint64, error) {

	kvData, resultIndex, err := this.GetList(Prefix, waitIndex)
	if err == nil && kvData != nil {
		log.Println("Data from the local store", waitIndex)

		for _, val := range kvData.KVPairs {
			log.Println("Data from the Local store", val.Key, string(val.Value))
		}
	} else {
		log.Println("Either the data is null or error", err, kvData)
		// if an error occurs return back the same WaitIndex so it re-tried
		return nil, waitIndex, nil
	}

	return kvData.KVPairs, resultIndex, nil

}

func (this *FederaConsulClient) GetData(key string) []byte {
	data, Meta, err := this.KV.Get(key, nil)

	if err != nil {
		log.Println("Err in KV Get")
	}

	log.Println(data, Meta)

	return data.Value
}

//WriteOptions can be used to override the datacenter

//key should have the prefix
func (this *FederaConsulClient) PutData(key string, value []byte) (error, bool) {
	//writeMeta
	_, err := this.Put(&api.KVPair{Key: Prefix + "/" + key, Value: value}, nil)
	if err != nil {
		log.Println("Error PutData: Unable to put data to", this.Name, err, *this.Client)
		return err, false
	}
	return nil, true
}

func (this *FederaConsulClient) GetList(prefix string, waitIndex uint64) (*KVData, uint64, error) {

	log.Println("GetList waitindex ", waitIndex)

	q := &api.QueryOptions{}
	q.WaitIndex = waitIndex
	q.WaitTime = time.Duration(time.Hour * 1)
	KeyValuelist, KeyValueMeta, err := this.List(Prefix+"/", q)
	if err != nil {
		log.Println("GetList failed", err, KeyValueMeta, this.Client)
		//we pass on the same index when it fails to fetch from the KV store
		// so that it can be re-tried
		return nil, waitIndex, err
	}
	newKVData := &KVData{}
	newKVData.KVPairs = KeyValuelist
	newKVData.QueryMeta = KeyValueMeta
	log.Println("GetList the total len of data", len(newKVData.KVPairs))
	for _, val := range newKVData.KVPairs {

		log.Println(string(val.Key), string(val.Value), val)
	}

	log.Println("GetList waitindex at END", KeyValueMeta.LastIndex)
	return newKVData, KeyValueMeta.LastIndex, nil
}
