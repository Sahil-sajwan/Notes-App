package main

import (
	dbcon "notes/db"
	"notes/handlers"
	"notes/middlewares"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

func main() {

	/*query := `INSERT INTO notes (title, content) VALUES ($1, $2)`
	res, err := db.Exec(query, "second", "value")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)*/
	db, err := dbcon.OpenCon()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	r := gin.Default()
	r.Use(middlewares.OptionsMiddleware())

	protectedRoutes := r.Group("/api", middlewares.AuthMiddleware())
	r.POST("/users", handlers.RegisterUserHandler(db))
	r.POST("/login", handlers.LoginUserHandler(db))
	protectedRoutes.POST("/notes", handlers.PostNoteHandler(db))
	protectedRoutes.GET("/notes", handlers.GetNotesHandler(db))
	protectedRoutes.GET("/notes/:id", handlers.GetNotesByIdHandler(db))
	protectedRoutes.DELETE("/notes/:id", handlers.DeleteNoteHandler(db))
	//r.POST("/notes", middlewares.AuthMiddleware(), handlers.PostNoteHandler(db))
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
