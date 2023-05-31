package handlers

const (
	Success = iota
	Failed
)

var statusMap = map[int]string{
	Success: "Success",
	Failed:  "Failed",
}