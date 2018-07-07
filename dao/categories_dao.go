package dao

import (
	"log"

	. "devgo/GoLanG/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CategoriesDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "categories"
)

// Establish a connection to database
func (m *CategoriesDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of Categorys
func (m *CategoriesDAO) FindAll() ([]Category, error) {
	var Category []Category
	err := db.C(COLLECTION).Find(bson.M{}).All(&Category)
	return Category, err
}

// Find a movie by its id
func (m *CategoriesDAO) FindById(id string) (Category, error) {
	var Category Category
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&Category)
	return Category, err
}

// Insert a movie into database
func (m *CategoriesDAO) Insert(Category Category) error {
	err := db.C(COLLECTION).Insert(&Category)
	return err
}

// Delete an existing movie
func (m *CategoriesDAO) Delete(Category Category) error {
	err := db.C(COLLECTION).Remove(&Category)
	return err
}

// Update an existing movie
func (m *CategoriesDAO) Update(Category Category) error {
	err := db.C(COLLECTION).UpdateId(Category.ID, &Category)
	return err
}
