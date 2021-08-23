package service

import (
	"crud/adapter"
	"crud/config"
	"crud/constant"
	"crud/dto/requests"
	"crud/dto/responses"
	"crud/entity"
	"crud/repository"
	"strings"

	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

func FetchCovidDataForIndiaAndState(location *requests.Location) (response *responses.FetchCasesResponse, err error) {
	state, err := adapter.GetStateFromLocation(location.Lat, location.Long)
	if err != nil {
		return nil, err
	}
	log.Println("user's state according to location is: ", state)
	response, err = fetchCovidData(state)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func UpdateCovidData(request *requests.UpdateCovidCasesRequest) (covidCases *responses.FetchCasesResponse, err error) {
	errorRes := errors.New("Exception while saving covid Data in DB")
	session, err := config.MongoClient.StartSession()
	if err != nil {
		log.Println(err)
		return nil, errorRes
	}
	if err = session.StartTransaction(); err != nil {
		log.Println(err)
		return nil, errorRes
	}
	ctx := context.Background()
	var stateCovidCases, indiaCovidCases *entity.CovidCases
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		stateCovidCases, err = repository.FindCovidDataByState(request.Region)
		if err != nil {
			log.Println("Exception while fetching covid cases for", request.Region, err)
			return errorRes
		}
		if stateCovidCases == nil {
			if request.Change < 0 {
				log.Println("cases can not be negative")
				return errors.New("Total cases can not be negative")
			}
			log.Printf("data for state: %s not found. creating new doc", request.Region)
			stateCovidCases = &entity.CovidCases{Region: request.Region, ActiveCases: request.Change}
		} else {
			stateCovidCases.ActiveCases += request.Change
		}
		currTime := time.Now()
		stateCovidCases.ModifiedAt = &currTime
		_, err := repository.UpsertCovidData(&sc, stateCovidCases)
		if err != nil {
			return errorRes
		}
		indiaCovidCases, err = repository.FindCovidDataByState(constant.INDIA_NAME)
		if err != nil {
			sc.AbortTransaction(sc)
			return errorRes
		}
		if indiaCovidCases == nil {
			indiaCovidCases = &entity.CovidCases{Region: constant.INDIA_NAME}
		}
		indiaCovidCases.ModifiedAt = &currTime
		indiaCovidCases.ActiveCases += request.Change
		_, err = repository.UpsertCovidData(&sc, indiaCovidCases)
		if err != nil {
			session.AbortTransaction(sc)
			return errorRes
		}
		if err = session.CommitTransaction(sc); err != nil {
			log.Println("Exception while commiting the upsert changes")
			return errorRes
		}
		return nil
	}); err != nil {
		log.Println(err)
		return nil, err
	}
	session.EndSession(ctx)

	conn := config.RedisPool.Get()
	err = conn.Send("MULTI")
	if err == nil {
		saveCovidDataInCache(&conn, stateCovidCases)
		saveCovidDataInCache(&conn, indiaCovidCases)
	} else {
		log.Println("exception while saving in cache")
	}
	_, err = conn.Do("EXEC")
	if err != nil {
		log.Println("exception while executing in cache")
	}

	log.Println(stateCovidCases, indiaCovidCases)
	response := responses.FetchCasesResponse{Result: [2]entity.CovidCases{*stateCovidCases, *indiaCovidCases}}
	return &response, nil
}

func fetchCovidData(state string) (data *responses.FetchCasesResponse, err error) {
	resError := errors.New("Error while fetching covid data")
	stateCovidCases, indiaCovidCases := getCovidCasesForStateAndIndiaFromCache(state)
	if stateCovidCases == nil || indiaCovidCases == nil {
		stateCovidCases, indiaCovidCases = fetchLiveCovidDataFromAPI(&state)
	}
	if stateCovidCases == nil {
		stateCovidCases, err = repository.FindCovidDataByState(state)
		if err != nil {
			return nil, resError
		}
		if stateCovidCases == nil {
			stateCovidCases = &entity.CovidCases{Region: state, Msg: "can not find covid data."}
		}
	}
	if indiaCovidCases == nil {
		indiaCovidCases, err = repository.FindCovidDataByState(constant.INDIA_NAME)
		if err != nil {
			return nil, resError
		}
		if indiaCovidCases == nil {
			indiaCovidCases = &entity.CovidCases{Region: constant.INDIA_NAME, Msg: "can not find covid data."}
		}
	}
	response := responses.FetchCasesResponse{Result: [2]entity.CovidCases{*stateCovidCases, *indiaCovidCases}}
	return &response, nil
}

func getCovidCasesForStateAndIndiaFromCache(region string) (stateCovidCases, indiaCovidCases *entity.CovidCases) {
	stateCovidCases = getCovidDataForStateFromCache(region)
	indiaCovidCases = getCovidDataForStateFromCache(constant.INDIA_NAME)
	return stateCovidCases, indiaCovidCases
}

func saveCovidDataInCache(conn *redis.Conn, data *entity.CovidCases) {
	if conn == nil {
		r := config.RedisPool.Get()
		conn = &r
		defer r.Close()
	}

	serialized, err := json.Marshal(data)
	if err != nil {
		log.Println("exception while serializing covid data in cache", err)
		return
	}
	if err == nil {
		_, err = (*conn).Do("SET", data.Region, serialized)
		if err != nil {
			log.Println("exception while saving covid data in cache", err)
			return
		}
		_, _ = (*conn).Do("EXPIRE", data.Region, 1800)
	}
}

func getCovidDataForStateFromCache(region string) *entity.CovidCases {
	conn := config.RedisPool.Get()
	defer conn.Close()
	value, err := redis.Bytes(conn.Do("GET", region))
	if err != nil {
		log.Println("exception while fetching covid data from cache for state", region, err)
		return nil
	} else if value == nil {
		return nil
	}
	var covidCases entity.CovidCases
	err = json.Unmarshal(value, &covidCases)
	if err != nil {
		log.Println("exception while converting covid data", err)
		return nil
	}
	return &covidCases
}

func fetchLiveCovidDataFromAPI(state *string) (stateCovidCases, indiaCovidCases *entity.CovidCases) {
	liveCovidDataResponse, _ := adapter.FetchLiveCovidData()
	if liveCovidDataResponse != nil {
		currTime := time.Now()
		indiaCovidCases = &entity.CovidCases{Region: constant.INDIA_NAME, ModifiedAt: &currTime, ActiveCases: liveCovidDataResponse.ActiveCases}
		liveCovidDataResponse.RegionData = append(liveCovidDataResponse.RegionData, indiaCovidCases)
		for i := 0; i < len(liveCovidDataResponse.RegionData); i++ {
			liveCovidDataResponse.RegionData[i].ModifiedAt = &currTime
			liveCovidDataResponse.RegionData[i].Region = strings.ToUpper(liveCovidDataResponse.RegionData[i].Region)
			saveCovidDataInCache(nil, liveCovidDataResponse.RegionData[i])
			repository.UpsertCovidData(nil, liveCovidDataResponse.RegionData[i])
		}
		for i := 0; i < len(liveCovidDataResponse.RegionData); i++ {
			if liveCovidDataResponse.RegionData[i].Region == *state {
				stateCovidCases = liveCovidDataResponse.RegionData[i]
				break
			}
		}
	}
	return stateCovidCases, indiaCovidCases
}
