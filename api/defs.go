package api

// InfoResponse is used as definition how info response should look like
type InfoResponse struct {
	PID             string
	ConfigFile      string
	RecipesInQueue  string
	MemoryAllocated string
}

// WorkersResponse is used as definition how workers response should look like
type WorkersResponse struct {
	Message string
	List [][]string `json:"list,omitempty"`
}