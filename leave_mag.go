package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	Name string `json:"name"`
	Id   int32  `json:"id"`
}
type StudentCred struct {
	Name     string `json:"name"`
	Id       int32  `json:"id"`
	Password string `json:"password"`
}
type AdminCred struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
type LeaveReq struct {
	Name    string `json:"name"`
	LeaveId int32  `json:"leaveid"`
	Reason  string `json:"reason"`
	Date    string `json:"date"`
	Status  string `json:"status"`
}
type LeaveApproval struct {
	LeaveId int32  `json:"leaveid"`
	Status  string `json:"status"`
}
type Signin struct {
	Id       int32  `json:"id"`
	Password string `json:"password"`
}
type Claims struct {
	Id int32
	jwt.StandardClaims
}
type JsonSigninRes struct {
	Status  string `json:"status"`
	Token   string `json:"token"`
	Invaild bool   `json:"invalid"`
	Message string `json:"message"`
}
type JsonResStudent struct {
	Status  string    `json:"status"`
	Data    []Student `json:"data"`
	Message string    `json:"message"`
}
type JsonResStudentCred struct {
	Status  string        `json:"status"`
	Data    []StudentCred `json:"data"`
	Message string        `json:"message"`
}
type JsonResAdminCred struct {
	Status  string      `json"status"`
	Data    []AdminCred `json:"data"`
	Message string      `json:"message"`
}
type JsonResLeaveReq struct {
	Status  string     `json:"status"`
	Data    []LeaveReq `json:"data"`
	Message string     `json:"message"`
}
type JsonResLeaveApproval struct {
	Status  string          `json:"status"`
	Data    []LeaveApproval `json:"data"`
	Message string          `json:"message"`
}
type ErrorRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type ErrorMsg struct {
	Status   string
	Message  string
	Response interface{}
}

var db *mongo.Client
var momgoCtx context.Context
var studentsdb *mongo.Collection
var leaverequestdb *mongo.Collection
var leaveapprovaldb *mongo.Collection
var userdetailsdb *mongo.Collection
var admindb *mongo.Collection

const studentCollection = "Student"
const leaverequestCollection = "Leaverequests"
const leaveapprovalCollection = "Leaveapprovals"
const userdetailsCollection = "Userdetails"
const leavemanagement = "Leavemanagement"
const adminCollection = "Admin"

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

func connectDB() {
	mongo_uri := goDotEnvVariable("MONGODB_URI")
	client, er := mongo.NewClient(options.Client().ApplyURI(mongo_url))
	handleError(err)
	fmt.Println("Connection Established")
	err = client.Connect(context.TODO())
	handleError(err)
	students_collection = client.Database("lms").Collection("Students")
	leave_collection = client.Database("lms").Collection("leaves")
}

func AddStudent(w http.ResponseWriter, r *http.Request) {
	a.Header().Set("Contents", "application/json")
	var student student
	json.NewDecoder(r.Body).Decode(&student)
	fmt.Println("student", student)

}
