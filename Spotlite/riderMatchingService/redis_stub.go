package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"

	insight "github.com/apganesh/Insight_Project/Spotlite/common"
	"github.com/go-redis/redis"
	"github.com/golang/geo/s2"
)

var redisClient *redis.Client

func RedisNewClient() error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     insight.DefaultConfig.RedisClusters,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		log.Fatal("Error Initializing connection to redis: ")
		return err
	}

	// start with a clean slate ....
	redisClient.FlushDB()

	fmt.Println(pong, err)
	return nil
}

func addDriverToRedis(in *insight.Driver) {

	// Encode the data
	drvSer := encodeDriver(in)

	driverId := strconv.FormatInt(in.Id, 10)

	redisClient.HSet("drivers", driverId, drvSer)

	// Cleanup the old values
	redisClient.SRem("georadius", driverId)
	redisClient.ZRem("geofree", driverId)

	if in.Radius != 0.0 {
		redisClient.SAdd("georadius", driverId)
	} else {
		redisClient.GeoAdd("geofree", &redis.GeoLocation{Longitude: in.Lng, Latitude: in.Lat, Name: driverId})
	}
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getNearestDrivers(lat, lng float64) []insight.Driver {
	var drvrs []insight.Driver

	// get all the pts with radius and filter them
	cmd := redisClient.GeoRadius("geofree", lng, lat, &redis.GeoRadiusQuery{
		Radius:      5,
		WithDist:    true,
		WithCoord:   true,
		WithGeoHash: true,
		Count:       10,
		Sort:        "ASC",
	})
	res, err := cmd.Result()
	if err != nil {
		log.Println("Cannot find nearest drivers")
		return drvrs
	}
	fmt.Println("Total matches for nearest drivers : ", len(res))

	var ids []int64

	for _, r := range res {
		fmt.Println(r.Name, r.Latitude, r.Longitude, r.Dist)
		id, _ := strconv.ParseInt(r.Name, 10, 64)
		ids = append(ids, id)
	}
	myres := getDriversFromIds(ids)
	return myres
}

////////////////////////////////////////////////////////////////////
// Decoder and Encoder for Driver from/to Redis
////////////////////////////////////////////////////////////////////

func encodeDriver(driver *insight.Driver) []byte {
	var res bytes.Buffer
	enc := gob.NewEncoder(&res)
	enc.Encode(driver)
	return res.Bytes()
}

func decodeDriver(in interface{}, driver *insight.Driver) {
	buf := bytes.NewBufferString(in.(string))
	dec := gob.NewDecoder(buf)
	dec.Decode(driver)
}

// Get the drivers from Redis with ids
func getDriversFromIds(ids []int64) []insight.Driver {
	var sids []string

	for _, id := range ids {
		sids = append(sids, strconv.FormatInt(id, 10))
	}

	mres := redisClient.HMGet("drivers", sids...)
	var drivers []insight.Driver
	res, _ := mres.Result()

	for _, r := range res {
		var d insight.Driver
		decodeDriver(r, &d)
		drivers = append(drivers, d)
	}
	return drivers
}

func getOverlappingDrivers(ll insight.LatLng, drivers []insight.Driver) []insight.Driver {
	var res []insight.Driver

	rll := s2.LatLngFromDegrees(ll.Lat, ll.Lng)

	for _, di := range drivers {
		if di.Status == 2 {
			continue
		}

		dll := s2.LatLngFromDegrees(di.Lat, di.Lng)
		l := s2.RegularLoop(s2.PointFromLatLng(dll), kmToAngle(di.Radius), 20)

		cont := l.ContainsPoint(s2.PointFromLatLng(rll))
		if cont == true {
			res = append(res, di)
			//fmt.Println("Added driver from geofence: ", di.Id, di.Lat, di.Lng)
		}
	}
	return res
}

func getOverlapDrivers(in *insight.Rider) []insight.Driver {
	res := redisClient.SMembers("georadius")
	var dids []int64

	for _, r := range res.Val() {
		id, _ := strconv.ParseInt(r, 10, 64)
		dids = append(dids, id)
	}

	inDrivers := getDriversFromIds(dids)
	fmt.Println("Total drivers with geofences ", len(inDrivers))
	ovdrivers := getOverlappingDrivers(insight.LatLng{in.SLat, in.SLng}, inDrivers)
	fmt.Println("Total drivers with overlap geofences", len(ovdrivers))
	return ovdrivers
}

////////////////////////////////////////////////////////////////////
// sort the drivers by haversine
////////////////////////////////////////////////////////////////////

type driverDistMap struct {
	Dist float64
	Drv  insight.Driver
}

type driversByDist []driverDistMap

func (a driversByDist) Len() int           { return len(a) }
func (a driversByDist) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a driversByDist) Less(i, j int) bool { return a[i].Dist < a[j].Dist }

func sortByHaverDistance(r *insight.Rider, drivers []insight.Driver) []insight.Driver {
	var ddMap driversByDist

	for _, drvr := range drivers {
		if drvr.Status == 2 {
			continue
		}

		hdist := getHaversineDist(r.SLat, r.SLng, drvr.Lat, drvr.Lng)
		ddMap = append(ddMap, driverDistMap{Drv: drvr, Dist: hdist})
	}

	sort.Sort(driversByDist(ddMap))
	var res []insight.Driver
	for _, dd := range ddMap {
		res = append(res, dd.Drv)
	}
	return res
}

////////////////////////////////////////////////////////////////////
// get  the drivers from redis
////////////////////////////////////////////////////////////////////

func getDrivers_redis(in *insight.Rider) []insight.Driver {
	nearDrivers := getNearestDrivers(in.SLat, in.SLng)
	overlapDrivers := getOverlapDrivers(in)

	var res []insight.Driver
	res = append(res, nearDrivers...)
	res = append(res, overlapDrivers...)

	fmt.Println("Total driver neighbors : ", len(res))
	sortedDrivers := sortByHaverDistance(in, res)

	var cDrivers []insight.Driver

	mm := math.Min(5.0, float64(len(sortedDrivers)))
	for i := 0; i < int(mm); i++ {
		cDrivers = append(cDrivers, sortedDrivers[i])
	}

	return cDrivers
}
