package main

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"

	. "devgo/GoLanG/config"
	. "devgo/GoLanG/dao"
	. "devgo/GoLanG/models"

	"github.com/gorilla/handlers"
)

var config = Config{}
var dao = CategoriesDAO{}

// GET list of categories
func AllCategoriesEndPoint(w http.ResponseWriter, r *http.Request) {
	categories, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, categories)
}

// GET a category by its ID
func FindCategoryEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	category, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}
	respondWithJson(w, http.StatusOK, category)
}

type Payload struct {
	Slot_temp string
	Data      string
	Time      string
	Device    string
	Signal    string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

// POST a new Category
func CreateCategoryEndPoint(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	defer r.Body.Close()
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	category.ID = bson.NewObjectId()
	if err := dao.Insert(category); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, category)
}

// PUT update an existing category
func UpdateCategoryEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(category); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing movie
func DeleteCategoryEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(category); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/categories", AllCategoriesEndPoint).Methods("GET")
	r.HandleFunc("/categories", CreateCategoryEndPoint).Methods("POST")
	r.HandleFunc("/categories", UpdateCategoryEndPoint).Methods("PUT")
	r.HandleFunc("/categories", DeleteCategoryEndPoint).Methods("DELETE")
	r.HandleFunc("/categories/{id}", FindCategoryEndpoint).Methods("GET")

	http.ListenAndServe(":9000", handlers.CORS()(r))
}
