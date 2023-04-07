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

func PostProjectColumn(newColumn models.ProjectColumn) models.ProjectColumn {
	db, err := sql.Open("postgres", ConnStr)

	if err != nil {
		log.Fatal(err)
	}
	
	newColumnPosition := 0
	err = db.QueryRow("SELECT COUNT(*) FROM projectcolumn WHERE projectid = $1;", newColumn.ProjectId).Scan(&newColumnPosition)	
	newColumnPosition = 1+newColumnPosition

	newId := 0
	err = db.QueryRow(`INSERT INTO projectcolumn (projectid, name, position) values ($1, $2, $3) 
	RETURNING (projectcolumnid);`, newColumn.ProjectId, newColumn.Name, newColumnPosition).Scan(&newId)
	db.Close();

	if err != nil {
		log.Fatal(err)
	}

	newColumnQueried := GetProjectColumn(strconv.Itoa(newId))
	
	return newColumnQueried
}

func GetProjectColumn(id string) models.ProjectColumn {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var i models.ProjectColumn

	err = db.QueryRow(`SELECT
	projectcolumnid, projectid, name, position
	FROM
		projectcolumn
	WHERE
	projectcolumnid = $1`, id).Scan(&i.ProjectColumnId, &i.ProjectId, &i.Name, &i.Position)

	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query(`SELECT
		taskid, projectcolumnid, name, position
		FROM
			task
		WHERE
		projectcolumnid = $1
		ORDER BY position ASC;`, id)

	if err != nil {
		log.Fatal(err)
	}

	var queriedTasks []models.Task

    for rows.Next() {
        var j models.Task

		err := rows.Scan(&j.TaskId, &j.ProjectColumnId, &j.Name, &j.Position)
		
        if err != nil {
			log.Fatal(err)
        }
        queriedTasks = append(queriedTasks, j)
    }

	i.Tasks = queriedTasks
	
    if err = rows.Err(); err != nil {
		log.Fatal(err)
    }

	return i;
}

func UpdateColumnPosition(newcolumn models.ProjectColumn) {
	db, err := sql.Open("postgres", ConnStr)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var currentProjectColumn models.ProjectColumn

	err = db.QueryRow(`SELECT position FROM projectcolumn WHERE projectcolumnid = $1;`, newcolumn.ProjectColumnId).Scan(&currentProjectColumn.Position)
	err = db.QueryRow(`SELECT projectid FROM projectcolumn WHERE projectcolumnid = $1;`, newcolumn.ProjectColumnId).Scan(&currentProjectColumn.ProjectId)

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if (newcolumn.Position < currentProjectColumn.Position) {
		_, err := tx.ExecContext(ctx,`UPDATE projectcolumn SET position=position+1 WHERE position >= $1 AND position < $2 AND projectid=$3;`, newcolumn.Position, currentProjectColumn.Position, currentProjectColumn.ProjectId)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}

	} else if (newcolumn.Position > currentProjectColumn.Position) {
		_, err := tx.ExecContext(ctx,`UPDATE projectcolumn SET position=position-1 WHERE position <= $1 AND position > $2 AND projectid=$3;`, newcolumn.Position,currentProjectColumn.Position, currentProjectColumn.ProjectId)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}

	_, err  = tx.ExecContext(ctx,`UPDATE projectcolumn SET position=$1 WHERE projectcolumnid = $2;`, newcolumn.Position, newcolumn.ProjectColumnId)

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

func DeleteColumn(id string) {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var allColumnsCount = 1
	err = db.QueryRow("SELECT COUNT(*) FROM projectcolumn;").Scan(&allColumnsCount)	

	if(allColumnsCount == 1) {
		log.Fatal(`The last column can't be deleted`)
	}

	var columnToDel models.ProjectColumn

	columnToDel = GetProjectColumn(id)
 	fmt.Println(columnToDel)

	var columnToInsertTasksIntoID = 1
	err = db.QueryRow(`SELECT projectcolumnid FROM projectcolumn WHERE projectcolumnid !=$1 ORDER BY position ASC limit 1;`, id).Scan(&columnToInsertTasksIntoID)

	var tasksCountInColumnToInsertTasksInto = 0
	err = db.QueryRow(`SELECT COUNT(*) FROM tasks WHERE projectcolumnid =$1;`, columnToInsertTasksIntoID).Scan(&tasksCountInColumnToInsertTasksInto)


	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = tx.ExecContext(ctx,`UPDATE task SET position=position+$1 WHERE projectcolumnid=$2;`, tasksCountInColumnToInsertTasksInto, id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx,`UPDATE task SET projectcolumnid=$1 WHERE projectcolumnid=$2;`, columnToInsertTasksIntoID, id)

	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx,`DELETE FROM projectcolumn WHERE projectcolumnid=$1;`, id)

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