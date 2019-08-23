package api

// TableResponse is used as definition how table response should like
// this is later used to render table to output
type TableResponse struct {
	Header []string `json:"header"`
	Columns [][]string `json:"columns"`
}

// InfoResponse is used as definition how info response should look like
type InfoResponse struct {
	Status         string
	RecipesInQueue string
}

// WorkersResponse is used as definition how workers response should look like
type WorkersResponse struct {
	MessageResponse
	List [][]string `json:"list,omitempty"`
}

// MessageResponse hold just message response message
type MessageResponse struct {
	Message string
}
