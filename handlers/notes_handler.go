package handlers

import (
	"database/sql"
	"net/http"
	"notes/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostNoteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")

		var note models.Note
		c.Bind(&note)
		query := `INSERT INTO notes (title, content, username) VALUES ($1,$2,$3) RETURNING id`
		row := db.QueryRow(query, note.Title, note.Content, username)
		err := row.Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error occured while inserting in database",
			})
			return
		}
		var id int
		err = row.Scan(&id)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error occured while retrieving id",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "note successfully posted",
			"id":      id,
		})

	}
}

func GetNotesHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		query := `SELECT id, title, content FROM notes WHERE username=$1`
		rows, err := db.Query(query, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "encountered some error in fetching notes",
			})
			return
		}
		var note models.Note
		var notes []models.Note
		for rows.Next() {
			rows.Scan(&note.ID, &note.Title, &note.Content)
			notes = append(notes, note)
		}
		c.JSON(http.StatusOK, notes)

	}
}

func GetNotesByIdHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error in converting route parameter to integer",
			})
		}
		query := `SELECT id, title, content FROM notes WHERE id=$1`
		row := db.QueryRow(query, id)
		err = row.Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error occured while retrieving note from database",
			})
			return
		}
		var note models.Note
		err = row.Scan(&note.ID, &note.Title, &note.Content)

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error occured while retrieving note from database",
			})
			return
		}
		c.JSON(http.StatusOK, note)

	}
}

func DeleteNoteHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		var err error
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error in converting route parameter to integer",
			})
		}
		query := `DELETE FROM notes WHERE id=$1`
		_, err = db.Exec(query, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error while deleting note from database",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "note successfully deleted",
		})

	}
}
