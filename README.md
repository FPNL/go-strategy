# go-strategy
Assemble strategy for GO, ex: Exponential backoff, API Group Sender, Rollkey

# Go Helpers

# Install
```bash 
go get -u github.com/FPNL/go-strategy
```

# Sync API
APIs that need to be called in sequence, if one of them fails, the rollback API will be called in reverse order.
```go
groupAPI := NewGroupAPI(
    WithWaiStrategy(WaitExponentialBackoff(32)),
    WithMaxRetry(maxRetry),
)

// If mockAPI2 fail, mockAPI1 and mockAPI2 will be executed.
apiErr, rollbackErr := groupAPI.
    Then(mockAPI1, mockRollbackAPI1).
    Then(mockAPI2, mockRollbackAPI2).
    Done()
```

# Roll Key
Make your API Credentials Rotational.
This is a simple helper that allows you to rotate your API credentials.
Let's say there is a free API key limit 5 request per second. And to avoid this limitation is to create multiple API keys.
and use them in rotation.
```go
APIKeys := []string{"api-key-1", "api-key-2"}
requestTimes := 50
rate := 2

// when
keys, err := NewRotationalSlice(APIKeys, rate)

eg := errgroup.Group{}

for i := 0; i < requestTimes; i++ {
    eg.Go(func() error {
        key, err := keys.Get(context.TODO())
        if err != nil {
            return err
        }
        
        return nil
    })
}
err = eg.Wait()
```
