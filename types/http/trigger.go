package typesHTTP

type TriggerRunRequest struct {
	SHA    string `json:"sha"`
	Branch string `json:"branch"`
}
