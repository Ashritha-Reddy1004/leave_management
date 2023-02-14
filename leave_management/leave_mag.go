package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Student struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Email string `json:"Email"`
}
type StudentCred struct {
	Name     string `json:"name"`
	Id       string `json:"id"`
	Password string `json:"password"`
}
type AdminCred struct {
	Password string `json:"password"`
	Id       string `json:"id"`
}
type LeaveReq struct {
	Name   string `json:"name"`
	Id     string `json:"leaveid"`
	Reason string `json:"reason"`
	Date   string `json:"date"`
	Status string `json:"status"`
}
type LeaveApproval struct {
	Id     string `json:"leaveid"`
	Status string `json:"status"`
}
type Signin struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
type Claims struct {
	Id string
	jwt.StandardClaims
}
type JsonSigninRes struct {
	Status  string `json:"status"`
	Token   string `json:"token"`
	Invalid bool   `json:"invalid"`
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
var mongo_uri context.Context
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
	if err != nil {
		log.Fatal(err)
	}
	return os.Getenv(key)
}

var students_collection *mongo.Collection
var leave_collection *mongo.Collection

func init() {
	mongo_uri = context.Background()
	mongo_uri := goDotEnvVariable("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Fatal(err)
	}
	//err = client.Ping(mongo_uri, nil)
	fmt.Println("Connection Established")
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	studentsdb = client.Database(leavemanagement).Collection(studentCollection)
	admindb = client.Database(leavemanagement).Collection(adminCollection)
	leaverequestdb = client.Database(leavemanagement).Collection(leaverequestCollection)
	leaveapprovaldb = client.Database(leavemanagement).Collection(leaverequestCollection)
	userdetailsdb = client.Database(leavemanagement).Collection(userdetailsCollection)

}
func PrintMessage(message string) {
	fmt.Println("---------------------------------------")
	fmt.Println(message)
	fmt.Println("----------------------------------------")
}
func CreateJWT(Id string) (response string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Id: Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err == nil {
		return tokenString, nil
	}
	return "", err
}
func VerifyToken(tokenString string) (id string, err error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return string(jwtSecretKey), nil
	})
	if token != nil {
		return claims.Id, nil
	}
	return "", err
}
func StudentLogin(w http.ResponseWriter, r *http.Request) {
	var Login Signin
	var result Student
	json.NewDecoder(r.Body).Decode(&Login)
	if Login.Id == "" {
		json.NewEncoder(w).Encode(ErrorRes{
			Status:  "400",
			Message: "Id cannot be alphabets",
		})
	} else if Login.Password == "" {
		json.NewEncoder(w).Encode(ErrorRes{
			Status:  "400",
			Message: "Can't add the details with null",
		})
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		hashpassword := Login.Password
		h := sha256.New()
		h.Write([]byte(hashpassword))
		Login.Password = hex.EncodeToString(h.Sum(nil))
		var err = userdetailsdb.FindOne(ctx, bson.M{
			"id":       Login.Id,
			"password": Login.Password,
		}).Decode(&result)
		defer cancel()
		if err != nil {
			json.NewEncoder(w).Encode(ErrorRes{
				Status:  "400",
				Message: fmt.Sprintf("Cannot add the student details with null values err= %s", err),
			})
		} else {
			tokenString, _ := CreateJWT(Login.Id)
			if tokenString == "" {
				json.NewEncoder(w).Encode(ErrorRes{
					Status:  "400",
					Message: "Cannot add the student data with null values",
				})
			}
			var ErrorMsg = ErrorMsg{
				Status:  string(http.StatusOK),
				Message: "You are already a user, try signing in",
				Response: JsonSigninRes{
					Status:  "200",
					Token:   tokenString,
					Invalid: false,
					Message: fmt.Sprintf("Successful login %s", Login.Id),
				},
			}
			successJsonResponse, jsonError := json.Marshal(ErrorMsg)
			if jsonError != nil {
				json.NewEncoder(w).Encode(ErrorRes{
					Status:  "400",
					Message: "Cannot add the student details with null values",
				})
			}
			w.Header().Set("Contents", "application/json")
			w.Write(successJsonResponse)
		}
	}
}
func AdminLogin(w http.ResponseWriter, r *http.Request) {
	var LoginRequest AdminCred
	var result Student
	json.NewDecoder(r.Body).Decode(&LoginRequest)
	if LoginRequest.Id == "" {
		json.NewEncoder(w).Encode(ErrorMsg{
			Status:  "400",
			Message: "Id cannot be null",
		})
	} else if LoginRequest.Password == "" {
		json.NewEncoder(w).Encode(ErrorMsg{
			Status:  "400",
			Message: "Cannot insert null details",
		})
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		hashpassword := LoginRequest.Password
		h := sha256.New()
		h.Write([]byte(hashpassword))
		LoginRequest.Password = hex.EncodeToString(h.Sum(nil))
		var err = userdetailsdb.FindOne(ctx, bson.M{
			"id":       LoginRequest.Id,
			"password": LoginRequest.Password,
		}).Decode(&result)
		defer cancel()
		if err != nil {
			json.NewEncoder(w).Encode(ErrorRes{
				Status:  "400",
				Message: fmt.Sprintf("Cannot add null student details. err=%s", err),
			})
		} else {
			tokenString, _ := CreateJWT(LoginRequest.Id)
			if tokenString == "" {
				json.NewEncoder(w).Encode(ErrorMsg{
					Status:  "400",
					Message: "Cannot add null student details",
				})
			}
			var ErrorMsg = ErrorMsg{
				Status:  string(http.StatusOK),
				Message: "Credentials in use, try signing in",
				Response: JsonSigninRes{
					Status:  "400",
					Token:   tokenString,
					Invalid: false,
					Message: fmt.Sprintf("Signin Successful %s", LoginRequest.Id),
				},
			}

			successJsonResponse, jsonError := json.Marshal(ErrorMsg)
			if jsonError != nil {
				json.NewEncoder(w).Encode(ErrorRes{
					Status:  "400",
					Message: "Cannot insert null student data",
				})
			}
			w.Header().Set("Contents", "application/json")
			w.Write(successJsonResponse)
		}
	}

}
func SetAdminCred(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contents", "application/json")
	var admin AdminCred
	json.NewDecoder(r.Body).Decode(&admin)
	fmt.Println("admin ", admin)
	if admin.Id == "" || admin.Password == "" {
		json.NewEncoder(w).Encode(JsonResStudentCred{
			Status:  "400",
			Message: "Cannot add null student details",
		})

	}
	hashpassword := admin.Password
	h := sha256.New()
	h.Write([]byte(hashpassword))
	admin.Password = hex.EncodeToString(h.Sum(nil))
	result, err := userdetailsdb.InsertOne(mongo_uri, admin)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResAdminCred{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
	}
	json.NewEncoder(w).Encode(JsonResAdminCred{
		Status:  "200",
		Data:    []AdminCred{admin},
		Message: fmt.Sprintf("Admin Inserted Successfully : %s", result.InsertedID),
	})
}
func SetStudentCred(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contents", "application/json")
	var student StudentCred
	json.NewDecoder(r.Body).Decode(&student)
	fmt.Println("student", student)
	if student.Id == "" || student.Password == "" {
		json.NewEncoder(w).Encode(JsonResStudentCred{
			Status:  "400",
			Message: "Cannot add null student details",
		})

	}
	hashpassword := student.Password
	h := sha256.New()
	h.Write([]byte(hashpassword))
	student.Password = hex.EncodeToString(h.Sum(nil))
	result, err := userdetailsdb.InsertOne(mongo_uri, student)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResStudentCred{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error :%v", err),
		})
	}
	json.NewEncoder(w).Encode(JsonResStudentCred{
		Status:  "200",
		Data:    []StudentCred{student},
		Message: fmt.Sprintf("Student added successfully :%s", result.InsertedID),
	})
}

func AddStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	json.NewDecoder(r.Body).Decode(&student)
	fmt.Println("student", student)
	if student.Id == "" || student.Name == "" || student.Email == "" {
		json.NewEncoder(w).Encode(JsonResStudent{
			Status:  "400",
			Message: "Cannot add null student details ",
		})
	}
	result, err := studentsdb.InsertOne(mongo_uri, student)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResStudent{
			Status:  "200",
			Data:    []Student{student},
			Message: fmt.Sprintf("Student added successfully: %s", result.InsertedID),
		})
	}
}
func AddLeaveRequest(w http.ResponseWriter, r *http.Request) {
	var student LeaveReq
	json.NewDecoder(r.Body).Decode(&student)
	fmt.Println("Leave Request", student)
	if student.Id == "" || student.Name == "" || student.Reason == "" || student.Date == "" {
		json.NewEncoder(w).Encode(JsonResStudent{
			Status:  "400",
			Message: "Cannot add null student details ",
		})
	}
	student.Status = "Pending"
	result, err := leaverequestdb.InsertOne(mongo_uri, student)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResLeaveReq{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
	}
	w.Header().Set("Contents", "application/json")
	json.NewEncoder(w).Encode(JsonResLeaveReq{
		Status:  "200",
		Data:    []LeaveReq{student},
		Message: fmt.Sprintf("Leave request added successfully :%s", result.InsertedID),
	})
}

