# Database Management System Final Project
This is Cason Parkinson's and my (Eli Bosch) submission for the DBMS final. The back-end is coded in Go, and the front-end is coded in HTML, CSS, and JS.

## Installation
To run the project, clone the repo, create a new database, update the .env file, and run the project.
``` bash
cd DBMS-Final/backend/cmd/server
go build main.go
./main
```


This auto-migrates the database schema to your machine. Next run this in your MySQL database.
```MySQL
ALTER TABLE students
MODIFY COLUMN student_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
ADD PRIMARY KEY (student_id);
```
There is some issues with the students table autoIncrement during the automigration. Then update the setup.sh file with your database information and run it.
``` bash
cd ../../..
chmod +x setup.sh
./setup.sh
```
This fills in some example entries based on the University of Arkansas in your database.

## Usage
TBD
