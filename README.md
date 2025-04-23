# Database Management System Final Project
This is Cason Parkinson's and my (Eli Bosch) submission for the DBMS final project. The back-end is coded in Go, and the front-end is coded in HTML, CSS, and JS. It's a basic housing assignment service for a University.

## Installation
To run the project on your local machine, clone the repo and create a .env file based on the structure of the .env.example file inside the server directory. Then create a database in MySQL. Compile the main.go file.
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
Once main is running, visit localhost:9010, and you'll visit the homepage to the website where you can go to the 7 subpages for different functionality.

## AWS Usage
To run the program on AWS, make sure that you have GoLang, MySQL, and Github installed on your instance. Create your database, and add your user. In main.go, change the http.listener from localhost:9010 to 0.0.0.0:80, and set up an S2 bucket utilizing the frontend folder from the repo. Change the CORS request in main.go to accept the S2 bucket. Then follow the steps in Installation. To run the backend, and fill it with example entries.
