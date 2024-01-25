package api

type RequestType int

const (
	Place RequestType = iota
	Pass
	Resign
)

var requestNames = map[RequestType]string{
	Place:  "Place",  // row, column, color as data
	Pass:   "Pass",   // no data
	Resign: "Resign", // no data
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

// Marshal a request to a string
func (r *Request) Marshal() string {
	return r.Data
}

// Unmarshal a string to a request
func (r *Request) Unmarshal(data []byte) {
	r.Data = string(data)
}

type Request struct {
	Type RequestType
	Data string
}
