package main

//IMPORT PACKAGES
import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//DATA STRUCTURE TO CREATE A USER
type User struct {
	Id    string `json: "id"`
	Name  string `json: "name"`
	Email string `json: "email"`
	Pwd   string `json: "pwd"`
}

//DATA STRUCTURE FOR POSTS FROM A USER
type UserPost struct {
	Id        string `json: "id"`
	Caption   string `json: "caption"`
	ImageURL  string `json: "imageurl"`
	TimeStamp int64  `json: "timestamp"`
}

//STORING USER DATA:
type usersHandler struct {
	sync.Mutex
	store map[string]User
}

//STORING POSTS DATA:
type userpostsHandler struct {
	sync.Mutex
	store map[string]User
}

//FUNCTION TO CLOSE CONNECTION
func close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			fmt.Println("b")
			panic(err)
		}
	}()
}

//FUNCTION TO CONNECT TO MONGODB
func connect(uri string) (*mongo.Client, context.Context,
	context.CancelFunc, error) {

	// ctx will be used to set deadline for process, here
	// deadline will of 30 seconds.
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)

	// mongo.Connect return mongo.Client method
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

//FUNCTION TO GET INSERT ONE PARTICULAR DOCUMENT TO MONGODB
func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

//FUNCTION TO GET ONE PARTICULAR DOCUMENT FROM MONGODB
func query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

//FUNCTION TO PING MONGODB TO CHECK CONNECTION STATUS
func ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil

}

//FUNCTION FOR SWITCH CASE OF USERS ENDPOINTS
func (h *usersHandler) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.postUser(w, r)
		return
	default:
		w.WriteHeader((http.StatusMethodNotAllowed))
		w.Write([]byte("Method not allowed"))
		return
	}
}

//FUNCTION FOR SWITCH CASE OF POSTS ENDPOINTS
func (h *userpostsHandler) userpost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getPosts(w, r)
		return
	case "POST":
		h.postPost(w, r)
		return
	default:
		w.WriteHeader((http.StatusMethodNotAllowed))
		w.Write([]byte("Method not allowed"))
		return
	}
}

//POST REQUEST WHICH ENABLES TO GET ALL USERS DETAILS ID VIA GET REQUEST
func (h *usersHandler) get(w http.ResponseWriter, r *http.Request) {
	users := make([]User, len(h.store))

	h.Lock()
	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Mongo connecting")
	//fmt.Printf("%T",client)
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	var filter, option interface{}

	filter = bson.D{}

	option = bson.D{}

	cursor, err2 := query(client, ctx, "instagram-users", "user-profile", filter, option)

	if err2 != nil {
		panic(err2)
	}
	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Reult/n")
	for _, doc := range results {
		fmt.Println(doc)
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//POST REQUEST WHICH ENABLES TO GET ALL POSTS OF ALL USER USING ID VIA GET REQUEST
//WHICH WILL BE PAGIANATED
func (h *userpostsHandler) getallpost(w http.ResponseWriter, r *http.Request) {
	users := make([]User, len(h.store))

	h.Lock()
	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Mongo connecting")
	//fmt.Printf("%T",client)
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	var filter, option interface{}

	filter = bson.D{}

	option = bson.D{}

	cursor, err2 := query(client, ctx, "instagram-users", "user-posts", filter, option)

	if err2 != nil {
		panic(err2)
	}
	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Reult/n")
	for _, doc := range results {
		fmt.Println(doc)
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//POST REQUEST WHICH ENABLES TO GET USER DETAILS USING ID VIA GET REQUEST
func (h *usersHandler) getUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		fmt.Println("a")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Lock()
	user, ok := h.store[parts[2]]
	id_to_get := parts[2]
	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Mongo connecting")

	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	var filter, option interface{}
	filter = bson.D{{"id", id_to_get}}
	option = bson.D{}

	cursor, err2 := query(client, ctx, "instagram-users", "user-profile", filter, option)

	if err2 != nil {
		panic(err2)
	}
	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//POST REQUEST WHICH ENABLES USERS TO POST AN IMAGE VIA POST REQUEST AND APPENDS IT TO MONGODB
func (h *userpostsHandler) postPost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", ct)))
		return
	}
	var userpost UserPost
	err = json.Unmarshal(bodyBytes, &userpost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	h.Lock()

	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Mongo connecting")
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)
	userpost.TimeStamp = time.Now().UnixNano()
	var document interface{}
	document = userpost
	insertOneResult, err := insertOne(client, ctx, "instagram-users", "user-posts", document)

	fmt.Println("Result of InsertOne")
	fmt.Println(insertOneResult.InsertedID)

	defer h.Unlock()
}

//POST REQUEST WHICH ENABLES TO GET ALL POSTS OF A USER USING ID VIA GET REQUEST
func (h *userpostsHandler) getPosts(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		fmt.Println("a")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.Lock()
	user, ok := h.store[parts[2]]
	id_to_get := parts[2]
	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Mongo connecting")

	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)

	var filter, option interface{}
	filter = bson.D{{"id", id_to_get}}
	option = bson.D{}

	cursor, err2 := query(client, ctx, "instagram-users", "user-posts", filter, option)

	if err2 != nil {
		panic(err2)
	}
	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		panic(err)
	}

	fmt.Println("Query Result")
	for _, doc := range results {
		fmt.Println(doc)
	}
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

//POST REQUEST WHICH ENABLES USERS TO CREATE A NEW ACCOUNT/USER VIA POST REQUEST AND APPENDS IT TO MONGODB
func (h *usersHandler) postUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("Need content-type 'application/json', but got '%s'", ct)))
		return
	}
	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}

	h.Lock()
	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Mongo connecting")
	//fmt.Printf("%T",client)
	defer close(client, ctx, cancel)

	// Ping mongoDB with Ping method
	ping(client, ctx)
	var document interface{}

	sum := sha256.Sum256([]byte(user.Pwd))
	user.Pwd = hex.EncodeToString(sum[:])
	document = user
	insertOneResult, err := insertOne(client, ctx, "instagram-users", "user-profile", document)

	fmt.Println("Result of InsertOne")
	fmt.Println(insertOneResult.InsertedID)

	defer h.Unlock()
}

//FUNCTION TO CREATER POINTER FOR USERSHANDLER
func newUsersHandler() *usersHandler {
	return &usersHandler{}
}

//FUNCTION TO CREATER POINTER FOR USERPOSTSHANDLER
func newUserpostsHandler() *userpostsHandler {
	return &userpostsHandler{}
}

//MAIN FUNCTION WHICH RUNS ON STARTING SERVER
func main() {
	fmt.Println("working lol")

	usersHandler := newUsersHandler()
	userpostsHandler := newUserpostsHandler()
	http.HandleFunc("/users", usersHandler.users)
	http.HandleFunc("/userpost", userpostsHandler.userpost)
	http.HandleFunc("/users/", usersHandler.getUser)
	http.HandleFunc("/userpost/", userpostsHandler.getPosts)
	http.HandleFunc("/posts", userpostsHandler.getallpost)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}

