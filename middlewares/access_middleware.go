package middlewares

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AccessMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error in converting route parameter to integer",
			})
		}
		query := `SELECT username FROM notes WHERE id=$1`
		row := db.QueryRow(query, id)
		err = row.Err()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error occured while querying database",
			})
			return
		}
		var noteuser string

		err = row.Scan(&noteuser)
		if err != nil {

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error occured while querying database",
			})
			return
		}
		if noteuser != username {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You do not have permission to access this resource",
			})
			return
		}
		c.Next()
	}
}