func AddApprovedLeaves(w http.ResponseWriter, r *http.Request) {
	var approved LeaveApproval
	var status LeaveReq
	json.NewDecoder(r.Body).Decode(&approved)
	fmt.Println("Approved Students", approved)
	if approved.Id == "" || approved.Status == "" {
		json.NewEncoder(w).Encode(JsonResStudent{
			Status:  "400",
			Message: "Cannot add null student details",
		})
	}
	status.Status = "Accepted"
	filter := bson.M{
		"$set": bson.M{
			"status": status.Status,
		},
	}
	query := bson.M{
		"StudentId": approved.Id,
	}
	result, err := leaveapprovaldb.InsertOne(mongo_uri, approved)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResLeaveApproval{
			Status:  "400",
			Message: fmt.Sprintf("Internal error: %v", err),
		})
	}
	_ = leaverequestdb.FindOneAndUpdate(mongo_uri, query, filter)
	w.Header().Set("Contents", "application/json")
	json.NewEncoder(w).Encode(JsonResLeaveApproval{
		Status:  "200",
		Data:    []LeaveApproval{approved},
		Message: fmt.Sprintf("Leave Approved Succesfully : %s", result.InsertedID),
	})
}
func GetAllStudents(w http.ResponseWriter, r *http.Request) {
	var students []Student
	cursor, err := studentsdb.Find(context.Background(), bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(JsonResStudent{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
		return
	}
	err = cursor.All(context.Background(), &students)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResStudent{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
		return
	}
	res := JsonResStudent{
		Status:  "200",
		Data:    students,
		Message: "Leave requests listed successfully",
	}
	defer cursor.Close(context.Background())
	w.Header().Set("Contents", "application.json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}
func GetAllApprovedLeaves(w http.ResponseWriter, r *http.Request) {
	var approves []LeaveApproval
	cursor, err := leaveapprovaldb.Find(context.Background(), bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(JsonResLeaveApproval{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
		return
	}
	err = cursor.All(context.Background(), &approves)
	if err != nil {
		json.NewEncoder(w).Encode(JsonResLeaveApproval{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
		return
	}
	res := JsonResLeaveApproval{
		Status:  "200",
		Data:    approves,
		Message: "Leave approved Succesfully",
	}
	defer cursor.Close(context.Background())
	w.Header().Set("Contents", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)

}
func GetLeaveRequest(w http.ResponseWriter, r *http.Request) {
	var leaves []LeaveReq
	cursor, err := leaverequestdb.Find(context.Background(), bson.M{})
	if err != nil {
		json.NewEncoder(w).Encode(JsonResLeaveReq{
			Status:  "400",
			Message: fmt.Sprintf("Internal Error : %v", err),
		})
		return
	}
	res := JsonResLeaveReq{
		Status:  "200",
		Data:    leaves,
		Message: "Leave requests listed",
	}
	defer cursor.Close(context.Background())
	w.Header().Set("contents", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&res)
}
func main() {
	a := mux.NewRouter()
	a.HandleFunc("/AddStudent", AddStudent).Methods("POST")
	a.HandleFunc("SetStudentCreden", SetStudentCred).Methods("POST")
	a.HandleFunc("/AddLeaveRequest", AddLeaveRequest).Methods("POST")
	a.HandleFunc("/GetLeaveRequest", GetLeaveRequest).Methods("GET")
	a.HandleFunc("/GetAllStudents", GetAllStudents).Methods("GET")
	a.HandleFunc("/AddLeaveApproval", AddApprovedLeaves).Methods("POST")
	a.HandleFunc("/GetAllApprovedLeaves", GetAllApprovedLeaves).Methods("GET")
	a.HandleFunc("/StudentLogin", StudentLogin).Methods("POST")
	a.HandleFunc("/SetAdminCredentials", SetAdminCred).Methods("POST")
	a.HandleFunc("/SetStudentCred", SetStudentCred).Methods("POST")
	a.HandleFunc("/AdminLogin", AdminLogin).Methods("POST")
	fmt.Println("Attempted to start the server")
	log.Fatal(http.ListenAndServe(":8008", a))

}
