package req

import (
	"encoding/json"
	"eos_bot/internal/models"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"eos_bot/internal/database"
	"eos_bot/internal/utils"
)

type EOSHTTPREQ struct {
	Bypass_sql bool
}

func (x *EOSHTTPREQ) HttpClient(req *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (x *EOSHTTPREQ) HttpEOSBasicRequest() models.APIRespBlockchain {
		var RespObj models.APIRespBlockchain
		req, err := http.NewRequest("GET", "https://eos.hyperion.eosrio.io/v2/history/get_actions?account=playuplandme&skip=0&limit=100&sort=desc", nil)
		if err != nil {
			log.Fatal(err)
		}
		res, err := x.HttpClient(req)
		if utils.IsJSON(string(res)) {
			err = json.Unmarshal([]byte(res), &RespObj)
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second * 5)
		} else {
			time.Sleep(time.Second * 1)
		}
		
		return RespObj
}

func (x *EOSHTTPREQ) HttpEOSgetAddress(prop_id string) models.UplandPropData {
	req, err := http.NewRequest("GET", "https://api.upland.me/properties/" + prop_id, nil)
	if err != nil {
		log.Fatal(err)
	}
	res, err := x.HttpClient(req)
	if err != nil {
		log.Fatal(err)
	}
	var PropData models.UplandPropData
	err = json.Unmarshal([]byte(res), &PropData)
	if err != nil {
		log.Printf("Error: ", err)
	}
	return PropData
}

func (x *EOSHTTPREQ) HttpEOSRespParser(req models.APIRespBlockchain) []models.DataPackageBLOCK {
	var myList []models.DataPackageBLOCK
	for _, v := range req.Actions {
		if v.Act.Name == "n2" {
			PropData := x.HttpEOSgetAddress(v.Act.Data.A45)
			myList = append(myList, models.DataPackageBLOCK{
				Type: v.Act.Name,
				ID: v.Act.Data.A45,
				Address: PropData.FullAddress,
				Lat: PropData.Centerlat,
				Long: PropData.Centerlng,
				UPX: PropData.OnMarket.Token,
				FIAT: PropData.OnMarket.Fiat,
			})
		}
	}
	if len(myList) == 0 {
		myList = append(myList, models.DataPackageBLOCK{
			Type: "NULL-DATA",
			ID: "",
		})
	}
	return myList
}

func (x *EOSHTTPREQ) CollectJsonFromAPI() [] models.DataPackageBLOCK {
	respObj := x.HttpEOSBasicRequest()
	parseDetails := x.HttpEOSRespParser(respObj)
	if !x.Bypass_sql && parseDetails[0].Type != "NULL-DATA" {
		database.EOSDatabaseMan("add", parseDetails)
	}
	return parseDetails
}