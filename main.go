package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"github.com/labstack/echo/v4"
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
	id := 1
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	data := Data{
		Items: []User{
			{
				Name:  "sachin",
				Email: "sachinbhankhar@gmail.com",
				Id:    id,
			},
		},
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", data)
	})

	e.POST("/add", func(c echo.Context) error {
		id += 1
		u := User{
			Name:  c.FormValue("name"),
			Email: c.FormValue("email"),
			Id:    id,
		}
		data.Items = append(data.Items, u)
		return c.Render(http.StatusOK, "user", u)
	})

	e.DELETE("/user/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusNotFound, "Not found")
		}
		index := -1
		for i, item := range data.Items {
			if item.Id == id {
				index = i
				break
			}
		}
		if index == -1 {
			return c.NoContent(http.StatusNotFound)
		}
		data.Items = append(data.Items[index:], data.Items[index+1:]...)
		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
