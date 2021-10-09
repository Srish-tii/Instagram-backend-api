# Instagram-backend-api
REST API of Instagram's functionalities; developed using GO (and Mongo).

## Constraints:

1. Complete API has been developed using Go 
2. MongoDB has been used to storage
3. Only standard libraries have been used

## Functionalities:

1. Create a User
* POST request
* JSON request body used
* Added data to the URL ‘/users'

2. Display details of all users
* GET request
* Display the user (unique) id, name, email and hash of the password
* URL: "/users"

3. Fetch details of the user using id
* GET request
* Displays user id, name, email and hash of the password
* URL: "/users/_id_"

4. Create a Post
* POST request
* JSON request body used
* Added data to the URL ‘/posts'

5. Display details of all posts (like scrolling through the feed)
* GET request
* Display the post (unique) id, caption, imageURL and timestamp
* URL: "/posts"

6. Fetch details of all the posts of a particular user using id
* GET request
* Display the post (unique) id, caption, imageURL and timestamp (for all posts of that user)
* URL: "/posts/_id_"

## Added functionalities:

1. Passwords have been securely stored such they can't be reverse engineered (using sha256)
2. Only JSON format is accepted (validation added)
3. Unique id is linking a user to his/her posts

## Attributes:

1. User:
* ID
* Name
* Email
* Password (saved in encrypted form; as a hash)

2. Post:
* Id
* Caption
* Image URL
* Posted Timestamp

## Directions to run the application:

### Prerequisites/ Software requirements:
  
* Go 
* Linux shell (eg. Gitbash)
* Install mongo dependencies (run `go get go.mongodb.org/mongo-driver/mongo`)
  
### Steps:

* Clone this repository 
* cd to the folder "cmd" in the terminal
* Run the command `go run server.go`; This should start the server
* Open another terminal in this folder (cmd)
* In this other terminal, run the commands for the specific functionality required: (make sure the path to curl is set under environment variables)
  * To create a new user, run  `curl localhost:3000/users -X POST -d '{"id":"<id>","name":"<name>","email":"<email>","pwd":"<password>"}' -H "Content-Type: application/json"`
  * To get details of all the users, run `curl localhost:3000/users`
  * To get the details of a particular user, run `curl localhost:3000/users/_id_`
  * To create a new post, run `curl localhost:3000/userpost -X POST -d '{"id":"<id>","caption":"<caption>","imageurl":"<url>"}' -H "Content-Type: application/json"`
  * To get details of all the posts on the feed, run `curl localhost:3000/posts`
  * To get all the posts of a particular user, run `curl localhost:3000/users/_id_`


