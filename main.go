package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type User struct {
	Name  string
	Email string
	Id    int
}
type Data struct {
	Items []User
}

func main() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	createTable()

	e.GET("/", func(c echo.Context) error {
		data := Data{}
		data.Items = getUsers()
		return c.Render(http.StatusOK, "index", data)
	})

	e.POST("/add", func(c echo.Context) error {
		u := User{
			Name:  c.FormValue("name"),
			Email: c.FormValue("email"),
		}
		createUser(&u)
		if u.Id == 0 {
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.Render(http.StatusOK, "user", u)
	})

	e.DELETE("/user/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusNotFound, "Not found")
		}
		deleted := deleteUser(id)

		if !deleted {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8081"))
}

func createUser(user *User) *User {
	db, _ := openDB()
	if db != nil {
		res, err := db.Exec("insert into User(name,email) values(?,?)", user.Name, user.Email)
		if err != nil {
			return user
		}

		id, err := res.LastInsertId()
		user.Id = int(id)
		return user
	}
	return user
}

func deleteUser(id int) bool {
	db, _ := openDB()
	if db != nil {
		res, _ := db.Exec("delete from user where id = ?", id)
		return res != nil
	}
	return false
}

func getUsers() []User {
	db, _ := openDB()
	if db != nil {
		res, err := db.Query("select * from user order by id desc")
		if err != nil {
			return nil
		}

		users := []User{}
		for res.Next() {
			var name string
			var email string
			var id int

			res.Scan(&id, &name, &email)

			u := User{
				Name:  name,
				Id:    id,
				Email: email,
			}
			users = append(users, u)
		}

		return users
	}
	return nil
}

func createTable() {
	db, _ := openDB()
	if db == nil {
		println("Unable to open db")

	}
	_, err := db.Exec("create table `user` (  `id` INTEGER PRIMARY KEY AUTOINCREMENT,  `name` VARCHAR(64) NULL,`email` VARCHAR(64) NULL )")
	fmt.Print(err)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	return db, nil
}
