# Database Management System Final Project
This is Cason Parkinson's and my (Eli Bosch) submission for the DBMS final. The back-end is coded in Go, and the front-end is coded in HTML, CSS, and JS. It's a basic housing assignment service for a University.

## Installation
To run the project on your local machine, clone the repo and create a .env file based on the .env.example file inside the server directory. Then create a database in MySQL. Compile the main.go file.
```bash
go build main.go
```
Then execute the file
```bash
./main
```
This migrates the database schema to your database using GORM's AutoMigrate function. Then, run the following commands during setup.sh file to fill the database with some example entries based on the University of Arkansas.
```bash
chmod +x setup.sh
./setup.sh
```
The database is now ready and filled for use by the front-end.

## Usage
TBD
