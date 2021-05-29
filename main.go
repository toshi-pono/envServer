package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type PostData struct {
	SecretKey string  `json:"secretkey"`
	EnvData   EnvData `json:"envData"`
}

type EnvData struct {
	Id          int       `json:"id,omitempty" db:"id"`
	Temperature float64   `json:"temperature" db:"temperature"`
	Humidity    float64   `json:"humidity" db:"humidity"`
	Pressure    float64   `json:"pressure" db:"pressure"`
	Battery     float64   `json:"battery" db:"battery"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// テンプレートのレンダラ―を作る
type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

var (
	db *sqlx.DB
)

func main() {
	fmt.Println("hello")
	_db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db = _db

	e := echo.New()

	e.Renderer = &Renderer{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	e.GET("/data/latest", getLatestDataHandler)
	e.POST("/data/postData", postDataHandler)

	// プロセスを起動
	e.Start(":" + os.Getenv("PORT"))
}

func getLatestDataHandler(c echo.Context) error {
	envData := EnvData{}
	err := db.Get(&envData, "SELECT * FROM weather ORDER BY created_at DESC LIMIT 1")
	if err == sql.ErrNoRows {
		return c.NoContent(http.StatusNotFound)
	} else if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.JSON(http.StatusOK, envData)
}

func postDataHandler(c echo.Context) error {
	data := new(PostData)
	if err := c.Bind(data); err != nil {
		return err
	}
	if data.SecretKey != os.Getenv("POST_DATA_KEY") {
		return c.String(http.StatusForbidden, "Forbidden")
	}
	req := data.EnvData
	// データベースに追加する
	fmt.Println(req)
	_, err := db.Exec("INSERT INTO weather (temperature, humidity, pressure, battery, created_at) VALUES ($1, $2, $3, $4, $5)", req.Temperature, req.Humidity, req.Pressure, req.Battery, req.CreatedAt)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.String(http.StatusOK, "OK")
}
