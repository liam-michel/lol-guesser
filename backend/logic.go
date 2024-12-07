package main

import (
	"encoding/json"
	"log"
	"os"

	"math/rand/v2"
)

//one function to read in the json file
//one to generate the random index
//one to pick out the champion and its image from the url
//DRY

// takes a list of all champion names, returns the name and image url of the chosen champion
func generateRandomNumber(length int64) int64 {
	//generate a random number between 0 and the length of the champion array
	randomIndex := rand.IntN(int(length))
	return int64(randomIndex)

}

func PickChampByName(name string, championsData *map[string]interface{}) (string, string, error) {
	//pick out the enemy champion from the champions data map
	champion := (*championsData)[name]
	imageUrl := champion.(map[string]interface{})["image"].(map[string]interface{})["full"].(string)
	return name, imageUrl, nil
}

func getChampionNames(championsData *map[string]interface{}) ([]string, error) {
	champions := make([]string, 0, len(*championsData))
	for champion := range *championsData {
		champions = append(champions, champion)
	}
	return champions, nil
}

func ReadChampionsJSON() (map[string]interface{}, error) {
	file, err := os.Open("./../static/champions.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data["data"].(map[string]interface{}), err

}

func PickRandomChampion() (name string, url string, err error) {
	//read in the json file
	championsData, err := ReadChampionsJSON()
	if err != nil {
		log.Fatal(err)
	}
	//get the champion names
	champion_names, err := getChampionNames(&championsData)
	if err != nil {
		log.Fatal(err)
	}
	//generate a random number
	randomIndex := generateRandomNumber(int64(len(champion_names)))
	//pick out the champion name and image url
	name, image_url, err := PickChampByName(champion_names[randomIndex], &championsData)
	if err != nil {
		log.Fatal(err)
	}
	return name, image_url, nil

}
