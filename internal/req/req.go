package req

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gary23w/uplandcli/internal/models"

	"github.com/gary23w/uplandcli/internal/database"
)

type EOSHTTPREQ struct {
	Bypass_sql bool
}

var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
	Timeout: 20 * time.Second,
}

func readResponseBody(resp *http.Response, maxBytes int64) ([]byte, error) {
	if resp.Body == nil {
		return nil, errors.New("empty response body")
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, maxBytes))
}

func (x *EOSHTTPREQ) do(req *http.Request) (*http.Response, error) {
	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := readResponseBody(resp, 256*1024)
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, string(b))
	}
	return resp, nil
}

func (x *EOSHTTPREQ) HttpClient(req *http.Request) (string, error) {
	resp, err := x.do(req)
	if err != nil {
		return "", err
	}
	body, err := readResponseBody(resp, 10*1024*1024)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (x *EOSHTTPREQ) HttpEOSBasicRequestContext(ctx context.Context) (models.APIRespBlockchain, error) {
	var RespObj models.APIRespBlockchain
	req, err := http.NewRequestWithContext(ctx, "GET", "https://eos.hyperion.eosrio.io/v2/history/get_actions?account=playuplandme&skip=0&limit=100&sort=desc", nil)
	if err != nil {
		return RespObj, err
	}
	resp, err := x.do(req)
	if err != nil {
		return RespObj, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(io.LimitReader(resp.Body, 10*1024*1024))
	if err := dec.Decode(&RespObj); err != nil {
		return RespObj, err
	}
	return RespObj, nil
}

func (x *EOSHTTPREQ) HttpEOSBasicRequest() models.APIRespBlockchain {
	resp, err := x.HttpEOSBasicRequestContext(context.Background())
	if err != nil {
		// keep legacy behavior: don't crash the CLI loop
		log.Println("http request failed:", err)
		time.Sleep(1 * time.Second)
		return models.APIRespBlockchain{}
	}
	// keep the original pacing
	time.Sleep(5 * time.Second)
	return resp
}

func (x *EOSHTTPREQ) HttpEOSgetAddressContext(ctx context.Context, prop_id string) (models.UplandPropData, error) {
	var PropData models.UplandPropData
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.upland.me/properties/"+prop_id, nil)
	if err != nil {
		return PropData, err
	}
	resp, err := x.do(req)
	if err != nil {
		return PropData, err
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(io.LimitReader(resp.Body, 10*1024*1024))
	if err := dec.Decode(&PropData); err != nil {
		return PropData, err
	}
	return PropData, nil
}

func (x *EOSHTTPREQ) HttpEOSgetAddress(prop_id string) models.UplandPropData {
	prop, err := x.HttpEOSgetAddressContext(context.Background(), prop_id)
	if err != nil {
		log.Println("property lookup failed:", err)
		return models.UplandPropData{}
	}
	return prop
}

func (x *EOSHTTPREQ) HttpEOSRespParserContext(ctx context.Context, req models.APIRespBlockchain) []models.DataPackageBLOCK {
	// Limit concurrency to avoid hammering upstream while still speeding up processing.
	const maxConcurrent = 4
	sem := make(chan struct{}, maxConcurrent)

	results := make([]models.DataPackageBLOCK, 0, len(req.Actions))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, v := range req.Actions {
		v := v
		if v.Act.Name != "n2" {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			prop, err := x.HttpEOSgetAddressContext(ctx, v.Act.Data.A45)
			if err != nil {
				return
			}
			item := models.DataPackageBLOCK{
				Type:    v.Act.Name,
				ID:      v.Act.Data.A45,
				Address: prop.FullAddress,
				Lat:     prop.Centerlat,
				Long:    prop.Centerlng,
				UPX:     prop.OnMarket.Token,
				FIAT:    prop.OnMarket.Fiat,
			}
			mu.Lock()
			results = append(results, item)
			mu.Unlock()
		}()
	}

	wg.Wait()
	if len(results) == 0 {
		return []models.DataPackageBLOCK{{Type: "NULL-DATA", ID: ""}}
	}
	return results
}

func (x *EOSHTTPREQ) HttpEOSRespParser(req models.APIRespBlockchain) []models.DataPackageBLOCK {
	return x.HttpEOSRespParserContext(context.Background(), req)
}

func (x *EOSHTTPREQ) CollectJsonFromAPIContext(ctx context.Context) ([]models.DataPackageBLOCK, error) {
	respObj, err := x.HttpEOSBasicRequestContext(ctx)
	if err != nil {
		return []models.DataPackageBLOCK{{Type: "NULL-DATA", ID: ""}}, err
	}
	parseDetails := x.HttpEOSRespParserContext(ctx, respObj)
	if !x.Bypass_sql && len(parseDetails) > 0 && parseDetails[0].Type != "NULL-DATA" {
		database.EOSDatabaseMan("add", parseDetails)
	}
	return parseDetails, nil
}

func (x *EOSHTTPREQ) CollectJsonFromAPI() []models.DataPackageBLOCK {
	data, err := x.CollectJsonFromAPIContext(context.Background())
	if err != nil {
		log.Println("api collection failed:", err)
		return []models.DataPackageBLOCK{{Type: "NULL-DATA", ID: ""}}
	}
	return data
}