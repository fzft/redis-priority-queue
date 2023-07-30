## Redis Priority Queue
The redis-priority-queue package provides an interface RedisQueueClient for interacting with a priority queue implemented in Redis. This package utilizes the power of Redis and Lua scripting to ensure atomicity and efficiency.

### Requirements

- Go 1.19 or higher, use generics
- Redis 2.6.0 or higher

### Installation

```bash
$ go get github.com/fzft/redis-priority-queue
```

### Usage

Below is an example of how to use RedisQueueClient:

```go
    package main
    import  pq "github.com/fzft/redis-priority-queue"
    
    ctx := context.Background()

    // Create a new Redis client
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    // Create a new RedisQueueClient
    queueClient := pq.NewRedisQueueClient[string](ctx, redis_priority_queue.SerializerJson, client)

    // Push an item to the queue
    err := queueClient.PushOne(ctx, "testData", 1, "testKey")
    if err != nil {
        panic(err)
    }

    // Pop an item from the queue
    var data string
    err = queueClient.PopOne(ctx, "testKey", &data)
    if err != nil {
        panic(err)
    }
    fmt.Println(data)
```


### Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. Please make sure to update tests as appropriate.

### Upcoming Features
We're actively working on enhancing our project and have exciting features planned for future releases. Here's a sneak peek into what's coming next:

1. **Feature One**: Add redis stream feature to provide a mechanism to ensure that every message is delivered at least once.
this feature relies on the redis stream feature, which is only available in redis 5.0 or higher.
   ensuring that messages are processed exactly once (idempotently) is up to the application.

