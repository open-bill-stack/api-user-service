package structure

type Event struct {
	Type     string `json:"event_type"`
	UserUUID string `json:"user_uuid"`
}
