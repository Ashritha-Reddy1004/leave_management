package main

type student struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}
type admin struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
type leave struct {
	Name   string `json:"name"`
	Id     int32  `json:"id"`
	Reason string `json:"reason"`
	Date   string `json:"date"`
	Status string `json:"status"`
}
type leave_approval struct {
}
