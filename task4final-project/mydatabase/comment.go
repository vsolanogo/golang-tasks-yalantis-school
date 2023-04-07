package mydatabase

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
	"task4/models"
	"strconv"
	"context"

)

func InsertComment(newComment models.Comment) models.Comment {
	db, err := sql.Open("postgres", ConnStr)

	if err != nil {
		log.Fatal(err)
	}

	newId := 0
	err = db.QueryRow(`INSERT INTO comment (taskid, text) values ($1, $2) 
	RETURNING (commentid);`, newComment.TaskId, newComment.Text).Scan(&newId)
	db.Close();

	if err != nil {
		log.Fatal(err)
	}

	newCommentQueried := SelectComment(strconv.Itoa(newId))
	
	return newCommentQueried
}

func SelectComment(id string) models.Comment {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	var i models.Comment

	err = db.QueryRow(`SELECT
	commentid, taskid, text, createdat, updatedat
	FROM
		comment
	WHERE
	commentid = $1`, id).Scan(&i.CommentId, &i.TaskId, &i.Text, &i.CreatedAt, &i.UpdatedAt)

	if err != nil {
		log.Fatal(err)
	}

	return i;
}

func DeleteComment(id string) {
	db, err := sql.Open("postgres", ConnStr)
	defer db.Close()

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.ExecContext(ctx,`DELETE FROM comment WHERE commentid=$1;`, id)
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