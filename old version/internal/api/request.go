package api

import "encoding/json"

type RequestType int

const (
	FindGame RequestType = iota
	Update
	Place
	Pass
	Resign
)

var requestNames = map[RequestType]string{
	FindGame: "FindGame", // json: {"size": 9, "player": "Robin", "human_opponent": true}
	Update:   "Update",   // no data
	Place:    "Place",    // row, column, color as data
	Pass:     "Pass",     // no data
	Resign:   "Resign",   // no data
}

func (c RequestType) String() string {
	return requestNames[c]
}

func StringToRequestType(s string) (RequestType, bool) {
	for k, v := range requestNames {
		if v == s {
			return k, true
		}
	}
	return 0, false // or return an error
}

type InitRequest struct {
	Size          int
	Player        string
	HumanOpponent bool
}

func (r *InitRequest) UnmarshalJSON(data []byte) error {
	type Alias InitRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	return json.Unmarshal(data, &aux)
}

type Request struct {
	Type RequestType
	Data string
}
