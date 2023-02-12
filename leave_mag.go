package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
	Id   int32 `json:"id"`
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
var jwtSecretKey = []byte("jwt_secret_key")

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
func PrintMessage(message string) {
	fmt.Println("---------------------------------------")
	fmt.Println(message)
	fmt.Println("----------------------------------------")
}
func CreateJWT(Id int32)(response string,err error){
	expirationTime := time.Now().Add(5*time.Minute)
	claims :=&Claims{
		Id : id,
		StaStandardClaims: jwt.StandardClaims{
			ExpiresAt : expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethod256,claims)
	tokeString,err :=token.SignedString(jwtSecretKey)
	if err==nil{
		return tokenString ,nil
	}
	return "",err
}
func VerifyToken(tokenString string)(id int32,err error){
	claims=&Claims{}
	token,err:=jwt.ParseWithClass(tokenString,claims,func(token *jwt.Token)(interface{},error){
		return string(jwtSecretKey),nil
})
if token !=nil{
	return claims.Id,nil
}
  return "",err
}
func StudentLogin(w http.ResponseWriter, r *http.Request){
	var LoginRequest Login 
	var result Student
	json.NewDecoder(r.Body).Decode(&LoginRequest)
	if LoginRequest.Id ==""{
		json.NewEncoder(w).Encode(ErrorResponse{
			Status :400,
			Message : "Id cannot be alphabets",
		})
	}else if LoginRequest.Password==""{
		json.NewEncoder(w).Encode(ErrorResponse{
			Status :400,
			Message:"Can't add the details with null",
		})
	}else{
		ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	    hashpassword :=LoginRequest.Password
		h:=sha256.New()
		h.Write([]byte(hashpassword))
		loginRequest.Password=hex.EncodeToString(h.Sum(nil))
		var err=userdetailsdb.FindONe(ctx,bson.M{
			"id":LoginRequest.Id,
			"password":LoginRequest.Password,
		}).Decode(&result)
		defer cancel()
		if err !=nil{
			json.NewEncoder(w).Encode(ErrorResponse{
				Status:400,
				Message:fmt.Sprintf("Cannot add the student details with null values err= ",err),
			})
		}else{
			tokenString,_:=CreateJWT(LoginRequest.Id)
			if tokenString ==""{
			json.NewEncoder(w).Encode(ErrorResponse{
				Status :400,
				Message :"Cannot add the student data with null values"
			})
		}
		var ErrorMsg = ErrorMsg{
			Status : http.StatusOK,
			Message :"You are already a user, try signing in",
			Response :JsonSigninRes{
				Status: 200,
				Token :tokenString,
				Invalid :false,
				Message : fmt.Sprintf("Successful login %s",&LoginRequest.Id)
			},
		}
	    ErrorMsg,jsonError :=json.Marshal(ErrorMsg)
		if jsonError !=nil{
			json.NewEncoder(w).Encode(EncodeResponse{
				Status :400,
				Message : "Cannot add the student details with null values",
			})
		}
		w.Header().Set("Content-Type","application/json")
		w.Write(successJsonResponse)
	}
}
}
func AdminLogin(w http.ResponseWriter, r *http.Request){
	var LoginRequest
}