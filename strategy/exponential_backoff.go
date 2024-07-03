package strategy

import (
	"math"
	"math/rand"
	"time"
)

// ExponentialBackoff
// An exponential backoff algorithm retries requests exponentially, increasing the waiting time between retries up to a
// maximum backoff time. An example is:
//
// 1. Make a request to Memorystore for Redis.
// 2. If the request fails, wait 1 + random_number_milliseconds seconds and retry the request.
// 3. If the request fails, wait 2 + random_number_milliseconds seconds and retry the request.
// 4. If the request fails, wait 4 + random_number_milliseconds seconds and retry the request.
// 5. And so on, up to a maximum_backoff time.
// 6. Continue waiting and retrying up to some maximum number of retries, but do not increase the wait period between
// retries.
//
// where:
//
// - The wait time is min(((2^n)+random_number_milliseconds), maximum_backoff), with n incremented by 1 for each
// iteration (request).
// - random_number_milliseconds is a random number of milliseconds less than or equal to 1000. This helps to avoid cases
// where many clients get synchronized by some situation and all retry at once, sending requests in synchronized waves.
// The value of random_number_milliseconds should be recalculated after each retry request.
// - maximum_backoff is typically 32 or 64 seconds. The appropriate value depends on the use case.
//
// It's okay to continue retrying once you reach the maximum_backoff time. Retries after this point do not need to
// continue increasing backoff time. For example, if a client uses an maximum_backoff time of 64 seconds, then after
// reaching this value, the client can retry every 64 seconds. At some point, clients should be prevented from retrying
// infinitely.
//
// The maximum backoff and maximum number of retries that a client uses depends on the use case and network conditions.
// For example, mobile clients of an application may need to retry more times and for longer intervals when compared to
// desktop clients of the same application
//
// reference: https://cloud.google.com/memorystore/docs/redis/exponential-backoff
func ExponentialBackoff(n int, maximumBackoff float64) time.Duration {
	x := 1 << n
	randomNumber := float64(x) + rand.Float64()

	waitTime := math.Min(randomNumber, maximumBackoff)

	return time.Duration(waitTime) * time.Second
}
