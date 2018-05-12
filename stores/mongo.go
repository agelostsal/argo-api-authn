package stores

import (
	"time"

	"github.com/ARGOeu/argo-api-authn/argo-consts"
	"github.com/ARGOeu/argo-api-authn/utils"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStore struct {
	Server   string
	Database string
	Session  *mgo.Session
}

// Initialize initializes the mongo stores struct
func (mongo *MongoStore) SetUp() {
	session, err := mgo.Dial(mongo.Server)
	if err != nil {
		log.Fatal("STORE", "\t", err.Error())
	}

	mongo.Session = session

	log.Info("STORE", "\t", "Connected to Mongo: ", mongo.Server)

}

func (mongo *MongoStore) Close() {
	mongo.Session.Close()
}

func (mongo *MongoStore) QueryServices(name string) ([]QService, error) {

	var qServices []QService
	var err error

	c := mongo.Session.DB(mongo.Database).C("services")
	query := bson.M{}

	if name != "" {
		query = bson.M{"name": name}
	}

	err = c.Find(query).All(&qServices)

	if err != nil {
		log.Fatal("STORE", "\t", err.Error())
		err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
		return []QService{}, err
	}

	return qServices, err
}

func (mongo *MongoStore) QueryAuthMethod(service string, host string, typeName string) (map[string]interface{}, error) {

	var qAuthType map[string]interface{}
	var err error

	c := mongo.Session.DB(mongo.Database).C("auth_types")
	err = c.Find(bson.M{"type": typeName, "service": service, "host": host}).One(&qAuthType)

	if err != nil {
		log.Fatal("STORE", "\t", err.Error())
		err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
		return qAuthType, err
	}

	return qAuthType, err
}

func (mongo *MongoStore) QueryBindingsByDN(dn string, host string) ([]QBinding, error) {

	var qBindings []QBinding
	var err error

	c := mongo.Session.DB(mongo.Database).C("bindings")
	err = c.Find(bson.M{"dn": dn, "host": host}).All(&qBindings)

	if err != nil {
		log.Fatal("STORE", "\t", err.Error())
		err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
		return []QBinding{}, err
	}

	return qBindings, err
}

func (mongo *MongoStore) QueryBindings(service string, host string) ([]QBinding, error) {

	var qbindings []QBinding
	var err error
	query := bson.M{}

	db := mongo.Session.DB(mongo.Database)
	c := db.C("bindings")

	if service != "" && host != "" {
		query = bson.M{"service": service, "host": host}
	}

	if err = c.Find(query).All(&qbindings); err != nil {
		err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
		return qbindings, err
	}
	return qbindings, err
}

//InsertBinding inserts a new binding into the datastore
func (mongo *MongoStore) InsertBinding(name string, service string, host string, dn string, oidcToken string, uniqueKey string) (QBinding, error) {

	var qBinding QBinding
	var err error

	qBinding = QBinding{Name: name, Service: service, Host: host, DN: dn, OIDCToken: oidcToken, UniqueKey: uniqueKey, CreatedOn: time.Now().Format(argo_consts.ZULU_FORM)}
	db := mongo.Session.DB(mongo.Database)
	c := db.C("bindings")

	if err := c.Insert(qBinding); err != nil {
		log.Fatal("STORE", "\t", err.Error())
		err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
		return QBinding{}, nil
	}

	return qBinding, err
}

//UpdateBinding updates the given binding
func (mongo *MongoStore) UpdateBinding(original QBinding, updated QBinding) (QBinding, error) {

	var err error

	db := mongo.Session.DB(mongo.Database)
	c := db.C("bindings")

	if err := c.Update(original, updated); err != nil {
		log.Fatal("STORE", "\t", err.Error())
		err = &utils.APIError{Status: "INTERNAL SERVER ERROR", Code: 500, Message: err.Error()}
		return QBinding{}, err
	}

	return updated, err
}
