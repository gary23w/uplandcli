package models

type DataPackageBLOCK struct {
	Type string 
	ID  string
}

type APIRespBlockchain struct {
	QueryTimeMs float64 `json:"query_time_ms"`
	Cached      bool    `json:"cached"`
	Lib         int     `json:"lib"`
	Total       struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	Actions []struct {
		AtTimestamp string `json:"@timestamp"`
		Timestamp   string `json:"timestamp"`
		BlockNum    int    `json:"block_num"`
		TrxID       string `json:"trx_id"`
		Act         struct {
			Account       string `json:"account"`
			Name          string `json:"name"`
			Authorization []struct {
				Actor      string `json:"actor"`
				Permission string `json:"permission"`
			} `json:"authorization"`
			Data struct {
				A54  string `json:"a54"`
				A45  string `json:"a45"`
				P11  string `json:"p11"`
				P3   string `json:"p3"`
				P31  bool   `json:"p31"`
				P32  bool   `json:"p32"`
				P33  bool   `json:"p33"`
				P51  string `json:"p51"`
				P45  string `json:"p45"`
				Memo string `json:"memo"`
			} `json:"data"`
		} `json:"act"`
		Notified             []string `json:"notified"`
		CPUUsageUs           int      `json:"cpu_usage_us,omitempty"`
		NetUsageWords        int      `json:"net_usage_words,omitempty"`
		GlobalSequence       int64    `json:"global_sequence"`
		Producer             string   `json:"producer"`
		ActionOrdinal        int      `json:"action_ordinal"`
		CreatorActionOrdinal int      `json:"creator_action_ordinal"`
		AccountRAMDeltas     []struct {
			Account string `json:"account"`
			Delta   int    `json:"delta"`
		} `json:"account_ram_deltas,omitempty"`
		Receiver string `json:"receiver,omitempty"`
	} `json:"actions"`
}