package game

type ClientRequest struct {
	Type       string  `json:"type"`
	X          float32 `json:"x"`
	Y          float32 `json:"y"`
	UUID       int     `json:"uuid"`
	TargetUUID int     `json:"target_uuid"`
}
