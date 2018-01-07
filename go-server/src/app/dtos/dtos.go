package dtos

// Position describes geographical position of the entity
type Position struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

type Vehicle struct {
	Position *Position `json:"position"`
	Id       string    `json:"id"`
	RouteId  string    `json:"routeId"`
}
