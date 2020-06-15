# Overview

## Test General description
It is required to build an API Rest to affiliate clients, allow them to make a payment of their
membership and allow them to see your payments made. In total there are 4 use cases, which
described below:

### Case 1. Registration
Description
The API must allow a client to register with the following basic data: Name
full, phone number, email and password. These data should be stored those
data in the database.
Validations
● The telephone number (cell phone only) must have a correct format for Mexico
(E.g. +525123456789).
● The full name field must have more than one word, and each word must
have initial capital letters (E.g. Pedro Rodríguez Álvarez).
● The password must have at least one lowercase letter, one uppercase letter and one
number and contain at least 8 characters.
Answers
The service must reply with:
● 200 - if registration is successful
● 400 - if there is an error with the input data
● 412 - if the customer (same phone number) already exists

### Case 2. Login
Description
The API must allow a login with the phone number and password.
Validations
After 5 failed attempts, the user must be blocked.

Answers
The service must reply with:
● 200 - if the login is successful
● 400 - if the data sent is not in the valid format
● 401 - if the login information is incorrect
● 403 - if the user was blocked
● 429 - if the maximum number of login attempts exceeded

### Case 3. Membership payment
Description
The API must allow you to make a payment from your own membership, where the amount to
pay and fictitious credit card information (card number, expiration date,
name of the owner and CVV).
Validations
● A client can only pay for her own membership.
● The amount must not exceed 100,000 Mexican pesos.
● Each payment attempt must be stored.
● If the card number ends in an even digit, you must successfully answer the
transaction. Conversely, if it ends in an odd number, you must answer rejected.
● After 3 unsuccessful payment attempts, the user must be blocked to login and to
pay membership.
● The card number must be a field consisting of 16 digits. The first field does not
it can be a zero.
● The expiration date must be a field consisting of 4 digits, in format
MMYY.
● The CVV must be a 3 digit field.
Answers
The service must reply with:
● 200 - if payment is successful
● 400 - if the data sent does not have the valid format or the amount is greater than 100
thousand Mexican pesos
● 401 - if you have not logged in
● 403 - if the user tries to pay a membership of another client or the user is
locked
● 422 - if the payment was not successful
● 429 - if the client exceeded the maximum number of attempts allowed to pay the
membership

### Case 4. Consultation of membership payments
Description
Finally, the API must allow consulting the history of successful payments of the membership of the
user who consults.
Validations
You should not allow showing payments from another user.
Answers
The service must reply with:
● 200 - if the consultation is successful
● 401 - if you have not logged in
● 403 - if the user tries to obtain the list of membership payments from another client
Note: The language to use must be Golang. The database can be Postgres or MySQL.



# Architecture
The architecture used for this project is for Go modules for microservices due to the easy implementation of DevOps, maintainance and flexibility for adding new features and implementation.

The ptoject uses the following folder structure:

* assets (readme assets and swagger)
    * swagger (swagger yaml)
* build (docker and k8 [not implemented])
* cmd (neccesary functions to start the applications but not the Business logic)
    * api (single tools for running an application, could be more than 1 application)
        * handlers (routes and handler)
* configs (.env)
* deployments (docker-compose.yml)
* internal (decoupled core modules for the app)
    * config
    * platform
        * auth  (jwt)
        * database  (databse functions decoupled from DB)
            * schema    (mysql schemas)
        * web (http codes, functions and tools for the application web handling)
    * responses (http responses structure)
    * utils
* scripts (needed for initialize the application )
    * db
* tests


## Diagrams

### ER database diagram

![ER-Database](./assets/DataBase_ER.jpeg)
Format: ![Alt Text](url)

**Description**
The previous image describes the ER diagram for the Database with the following reasons for this particular choice
1. The application only needs 3 main entities for operations (Payments) and catalog (Users, Memberships)
2. The scope of the project does not require major operations such as login historial, password historial or multiple user memberships
3. The Users entity contains the login failures count because there's no need to implement a login historial table, at least not for this case
4. Memberships entity contains ExpDate, CVV fields which are not required fields due to security reasons and are disabled inside the code.
 Tier, Pricing and IsBloccked fields are not part of the scope but are needed because of common sense, the first indicates whether is normal or not, the second is the amount to pay when that feature is developed at last the IsBlocked field indicates if the membership is blocked for external reasons (not implemented)
 5. Payments entity is a many to many relationship which stoores all user payments attemprs whether are successfull or not and it only stores the UserID, MembershipID, Amount and the status (success or unsuccess).
 6. Status entity is not implemented due to time reasons but is neccesary as a Catalog for Users, Memberships and Payments

 **Improvements**
 1. Status table
 2. Login historial
 3. Memberships_Users table is the user needs to attach one more Card to his account
 4. Password historial

### Flow

![Register-Login](./assets/Register-Login.jpeg)
Format: ![Alt Text](url)

![Membership Payment](./assets/Membership_payment.jpeg)
Format: ![Alt Text](url)

![Payments Consult](./assets/Payments_consult.jpeg)
Format: ![Alt Text](url)


# RUN APLLICATION
## Requisites

* [Docker installed](https://docs.docker.com/get-docker/)
* [Docker-Compose installed](https://docs.docker.com/compose/install/)
* [Go installed](https://golang.org/doc/install)

## Export Env variables
* cd current_root_project_path
* export $(xargs <configs/example.env)

## Database

### Docker Compose
* cd current_root_project_path/deployments
* docker-compose up -d

## Run application
* cd current_root_project_path
* go mod download
* go run cmd/api/main.go