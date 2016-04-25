package consulLib

import (
	"log"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
)

type KVData struct {
	api.KVPairs
	*api.QueryMeta
}

const Prefix = "Fedra"

//this file will have all the kv store warps/decorator
//this will have the list of ClientKV for the uderlying client

//poll the local store and updaet the other DC's
// StorePreFix Federa
// TODO:blockFetch to use cas for fetching only when a new change is avaliable in the KV store
// This function will be called only by Leader of the consul replicate
func (this *FederaConsulClient) PollAndUpdateKV(blockFetch bool) {

	log.Println("[INFO] PollAndUpdateKV: KV store dosent exist")
	for {

		//read from local store
		if data := this.GetList(Prefix); data != nil {
			log.Println("[INFO] PollAndUpdateKV: Data received")
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
	for _, dc := range this.DClist.DCClientConnection {
		//update all the DC's KV store in seperate goroutines
		localWG.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			for _, kvpair := range data.KVPairs {
				dc.PutData(kvpair.Key, kvpair.Value)
			}

		}(localWG)
		time.Sleep(5 * time.Second)
	}
	localWG.Wait()
	this.DClist.Unlock()
}

//GetDataFromLocalKVStore this will be used by the non-leader gossiper module
func (this *FederaConsulClient) GetDataFromLocalKVStore() {

	go func() {
		for {
			kvData := this.GetList(Prefix)
			if kvData != nil {
				log.Println("Data from the local store")
				for _, val := range kvData.KVPairs {
					log.Println("Data from the Local store", val.Key, string(val.Value))
				}
			}
			<-time.After(time.Second * 2)
		}
	}()
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

func (this *FederaConsulClient) GetList(prefix string) *KVData {

	newKVData := &KVData{}
	KeyValuelist, KeyValueMeta, err := this.List(Prefix+"/", nil)
	if err != nil {
		log.Println("GetList failed", err, KeyValueMeta, this.Client)
		return nil
	}

	newKVData.KVPairs = KeyValuelist
	newKVData.QueryMeta = KeyValueMeta
	log.Println("GetList the total len of data", len(newKVData.KVPairs))
	for _, val := range newKVData.KVPairs {

		log.Println(string(val.Key), string(val.Value), val)
	}

	return newKVData
}
