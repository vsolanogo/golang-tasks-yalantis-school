package mydatabase

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"fmt"
)

var ConnStr = "user=postgres password=33123 sslmode=disable database=taskmanager"

func Initdb() {
	db, err := sql.Open("postgres", ConnStr)

	defer db.Close();

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("The connection to the DB was successfully initialized!")
	}

	_, err = db.Exec("DROP TABLE IF EXISTS comment")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS task")
	if err != nil { 
		fmt.Println(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS projectcolumn")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("DROP TABLE IF EXISTS project")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE project 
	(
		projectid INTEGER generated always as identity PRIMARY KEY,
		name VARCHAR(500),
		description VARCHAR(1000)
	)
	`)
	
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`INSERT INTO project 
	(
		name, description
	) values (
		'projectnumber one', 'this is description'
	), (
		'projectnumber two1', 'this is description'
	)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE projectcolumn 
	(
		projectcolumnid INTEGER generated always as identity PRIMARY KEY,
		projectid INT NOT NULL,
		name VARCHAR(255),
		position int,
		CONSTRAINT fkproject FOREIGN KEY(projectid) REFERENCES project(projectid) ON DELETE CASCADE
	)
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`INSERT INTO projectcolumn (projectid, name, position) values (1, 'firstcol', 1), (1, 'c2', 2),(1, 'c3', 3);`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE task 
	(
		taskid INTEGER generated always as identity PRIMARY KEY,
		projectcolumnid INT,
		name VARCHAR(255),
		description VARCHAR(5000),
		position int,
		CONSTRAINT fkprojectcolumn FOREIGN KEY(projectcolumnid) REFERENCES projectcolumn(projectcolumnid) ON DELETE CASCADE
	)
	`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`INSERT INTO task 
	(
		projectcolumnid, name, position
	) values ( 1, 'col2task', 1 ), ( 1, 'col2taskxxx', 2 ), ( 1, 'col2tttt', 3 ), ( 2, 'col2taskZZZ', 1 ), ( 2, 'col3taskZZZ', 2 )
	`)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`CREATE TABLE comment 
	(
		commentid INTEGER generated always as identity PRIMARY KEY,
		taskid INT,
		text VARCHAR(5000) NOT NULL,
		createdat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updatedat TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fktask FOREIGN KEY(taskid) REFERENCES task(taskid) ON DELETE CASCADE
	)
	`)

	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec(`INSERT INTO comment
	(
		taskid, text
		) values ( 1, 'comment' ), ( 1, 'comment' ), ( 1, 'comment' )
	`)
	if err != nil {
		fmt.Println(err)
	}

}