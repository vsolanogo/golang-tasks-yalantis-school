package mydatabase

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"fmt"
	"strconv"
	"task4/models"
)


func DeleteProject(id string) error {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		return err
	}

	res, err := db.Query(`DELETE FROM project WHERE projectid=$1;`, id)

	if err != nil {
		return err
	} else {
		fmt.Println(res)
	}

	return nil
}


func PutProject(updatedProject models.Project) error {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		return err
	}

	fmt.Println(updatedProject)

	res, err := db.Query(`UPDATE project SET name =$1, description = $2 WHERE projectid = $3`, updatedProject.Name, updatedProject.Description, updatedProject.ID)

	if err != nil {
		return err
	} else {
		fmt.Println(res)
	}

	return nil
}

func GetProject(id string) models.Project {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var i models.Project

	err = db.QueryRow(`SELECT
	projectid, name, description
	FROM
		project
	WHERE
	projectid = $1`, id).Scan(&i.ID, &i.Name, &i.Description)

	rows, err := db.Query(`SELECT
		projectcolumnid, projectid, name, position
		FROM
			projectcolumn
		WHERE
		projectid = $1;`, id)

	if err != nil {
		log.Fatal(err)
	}
	
	var queriedColumns []models.ProjectColumn

    for rows.Next() {
        var j models.ProjectColumn

		err := rows.Scan(&j.ProjectColumnId, &j.ProjectId, &j.Name, &j.Position)
        if err != nil {
			log.Fatal(err)
        }
        queriedColumns = append(queriedColumns, j)
    }

	i.ProjectColumns = queriedColumns
	
    if err = rows.Err(); err != nil {
		log.Fatal(err)
    }

	return i;
}

func GetProjects() []models.Project {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	} 

	rows, err := db.Query(`SELECT projectid, name, description FROM project 
	ORDER BY
	name ASC`)
	defer rows.Close()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rows)
	}

	var queriedProjects []models.Project
	
	for rows.Next() {
		var i models.Project

        err := rows.Scan(&i.ID, &i.Name, &i.Description)
        if err != nil {
			log.Fatal(err)
        }

        queriedProjects = append(queriedProjects, i)
    }
	
    if err = rows.Err(); err != nil {
		log.Fatal(err)
    }
	
	return queriedProjects
}


func PostProject(newProject models.Project) models.Project {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close();

	if err != nil {
		log.Fatal(err)
	}

	newId := 0
	err = db.QueryRow("INSERT INTO project (name, description) values ($1, $2) RETURNING projectid", newProject.Name, newProject.Description).Scan(&newId)

	if err != nil {
		log.Fatal(err)
	}

	newProjectQueried := GetProject(strconv.Itoa(newId))
	
	return newProjectQueried
}
