package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/lambda-lama/user-api/aws"
	"github.com/lambda-lama/user-api/db"
	"github.com/lambda-lama/user-api/redis"
)

const (
	redisCacheTTL = 5 * time.Minute
)

// get topic data between two dates
func GetByTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	topic := vars["topic"]
	if topic == "" {
		sendError(w, http.StatusBadRequest, "Missing 'topic' in path")
		return
	}

	startDateStr := r.URL.Query().Get("start")
	endDateStr := r.URL.Query().Get("end")
	if startDateStr == "" && endDateStr == "" {
		cache, _ := redis.Get(topic)
		if cache != "" {
			sendResponseRaw(w, http.StatusOK, json.RawMessage(cache))
			return
		}
		data, err := db.GetTopicData(topic)
		if err != nil {
			fmt.Printf("Error getting topic data: %v", err)
			sendError(w, http.StatusInternalServerError, "Error getting data")
			return
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Error marshaling: %v", err)
			sendError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		redis.Set(topic, jsonData, redisCacheTTL)
		sendResponseRaw(w, http.StatusOK, jsonData)
		return
	}

	var startDate, endDate time.Time
	var err error
	if startDateStr != "" {
		startDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Invalid start date format")
			return
		}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			sendError(w, http.StatusBadRequest, "Invalid end date format")
			return
		}
	}

	cache, _ := redis.Get(topic + startDateStr + endDateStr)
	if cache != "" {
		sendResponseRaw(w, http.StatusOK, json.RawMessage(cache))
		return
	}
	data, err := db.GetWebcamDataByDateRange(topic, startDate, endDate)

	if err != nil {
		fmt.Printf("Error getting topic data: %v", err)
		sendError(w, http.StatusInternalServerError, "Error getting data")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	redis.Set(topic, jsonData, redisCacheTTL)
	sendResponseRaw(w, http.StatusOK, jsonData)
}

// get video by folder
func GetVideoFromFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	folder := vars["folder"]
	if folder == "" {
		sendError(w, http.StatusBadRequest, "Missing 'folder' in path")
		return
	}

	cache, _ := redis.Get(folder)
	if cache != "" {
		sendResponseRaw(w, http.StatusOK, json.RawMessage(cache))
		return
	}

	data, err := aws.GetVideoByFolder(folder)

	if err != nil {
		fmt.Printf("Error getting topic data: %v", err)
		sendError(w, http.StatusInternalServerError, "Error getting data")
		return
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"links": data,
	})
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	redis.Set(folder, jsonData, redisCacheTTL)
	sendResponseRaw(w, http.StatusOK, jsonData)
}

// get all videos for topic
func GetVideosByTopic(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	if topic == "" {
		sendError(w, http.StatusBadRequest, "Missing 'topic' parameter")
		return
	}

	cache, _ := redis.Get("videos/" + topic)
	if cache != "" {
		sendResponseRaw(w, http.StatusOK, json.RawMessage(cache))
		return
	}

	data, err := db.GetVideoDataByTopic(topic)
	if err != nil {
		fmt.Printf("Error getting topic data: %v", err)
		sendError(w, http.StatusInternalServerError, "Error getting data")
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling: %v", err)
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	redis.Set("videos/"+topic, jsonData, redisCacheTTL)
	sendResponseRaw(w, http.StatusOK, jsonData)
}

// get all topics
func GetDataTopics(w http.ResponseWriter, r *http.Request) {
	data, err := db.GetAllTopics()
	if err != nil {
		fmt.Printf("Error getting topic data: %v\n", err)
		sendError(w, http.StatusInternalServerError, "Error getting data")
		return
	}

	cache, _ := redis.Get("images")
	if cache != "" {
		sendResponseRaw(w, http.StatusOK, json.RawMessage(cache))
		return
	}

	fmt.Println(data)
	jsonData, err := json.Marshal(struct {
		Topics []string `json:"topics"`
	}{Topics: data})

	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		sendError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	redis.Set("images", jsonData, redisCacheTTL)
	sendResponseRaw(w, http.StatusOK, jsonData)
}
