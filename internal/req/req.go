package req

import (
	"encoding/json"
	"eos_bot/internal/models"
	"io/ioutil"
	"log"
	"net/http"

	"eos_bot/internal/database"
)

func httpClient(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// read response
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func httpEOSBasicRequest() models.APIRespBlockchain {
		var RespObj models.APIRespBlockchain
		req, err := http.NewRequest("GET", "https://eos.hyperion.eosrio.io/v2/history/get_actions?account=playuplandme&skip=0&limit=100&sort=desc", nil)
		if err != nil {
			log.Fatal(err)
		}
		res, err := httpClient(req)
		err = json.Unmarshal([]byte(res), &RespObj)
		if err != nil {
			log.Fatal(err)
		}
		return RespObj
}

func httpEOSRespParser(req models.APIRespBlockchain) []models.DataPackageBLOCK {
	var myList []models.DataPackageBLOCK
	for _, v := range req.Actions {
		if v.Act.Name == "n2" {
			myList = append(myList, models.DataPackageBLOCK{
				Type: v.Act.Name,
				ID: v.Act.Data.A45,
				Address: "NULL",
				Price: "NULL",
			})
		}
	}
	// add list to database
	database.AddPropertiesToDatabase(myList)
	if len(myList) == 0 {
		myList = append(myList, models.DataPackageBLOCK{
			Type: "No data found",
			ID: "",
		})
	}
	return myList
}

func CollectJsonFromAPI() [] models.DataPackageBLOCK {
	// do basic request
	respObj := httpEOSBasicRequest()
	// parse response
	parseDetails := httpEOSRespParser(respObj)
	// more data manipulation?
	//////
	return parseDetails
}