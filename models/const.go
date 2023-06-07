package models

const (
	Success = iota
	Failed
)

var StatusMap = map[int]string{
	Success: "Success",
	Failed:  "Failed",
}