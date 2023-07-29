package redisPriorityQueue

import (
	"context"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var setup = func() (mr *miniredis.Miniredis, err error) {
	// Mock Redis server
	return miniredis.Run()
}

type ClientTest struct {
	testKey string

	testOneData   string
	testBatchData []string

	mr *miniredis.Miniredis
}

func TestClientAll(t *testing.T) {
	// setup function
	mr, err := setup()
	defer mr.Close()

	if err != nil {
		os.Exit(1)
	}

	ct := ClientTest{mr: mr, testKey: "testKey", testOneData: "testOneData",
		testBatchData: []string{"testBatchData1", "testBatchData2"}}

	ct.TestPushOne(t)
	ct.TestPopOne(t)
	ct.TestBatchPush(t)
	ct.TestBatchPop(t)

}

func (ct ClientTest) TestPushOne(t *testing.T) {

	// Setup redisQueueClient
	rdb := redis.NewClient(&redis.Options{Addr: ct.mr.Addr()})
	rqc := NewRedisQueueClient[string](context.Background(), SerializerJson, rdb)

	// Test PushOne
	err := rqc.PushOne(context.Background(), ct.testOneData, 1, ct.testKey)
	assert.Nil(t, err)

	// Check if testData has been stored in testKey
	storedData, _ := ct.mr.ZMembers(ct.testKey)
	assert.NotEmpty(t, storedData)

	var storedDataString string
	json.Unmarshal([]byte(storedData[0]), &storedDataString)

	assert.Equal(t, ct.testOneData, storedDataString)
}

func (ct ClientTest) TestPopOne(t *testing.T) {

	// Setup redisQueueClient
	rdb := redis.NewClient(&redis.Options{Addr: ct.mr.Addr()})
	rqc := NewRedisQueueClient[string](context.Background(), SerializerJson, rdb)

	var resultData string
	err := rqc.PopOne(context.Background(), ct.testKey, &resultData)
	assert.Nil(t, err)

	// Check if testData has been stored in testKey
	storedData, _ := ct.mr.ZMembers(ct.testKey)
	assert.Empty(t, storedData)
	assert.Equal(t, ct.testOneData, resultData)
}

func (ct ClientTest) TestBatchPush(t *testing.T) {

	// Setup redisQueueClient
	rdb := redis.NewClient(&redis.Options{Addr: ct.mr.Addr()})
	rqc := NewRedisQueueClient[string](context.Background(), SerializerJson, rdb)

	err := rqc.BatchPush(context.Background(), 1, ct.testKey, ct.testBatchData)
	assert.Nil(t, err)

	// Check if testData has been stored in testKey
	storedData, _ := ct.mr.ZMembers(ct.testKey)

	result := make([]string, 0, len(storedData))
	for _, data := range storedData {
		var tmp string
		_ = json.Unmarshal([]byte(data), &tmp)
		result = append(result, tmp)
	}

	assert.NotEmpty(t, storedData)
	assert.Equal(t, ct.testBatchData, result)
}

func (ct ClientTest) TestBatchPop(t *testing.T) {

	// Setup redisQueueClient
	rdb := redis.NewClient(&redis.Options{Addr: ct.mr.Addr()})
	rqc := NewRedisQueueClient[string](context.Background(), SerializerJson, rdb)

	var count = len(ct.testBatchData)
	resultData := make([]*string, count, count)
	err := rqc.BatchPop(context.Background(), ct.testKey, count, resultData)
	assert.Nil(t, err)

	// Check if testData has been stored in testKey
	storedData, _ := ct.mr.ZMembers(ct.testKey)
	assert.Empty(t, storedData)

	// Dereference the pointers to get a []string
	actual := make([]string, len(resultData))
	for i, ptr := range resultData {
		actual[i] = *ptr
	}

	assert.Equal(t, ct.testBatchData, actual)
}
