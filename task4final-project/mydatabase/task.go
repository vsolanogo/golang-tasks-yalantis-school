package mydatabase

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"task4/models"
	"strconv"
	"context"
	"fmt"
)

func InsertTask(newTask models.Task) models.Task {
	db, err := sql.Open("postgres", ConnStr)

	if err != nil {
		log.Fatal(err)
	}
	
	newTaskPosition := 0
	err = db.QueryRow("SELECT COUNT(*) FROM task WHERE projectcolumnid=$1;", newTask.ProjectColumnId).Scan(&newTaskPosition)	
	newTaskPosition = 1+newTaskPosition

	newId := 0
	err = db.QueryRow(`INSERT INTO task (projectcolumnid, name, description, position) values ($1, $2, $3, $4) 
	RETURNING (taskid);`, newTask.ProjectColumnId, newTask.Name, newTask.Description, newTaskPosition).Scan(&newId)
	db.Close();

	if err != nil {
		log.Fatal(err)
	}

	newTaskQueried := SelectTask(strconv.Itoa(newId))
	
	return newTaskQueried
}

func SelectTask(id string) models.Task {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var i models.Task
	 
	err = db.QueryRow(`SELECT
	taskid, projectcolumnid, name, position, COALESCE(description, '') as description
	FROM
		task
	WHERE
	taskid = $1`, id).Scan(&i.TaskId, &i.ProjectColumnId, &i.Name, &i.Position, &i.Description)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT
		commentid, taskid, text, createdat, updatedat
		FROM
			comment
		WHERE
		TaskId = $1 
		ORDER BY createdat ASC;`, id)

	if err != nil {
		log.Fatal(err)
	}

	var queriedComments []models.Comment

    for rows.Next() {
        var j models.Comment

		err := rows.Scan(&j.CommentId, &j.TaskId, &j.Text, &j.CreatedAt, &j.UpdatedAt)
		
        if err != nil {
			log.Fatal(err)
        }
        queriedComments = append(queriedComments, j)
    }

	i.Comments = queriedComments
	
    if err = rows.Err(); err != nil {
		log.Fatal(err)
    }

	return i;
}

func UpdateTaskPositionInScopeOfTheColumn(newTask models.Task) {
	db, err := sql.Open("postgres", ConnStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var currentTask models.Task

	err = db.QueryRow(`SELECT position FROM task WHERE taskid = $1;`, newTask.TaskId).Scan(&currentTask.Position)
	err = db.QueryRow(`SELECT projectcolumnid FROM task WHERE taskid = $1;`, newTask.TaskId).Scan(&currentTask.ProjectColumnId)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if (newTask.Position < currentTask.Position) {
		_, err := tx.ExecContext(ctx,`UPDATE task SET position=position+1 WHERE position >= $1 AND position < $2 AND projectcolumnid=$3;`, newTask.Position, currentTask.Position, currentTask.ProjectColumnId)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

	} else if (newTask.Position > currentTask.Position) {
		_, err := tx.ExecContext(ctx,`UPDATE task SET position=position-1 WHERE position <= $1 AND position > $2 AND projectcolumnid=$3;`, newTask.Position, currentTask.Position, currentTask.ProjectColumnId)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	_, err  = tx.ExecContext(ctx,`UPDATE task SET position=$1 WHERE taskid = $2;`, newTask.Position, newTask.TaskId)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
}


func UpdateTaskPosition(newTask models.Task) {
	db, err := sql.Open("postgres", ConnStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	
	var temporaryPositionInNewColumn = 0
	err = db.QueryRow(`SELECT COUNT(*) FROM task WHERE projectcolumnid = $1;`, newTask.ProjectColumnId).Scan(&temporaryPositionInNewColumn)
	temporaryPositionInNewColumn = temporaryPositionInNewColumn + 1

	var currentTask models.Task

	err = db.QueryRow(`SELECT position FROM task WHERE taskid = $1;`, newTask.TaskId).Scan(&currentTask.Position)
	err = db.QueryRow(`SELECT projectcolumnid FROM task WHERE taskid = $1;`, newTask.TaskId).Scan(&currentTask.ProjectColumnId)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx,`UPDATE task SET position=$1 WHERE taskid=$2;`, temporaryPositionInNewColumn, newTask.TaskId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx,`UPDATE task SET projectcolumnid=$1 WHERE taskid=$2;`, newTask.ProjectColumnId, newTask.TaskId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	if (newTask.Position < temporaryPositionInNewColumn) {
		_, err := tx.ExecContext(ctx,`UPDATE task SET position=position+1 WHERE position >= $1 AND position < $2 AND projectcolumnid=$3;`, newTask.Position, temporaryPositionInNewColumn, newTask.ProjectColumnId)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	_, err  = tx.ExecContext(ctx,`UPDATE task SET position=$1 WHERE taskid = $2;`, newTask.Position, newTask.TaskId)

	// clean empty position in origin column
	_, err = tx.ExecContext(ctx,`UPDATE task SET position=position-1 WHERE position > $1 AND projectcolumnid=$2;`, currentTask.Position, currentTask.ProjectColumnId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
}


func DeleteTask(id string) {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	var currentTask models.Task

	fmt.Println(id)

	err = db.QueryRow(`SELECT position FROM task WHERE taskid = $1;`, id).Scan(&currentTask.Position)
	if err != nil {
		log.Fatal(err)
	}

	err = db.QueryRow(`SELECT projectcolumnid FROM task WHERE taskid = $1;`, id).Scan(&currentTask.ProjectColumnId)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(currentTask.Position)
	fmt.Println(currentTask.ProjectColumnId)

	_, err = tx.ExecContext(ctx,`UPDATE task SET position=position-1 WHERE position > $1 AND projectcolumnid=$2;`, currentTask.Position, currentTask.ProjectColumnId)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx,`DELETE FROM task WHERE taskid=$1;`, id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
}