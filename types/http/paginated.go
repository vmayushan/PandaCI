package typesHTTP

type Paginated struct {
	Next bool `json:"next"`
	Data any  `json:"data"`
}
