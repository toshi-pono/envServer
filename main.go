package main

import (
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

type EnvData struct {
	Id          int       `json:"id,omitempty" db:"id"`
	Temperature float32   `json:"temperature" db:"temperature"`
	Humidity    float32   `json:"humidity" db:"humidity"`
	Pressure    float32   `json:"pressure" db:"pressure"`
	Battery     float32   `json:"battery" db:"battery"`
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
	e.POST("/api/data/postDatas", HandleAPIPostDatas)

	// プロセスを起動
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func HandleAPIPostDatas(c echo.Context) error {
	req := new(EnvData)
	if err := c.Bind(req); err != nil {
		return err
	}
	// データベースに追加
	_, err := db.Exec("INSERT INTO weather (temperature, humidity, pressure, battery, created_at) VALUES (?, ?, ?, ?, ?)", req.Temperature, req.Humidity, req.Pressure, req.Battery, req.CreatedAt)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.String(http.StatusOK, "OK")
}
