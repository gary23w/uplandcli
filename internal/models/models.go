package models

import "time"

type DataPackageBLOCK struct {
	Type    string
	ID      string
	Address string
	Lat     string
	Long    string
	UPX     string
	FIAT    string
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

type UplandPropData struct {
	PropID             int64       `json:"prop_id"`
	FullAddress        string      `json:"full_address"`
	Centerlat          string      `json:"centerlat"`
	Centerlng          string      `json:"centerlng"`
	Boundaries         string      `json:"boundaries"`
	Area               int         `json:"area"`
	Status             string      `json:"status"`
	LockedUntil        interface{} `json:"locked_until"`
	LockedUserID       interface{} `json:"locked_user_id"`
	TransactionID      string      `json:"transaction_id"`
	LastTransactionID  string      `json:"last_transaction_id"`
	LastPurchasedPrice int         `json:"last_purchased_price"`
	TerminalID         interface{} `json:"terminal_id"`
	Feature            interface{} `json:"feature"`
	Labels             struct {
		FsaAllow bool `json:"fsa_allow"`
	} `json:"labels"`
	OnMarket struct {
		Token            string `json:"token"`
		Fiat             string `json:"fiat"`
		AllowFiatOffers  bool   `json:"allow_fiat_offers"`
		AllowTokenOffers bool   `json:"allow_token_offers"`
		AllowPropOffers  bool   `json:"allow_prop_offers"`
		Currency         string `json:"currency"`
	} `json:"on_market"`
	OffersMadeByYou    interface{} `json:"offers_made_by_you"`
	OffersWithProperty interface{} `json:"offers_with_property"`
	OffersToProperty   interface{} `json:"offers_to_property"`
	FiatPrice          interface{} `json:"fiat_price"`
	IsOffering         bool        `json:"is_offering"`
	IsBlocked          bool        `json:"is_blocked"`
	Owner              string      `json:"owner"`
	OwnerUsername      string      `json:"owner_username"`
	YieldPerHour       float64     `json:"yield_per_hour"`
	Price              int         `json:"price"`
	CanMakeOffer       bool        `json:"can_make_offer"`
	CollectionBoost    int         `json:"collection_boost"`
	Street             struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"street"`
	City struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"city"`
	TeleportPrice           int       `json:"teleport_price"`
	BuyerTransactionFee     float64   `json:"buyer_transaction_fee"`
	SellerTransactionFee    float64   `json:"seller_transaction_fee"`
	IsConstructionForbidden bool      `json:"is_construction_forbidden"`
	OwnershipChangedAt      time.Time `json:"ownership_changed_at"`
	Building                struct {
		ConstructionStatus string      `json:"constructionStatus"`
		Construction       interface{} `json:"construction"`
		Details            struct {
			MaxStackedSparks int     `json:"maxStackedSparks"`
			MinStackedSparks float64 `json:"minStackedSparks"`
			StepSparks       float64 `json:"stepSparks"`
		} `json:"details"`
		PropModelID int   `json:"propModelID"`
		PropertyID  int64 `json:"propertyID"`
		NftID       int   `json:"nftID"`
		IsLocked    bool  `json:"isLocked"`
	} `json:"building"`
	IsOwnerInJail bool `json:"is_owner_in_jail"`
	State         struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"state"`
}