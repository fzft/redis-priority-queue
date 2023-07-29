package redisPriorityQueue

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// RedisQueueClient Interface for RedisQueueClient,use generics to support different types of data
type RedisQueueClient[T any] interface {

	// PushOne Push an item to the queue
	PushOne(ctx context.Context, data T, priority int, key string) error

	// BatchPush Push a batch of items to the queue, use lua script to ensure atomicity and efficiency
	BatchPush(ctx context.Context, priority int, key string, data []T) error

	// PopOne Pop a message from the queue, use lua script to ensure atomicity and efficiency
	PopOne(ctx context.Context, key string, v *T) error

	// BatchPop Pop a batch of messages from the queue, use lua script to ensure atomicity and efficiency
	BatchPop(ctx context.Context, key string, count int, vs []*T) error

	// Sub Subscribe to the queue, return a channel to receive messages, use lua script to ensure atomicity and efficiency
	// every time poll we peek all the messages in the queue and return the highest priority message
	Sub(ctx context.Context, key string) (chan T, error)

	// UnSub Unsubscribe from the queue
	UnSub(ctx context.Context, key string) error
}

type redisQueueClient[T any] struct {
	Serializer
	logger Logger

	client       *redis.Client
	messageQueue chan []byte
}

// NewRedisQueueClient Create a new RedisQueueClient
// Params:
// - ctx: context.Context
// - serializer: SerializerType, eg: SerializerJson
// - client: *redis.Client
// Return:
// - RedisQueueClient
func NewRedisQueueClient[T any](_ context.Context, serializer SerializerType, client *redis.Client) RedisQueueClient[T] {
	return &redisQueueClient[T]{
		Serializer: NewSerializer(serializer),
		client:     client,
	}
}

// NewRedisQueueClientWithLogger Create a new RedisQueueClient with logger
// Params:
// - ctx: context.Context
// - serializer: SerializerType, eg: SerializerJson
// - client: *redis.Client
// - logger: Logger
// Return:
// - RedisQueueClient
func NewRedisQueueClientWithLogger[T any](_ context.Context, serializer SerializerType, client *redis.Client, logger Logger) RedisQueueClient[T] {
	return &redisQueueClient[T]{
		Serializer: NewSerializer(serializer),
		client:     client,
		logger:     logger,
	}
}

func (r *redisQueueClient[T]) PushOne(ctx context.Context, data T, priority int, key string) error {
	byteData, err := r.Serializer.Serialize(data)
	if err != nil {
		return err
	}
	script := redis.NewScript(ScriptPushOne)

	// Run the Lua script
	result, err := script.Run(ctx, r.client, []string{key}, priority, byteData).Result()
	if err != nil {
		return err
	}

	r.log("PushOne result: %v", result)
	return nil
}

func (r *redisQueueClient[T]) BatchPush(ctx context.Context, priority int, key string, data []T) error {

	// Convert data to byte slice, and insert priority before each data
	// eg. data = []string{"Alice", "Bob", "Charlie"} => args = []interface{}{priority, []byte{1, 65, 108, 105, 99, 101}, priority, []byte{1, 66, 111, 98}, priority, []byte{1, 67, 104, 97, 114, 108, 105, 101}}
	var args []interface{}
	for _, v := range data {
		b, err := r.Serializer.Serialize(v)
		if err != nil {
			return err
		}
		args = append(args, priority, b)
	}

	script := redis.NewScript(ScriptBatchPush)

	// Run the Lua script
	result, err := script.Run(ctx, r.client, []string{key}, args...).Result()
	if err != nil {
		return err
	}

	r.log("BatchPush result: %v", result)
	return nil
}

func (r *redisQueueClient[T]) PopOne(ctx context.Context, key string, v *T) error {
	script := redis.NewScript(ScriptPopOne)

	// Run the Lua script
	results, err := script.Run(ctx, r.client, []string{key}, 0, "+inf").Slice()
	if err != nil {
		return err
	}

	// Check if the result is empty
	if len(results) == 0 {
		return fmt.Errorf("result is empty")
	}

	r.log("PopOne result: %v", results[0])

	// Deserialize the message
	return r.Serializer.Deserialize([]byte(results[0].(string)), &v)
}

// BatchPop ...
// Params:
// - ctx: context.Context
// - key: string
// - count: int
// - vs: []*T, should be a pointer to a array, eg: []*string{}
func (r *redisQueueClient[T]) BatchPop(ctx context.Context, key string, count int, vs []*T) error {
	script := redis.NewScript(ScriptBatchPop)

	// Run the Lua script
	results, err := script.Run(ctx, r.client, []string{key}, 0, "+inf", count).Slice()
	if err != nil {
		return err
	}

	// Check if the result is empty
	if len(results) == 0 {
		return fmt.Errorf("result is empty")
	}

	r.log("BatchPop results: %v", results)

	for i, result := range results {
		// Deserialize the message
		if err = r.Serializer.Deserialize([]byte(result.(string)), &vs[i]); err != nil {
			return err
		}

	}

	return nil
}

func (r *redisQueueClient[T]) Sub(ctx context.Context, key string) (chan T, error) {
	//TODO implement me
	panic("implement me")
}

func (r *redisQueueClient[T]) UnSub(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func (r *redisQueueClient[T]) log(format string, args ...interface{}) {
	if r.logger != nil {
		switch r.logger.GetLevel() {
		case Info:
			r.logger.Info(format, args...)
		case Error:
			r.logger.Error(format, args...)
		case Debug:
			r.logger.Debug(format, args...)
		case Trace:
			r.logger.Trace(format, args...)
		default:
		}
	}
}
