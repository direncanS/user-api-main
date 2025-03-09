package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lambda-lama/user-api/models"
)

func GetWebcamDataByDateRange(topic string, startDate time.Time, endDate time.Time) ([]models.WebcamData, error) {
	conn, err := GetConnection()
	if err != nil {
		fmt.Print("Unable to connect to DB")
		return nil, errors.New("unable to connect to DB")
	}
	defer conn.Close(context.Background())
	fmt.Printf("Connection available: %s", conn.Config().Host)

	stmt := `SELECT id, topic, metadata FROM webcam_data WHERE topic = $1 AND created_at >= $2 AND created_at <= $3`

	rows, err := conn.Query(context.Background(), stmt, topic, startDate, endDate)
	if err != nil {
		fmt.Println("Unable to retrieve webcam data: ", err)
		return nil, err
	}
	defer rows.Close()

	var list []models.WebcamData

	for rows.Next() {
		var webcamData models.WebcamData
		err := rows.Scan(&webcamData.Id, &webcamData.Topic, &webcamData.Metadata)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		list = append(list, webcamData)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over result set:", err)
		return nil, err
	}

	return list, nil
}

func GetVideoDataByTopic(topic string) ([]models.VideoData, error) {
	conn, err := GetConnection()
	if err != nil {
		fmt.Print("Unable to connect to DB")
		return nil, errors.New("unable to connect to DB")
	}
	defer conn.Close(context.Background())
	fmt.Printf("Connection available: %s", conn.Config().Host)

	stmt := `SELECT id, folder, start_date, end_date FROM video_data WHERE topic = $1`

	rows, err := conn.Query(context.Background(), stmt, topic)
	if err != nil {
		fmt.Println("Unable to retrieve webcam data: ", err)
		return nil, err
	}
	defer rows.Close()

	var list []models.VideoData

	for rows.Next() {
		var videoData models.VideoData
		err := rows.Scan(&videoData.Id, &videoData.Folder, &videoData.StartDate, &videoData.EndDate)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		list = append(list, videoData)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error iterating over result set:", err)
		return nil, err
	}

	return list, nil
}

func GetAllTopics() ([]string, error) {
	conn, err := GetConnection()
	if err != nil {
		fmt.Print("Unable to connect to DB")
		return nil, errors.New("unable to connect to DB")
	}
	defer conn.Close(context.Background())
	fmt.Printf("Connection available: %s\n", conn.Config().Host)

	stmt := `SELECT DISTINCT topic FROM webcam_data`

	rows, err := conn.Query(context.Background(), stmt)
	if err != nil {
		fmt.Println("Unable to insert ", err)
		return nil, err
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topic string
		err := rows.Scan(&topic)
		if err != nil {
			fmt.Print(err)
		}
		topics = append(topics, topic)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	return topics, nil
}

func GetTopicData(topic string) ([]models.WebcamData, error) {
	conn, err := GetConnection()
	if err != nil {
		fmt.Println("Unable to connect to DB")
		return nil, errors.New("unable to connect to DB")
	}
	defer conn.Close(context.Background())
	fmt.Printf("Connection available: %s", conn.Config().Host)

	stmt := `SELECT id, metadata, created_at FROM webcam_data WHERE topic=$1`

	rows, err := conn.Query(context.Background(), stmt, topic)
	if err != nil {
		fmt.Println("Unable to select", err)
		return nil, err
	}
	defer rows.Close()

	var list []models.WebcamData
	for rows.Next() {
		var webcamData models.WebcamData
		err := rows.Scan(&webcamData.Id, &webcamData.Metadata, &webcamData.CreateAt)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		list = append(list, webcamData)
	}

	return list, nil
}
