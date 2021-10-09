package all_user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	Id    string `json: "id"`
	Name  string `json: "name"`
	Email string `json: "email"`
	Pwd   string `json: "pwd"`
}

type usersHandler struct {
	sync.Mutex
	store map[string]User
}

func (h *usersHandler) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader((http.StatusMethodNotAllowed))
		w.Write([]byte("Method not allowed"))
		return
	}
}

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

func insertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	// select database and collection ith Client.Database method
	// and Database.Collection method
	collection := client.Database(dataBase).Collection(col)

	// InsertOne accept two argument of type Context
	// and of empty interface
	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

func query(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	// select database and collection.
	collection := client.Database(dataBase).Collection(col)
	result, err = collection.Find(ctx, query, options.Find().SetProjection(field))
	return
}

func ping(client *mongo.Client, ctx context.Context) error {

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	fmt.Println("connected successfully")
	return nil

}

func (h *usersHandler) get(w http.ResponseWriter, r *http.Request) {
	users := make([]User, len(h.store))

	h.Lock()
	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		fmt.Println("c")
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
		fmt.Println("d")
		panic(err2)
	}
	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		// handle the error
		fmt.Println("e")
		panic(err)
	}

	// printing the result of query.
	fmt.Println("Query Reult/n")
	for _, doc := range results {
		fmt.Println("/n")
		fmt.Println(doc)
	}
	h.Unlock()
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		fmt.Println("f")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

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
		fmt.Println("c")
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
		fmt.Println("d")
		panic(err2)
	}
	var results []bson.D

	if err := cursor.All(ctx, &results); err != nil {

		fmt.Println("e")
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
		fmt.Println("g")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *usersHandler) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("h")
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
		fmt.Println("i")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	user.Id = fmt.Sprintf("%d", time.Now().UnixNano())
	fmt.Println(user.Id)

	h.Lock()
	h.store[user.Id] = user

	client, ctx, cancel, err1 := connect("mongodb+srv://srishti:XqH8QC4le5e4YvYQ@all-user-data.5slh5.mongodb.net/instagram-users?retryWrites=true&w=majority")
	if err1 != nil {
		fmt.Println("j")
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

func newUsersHandler() *usersHandler {
	return &usersHandler{
		store: map[string]User{
			"id1": User{
				Id:    "1234",
				Name:  "TestName",
				Email: "test@gmail.com",
				Pwd:   "test",
			},
		},
	}
}
