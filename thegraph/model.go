package thegraph

type Query struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}
