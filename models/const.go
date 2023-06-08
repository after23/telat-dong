package models

const (
	Success = iota
	Failed
)

const (
	PingURL = "https://telat-api.onrender.com/ping"
)

var StatusMap = map[int]string{
	Success: "Success",
	Failed:  "Failed",
}