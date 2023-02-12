package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

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
	LeaveId int32  `json:"leaveid"`
	Status  string `json:"status"`
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	handleError(err)
	return os.Getenv(key)
}
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
