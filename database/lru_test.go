package database

import (
	"fmt"
	"github.com/hdt3213/godis/config"
	"github.com/hdt3213/godis/database/eviction"
	"github.com/hdt3213/godis/lib/mem"
	"github.com/hdt3213/godis/lib/utils"
	"testing"
)

func TestLRUEvictionKey(t *testing.T) {
	testDB.Flush()
	setLRUConfig()
	marks := make([]eviction.KeyMark, 10)
	for i := 0; i < 10; i++ {

		marks[i] = eviction.KeyMark{
			Mark: int32(i),
			Key:  fmt.Sprint(i),
		}
	}
	s := testDB.evictionPolicy.Eviction(marks)
	if s != "0" {
		t.Errorf("eviction key is wrong")
	}
	config.Properties = &config.ServerProperties{}
}

func TestLRU(t *testing.T) {
	testDB.Flush()
	setLRUConfig()
	for i := 0; i < 1000; i++ {
		key := utils.RandString(10)
		value := utils.RandString(10)
		testDB.Exec(nil, utils.ToCmdLine("SET", key, value))

		if mem.GetMaxMemoryState(nil) {
			t.Errorf("memory out of config")
		}
	}
	config.Properties = &config.ServerProperties{}
}

func setLRUConfig() {
	config.Properties = &config.ServerProperties{
		Maxmemory:        3000,
		MaxmemoryPolicy:  "allkeys-lru",
		LfuLogFactor:     5,
		MaxmemorySamples: 5,
	}
	testDB.evictionPolicy = makeEvictionPolicy()
}
