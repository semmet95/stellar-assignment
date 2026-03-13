package cache

import (
	"log"
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

// init sets cache duration and initializes the cache map.
func init() {
	cd, ok := os.LookupEnv("CACHE_DURATION_MINS")
	if !ok {
		cacheDuration = time.Duration(defaultCacheDuration) * time.Minute
	} else {
		duration, err := strconv.Atoi(cd)
		if err != nil {
			log.Printf("failed to parse cache duration: %v; switching to default %d minutes\n", err, defaultCacheDuration)
			cacheDuration = time.Duration(defaultCacheDuration) * time.Minute
		} else {
			cacheDuration = time.Duration(duration) * time.Minute
		}
	}
	cachedMeasurements = make(map[string]Measurement)
}

// UpdateCache updates the cache for a client.
func UpdateCache(clientID, assetID string, measurement *Measurement) {
	delete(cachedMeasurements, clientID+"_"+assetID)
	cachedMeasurements[clientID+"_"+assetID] = *measurement
}

// GetCachedMeasurement returns a cached measurement if valid.
func GetCachedMeasurement(clientID, assetID string) *Measurement {
	measurement, ok := cachedMeasurements[clientID+"_"+assetID]
	if !ok {
		return nil
	}

	cacheAge := time.Since(measurement.ReqTime)
	if cacheAge > cacheDuration {
		return nil
	}
	return &measurement
}
