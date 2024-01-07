package tasks

import "time"

type Object struct {
	Oid           string  `json:"oid,omitempty"`
	Size          int64   `json:"size,omitempty"`
	Authenticated bool    `json:"authenticated,omitempty"`
	Actions       Actions `json:"actions,omitempty"`
}

type BatchReq struct {
	Operation string   `json:"operation"`
	Objects   []Object `json:"objects"`
	Transfers []string `json:"transfers"`
	// Ref       struct {
	// 	Name string `json:"name"`
	// } `json:"ref"`
	HashAlgo string `json:"hash_algo"`
}

type BatchResp struct {
	Transfer string   `json:"transfer,omitempty"`
	Objects  []Object `json:"objects,omitempty"`
	HashAlgo string   `json:"hash_algo,omitempty"`
}

type Download struct {
	Href      string            `json:"href,omitempty"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresAt time.Time         `json:"expires_at,omitempty"`
}

type Actions struct {
	Download Download `json:"download,omitempty"`
}
