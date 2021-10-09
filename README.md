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
* URL: "users/_id_"

## Added functionalities:

1. Passwords have been securely stored such they can't be reverse engineered (using sha256)
2. Only JSON format is used

## Attributes:

1. User:
* Unique ID
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
* Install mongo dependencies (run: go get go.mongodb.org/mongo-driver/mongo)
  
### Steps:

* Clone this repository 
* cd to the folder "cmd" in the terminal
* Run the command "go run server.go"; This should start the server
* Open another terminal in this folder (cmd)
* In this other terminal, run the commands for the specific functionality required: (make sure the path to curl is set under environment variables)
  * To create a new user, run :  curl localhost:3000/users -X POST -d '{"id":"<user id>","name":"<name>","email":"<email>","pwd":"<password>"}' -H "Content-Type: application/json"
  * To get details of all the users, run : curl localhost:3000/users
  * To get the details of a particular user, run : curl localhost:3000/users/_id_


