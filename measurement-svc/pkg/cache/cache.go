package cache

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultCacheDuration = 5
)

var (
	cacheDuration      time.Duration
	cachedMeasurements map[string]Measurement
)

type Measurement struct {
	ReqTime     time.Time
	ActivePower int64
	Setpoint    int64
}

func init() {
	cd, ok := os.LookupEnv("CACHE_DURATION_MINS")
	if !ok {
		cacheDuration = time.Duration(defaultCacheDuration) * time.Minute
	} else {
		duration, err := strconv.Atoi(cd)
		if err != nil {
			fmt.Printf("failed to parse cache duration with error: %v\nswitching to default %d minutes\n", err, defaultCacheDuration)
			cacheDuration = time.Duration(defaultCacheDuration) * time.Minute
		} else {
			cacheDuration = time.Duration(duration) * time.Minute
		}
	}
}

func UpdateCache(clientID string, measurement *Measurement) {
	delete(cachedMeasurements, clientID)
	cachedMeasurements[clientID] = *measurement
}

func GetCachedMeasurement(cliendID string) *Measurement {
	measurement, ok := cachedMeasurements[cliendID]
	if !ok {
		return nil
	}

	cacheAge := time.Since(measurement.ReqTime)
	if cacheAge > cacheDuration {
		return nil
	}
	return &measurement
}
