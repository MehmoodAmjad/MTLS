package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func connected(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "Hello World!")
}
func fun2(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/2" {
		http.NotFound(w, req)
		return
	}
	fmt.Fprintf(w, "The second function has been called")
}
func main() {
	addr := flag.String("addr", ":4000", "HTTPS network address")
	certFile := flag.String("certfile", "../servercert/cert.pem", "certificate PEM file")
	keyFile := flag.String("keyfile", "../servercert/key.pem", "key PEM file")
	clientCertFile := flag.String("clientcert", "../clientcert/clientcert.pem", "certificate PEM for client authentication")
	flag.Parse()

	//mux := http.NewServeMux()
	//mux.HandleFunc("/", connected)
	//mux.HandleFunc("/2", fun2)
	// seeding
	courses = append(courses, Course{CourseID: "2", CourseName: "GOLANG", CoursePrice: 299, Author: &Author{Fullname: "Mehmood Amjad", Website: "securiti.go"}})
	courses = append(courses, Course{CourseID: "4", CourseName: "Docker", CoursePrice: 399, Author: &Author{Fullname: "Mehmood Amjad", Website: "foundri.go"}})

	//router.HandleFunc("/test", testVal).Methods("GET")
	// routing
	router := mux.NewRouter()

	router.HandleFunc("/", serveHome).Methods("GET")
	router.HandleFunc("/courses", getAllCourses).Methods("GET")
	router.HandleFunc("/course/{courseid}", getOneCourse).Methods("GET")
	router.HandleFunc("/course", createOneCourse).Methods("POST")
	router.HandleFunc("/course/{courseid}", updateOneCourse).Methods("PUT")
	router.HandleFunc("/course/{courseid}", deleteOneCourse).Methods("DELETE")

	// Trusted client certificate.
	clientCert, err := os.ReadFile(*clientCertFile)
	if err != nil {
		log.Fatal(err)
	}
	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCert)

	srv := &http.Server{
		Addr:    *addr,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
			ClientCAs:                clientCertPool,
			ClientAuth:               tls.RequireAndVerifyClientCert,
		},
	}

	log.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS(*certFile, *keyFile)
	log.Fatal(err)
}

// Model for course (goes in file)
type Course struct {
	CourseID    string  `json:"courseid"`
	CourseName  string  `json:"coursename"`
	CoursePrice int     `json:"-"`
	Author      *Author `json:"-"`
}

type Author struct {
	Fullname string `json:"fullname"`
	Website  string `json:"website"`
}

// fake DB
var courses []Course

// middleware, helper(goes in file)
func (c *Course) IsEmpty() bool {
	//return c.CourseID == "" && c.CourseName == ""
	return c.CourseName == ""
}

// conntrollers (goes in seperate files)

// func testVal(w http.ResponseWriter, r *http.Request){
// 	json.NewEncoder(w).Encode([]byte{})
// }

// serve home route
func serveHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Welcome to API</h1>"))
}

func getAllCourses(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get all Courses")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	w.Header().Set("Content=Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}

func getOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get One Course")
	w.Header().Set("Content=Type", "application/json")
	// grab id from request
	params := mux.Vars(r)
	// loop through courses and find matchingn id then return the reponse
	for _, course := range courses {
		if course.CourseID == params["courseid"] {
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	json.NewEncoder(w).Encode("no Course Founnd with given id")
	return
}

func createOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create onne course")
	w.Header().Set("Content=Type", "application/json")

	// what if body is empty
	if r.Body == nil {
		json.NewEncoder(w).Encode("Please send some data")
	}

	// what about data (being sent in form of {})
	var course Course

	_ = json.NewDecoder(r.Body).Decode(&course)
	if course.IsEmpty() {
		json.NewEncoder(w).Encode("No data inside JSON")
		return
	}

	// gennerate unnique id, string
	// append course into courses

	rand.Seed(time.Now().UnixNano())
	course.CourseID = strconv.Itoa(rand.Intn(100))
	courses = append(courses, course)
	json.NewEncoder(w).Encode(course)
	return
}

func updateOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update one course")
	w.Header().Set("Content=Type", "application/json")

	// first - grab id fro request
	params := mux.Vars(r)

	// loop through value to get id then remove then add with ID

	for index, course := range courses {
		if course.CourseID == params["courseid"] {
			courses = append(courses[:index], courses[index+1:]...)
			var course Course
			_ = json.NewDecoder(r.Body).Decode(&course)
			course.CourseID = params["courseID"]
			courses = append(courses, course)
			json.NewEncoder(w).Encode(course)
			return
		}
	}
	// send a response whenn id not found
}

func deleteOneCourse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete one course")
	w.Header().Set("Content=Type", "application/json")

	params := mux.Vars(r)
	// loop, finnd id, remove(index,index+1)
	for index, course := range courses {
		if course.CourseID == params["courseid"] {
			courses = append(courses[:index], courses[index+1:]...)
			break
		}
	}
}
