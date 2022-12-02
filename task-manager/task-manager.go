package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"net/url"
	"os"

	"log"
	model "task-manager/models"
)

type ResponseSourceModel struct {
	Status  string `json:"status" default:"true"`
	Data    any    `json:"data" default:"[]"`
	Message string `json:"message" default:""`
	Error   string `json:"error" default:""`
}

var db *gorm.DB

func main() {
	var err error

	log.Println("Start app")
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSL_MODE")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}

	err = createSchema()
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	s := r.PathPrefix("/tasks/").Subrouter()
	s.HandleFunc("/index", getAllTasks).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func createSchema() error {
	err := db.AutoMigrate(&model.Task{})
	if err != nil {
		return err
	}
	return nil
}

func stateActive() {
	u := url.URL{Scheme: "wss", Host: "localhost:8080", Path: "/task-manager-active"}
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Printf("handshake failed with status %d", resp.StatusCode)
		log.Fatal("dial:", err)
	}
	//When the program closes close the connection
	defer c.Close()

}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	//specify status code
	w.WriteHeader(http.StatusOK)

	//update response writer
	fmt.Fprintf(w, "API is up and running")
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []model.Task

	err := db.Model(&model.Task{}).Preload(clause.Associations, model.PreloadTasks).Where("parent_id is null").Find(&tasks).Error
	if err != nil {
		log.Fatalln(err)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	success(w, "true", tasks, "", "")
}

func success(w http.ResponseWriter, status string, data any, message string, error string) {
	var response = ResponseSourceModel{
		Status:  status,
		Data:    data,
		Message: message,
		Error:   error,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}

	w.Write(jsonResponse)
}
