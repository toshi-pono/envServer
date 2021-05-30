package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type PostData struct {
	SecretKey   string  `json:"secretkey"`
	TimeSetting int     `json:"timeSetting"`
	EnvData     EnvData `json:"envData"`
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

const (
	verifyToken = "00000000000000000000000000000000"
)

var (
	db  *sqlx.DB
	bot *linebot.Client
)

func init() {
	rand.Seed(time.Now().UnixNano())
	_db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	db = _db

	// LINEのAPIを利用する設定
	_bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
	bot = _bot
}

func main() {
	e := echo.New()

	e.Renderer = &Renderer{
		templates: template.Must(template.ParseGlob("static/*.html")),
	}

	// LINE LIFE
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})

	// DATA
	e.GET("/data/latest", getLatestDataHandler)
	e.POST("/data/postData", postDataHandler)

	// LINEbot
	e.POST("/callback", linebotHandler)

	// プロセスを起動
	e.Start(":" + os.Getenv("PORT"))
}

// LINEbotのリプライメッセージを扱う
func linebotHandler(c echo.Context) error {
	log.Println("Accessed")

	events, err := bot.ParseRequest(c.Request())
	switch err {
	case nil:
	case linebot.ErrInvalidSignature:
		log.Println("PaeseRequest error:", err)
		return c.NoContent(http.StatusBadRequest)
	default:
		log.Println("PaeseRequest error:", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	// メッセージの種類によって処理を変更
	for _, event := range events {
		// LINEサーバーからのverify時
		if event.ReplyToken == verifyToken {
			return c.NoContent(http.StatusOK)
		}

		switch event.Type {
		// メッセージ受信時
		case linebot.EventTypeMessage:
			replyMessage := getReplyMessage(event)
			// 返信を送信
			if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
				log.Print(err)
			}
		// それ以外
		default:
			continue
		}
	}
	return c.NoContent(http.StatusOK)
}

const helpMessage = `使い方
temp: 最新の温度・湿度・気圧情報を取得します
dbstatus: データベースのレコード数を取得します
m5status: 観測機器のバッテリー情報を取得します`

func getReplyMessage(event *linebot.Event) string {
	switch message := event.Message.(type) {
	// テキストメッセージ
	case *linebot.TextMessage:
		return createReplyText(message.Text)
	// スタンプ
	case *linebot.StickerMessage:
		return fmt.Sprintf("sticker id is %v, stickerResourceType is %v", message.StickerID, message.StickerResourceType)
	// それ以外
	default:
		return helpMessage
	}
}

// テキストメッセージを分析して返信を作成する
func createReplyText(message string) string {
	if strings.Contains(message, "help") || strings.Contains(message, "使い方") {
		return helpMessage
	}
	if strings.Contains(message, "temp") {
		return latestRoomDataMessage()
	}
	if strings.Contains(message, "dbstatus") {
		return dbStatusMessage()
	}
	if strings.Contains(message, "m5status") {
		return m5StatusMessage()
	}
	return createRandomReply()
}

// ランダムな返信を返す
func createRandomReply() string {
	// TODO:  ランダムな返答を作成する
	return helpMessage
}

func latestRoomDataMessage() string {
	envData := EnvData{}
	err := db.Get(&envData, "SELECT * FROM weather ORDER BY created_at DESC LIMIT 1")
	if err == sql.ErrNoRows {
		return "データがありません"
	} else if err != nil {
		return fmt.Sprintf("db error: %v", err)
	}

	const format1 = "2006/01/02 15:04:05"
	return fmt.Sprintf("取得時刻: %s\n気温: %.1f度\n湿度: %.1f%%\n大気圧: %.1f hPa", envData.CreatedAt.Format(format1), envData.Temperature, envData.Humidity, envData.Pressure/100)
}

func dbStatusMessage() string {
	var count int
	err := db.Get(&count, "SELECT count(*) FROM weather")
	if err != nil {
		return fmt.Sprintf("db error: %v", err)
	}
	return fmt.Sprintf("weather: %d 件", count)
}

func m5StatusMessage() string {
	envData := EnvData{}
	err := db.Get(&envData, "SELECT * FROM weather ORDER BY created_at DESC LIMIT 1")
	if err == sql.ErrNoRows {
		return "データがありません"
	} else if err != nil {
		return fmt.Sprintf("db error: %v", err)
	}

	m5stickCStatus := "受信に問題があります"
	if envData.CreatedAt.After(time.Now().Add(-2 * time.Minute)) {
		m5stickCStatus = "正常稼働中"
	} else if envData.CreatedAt.After(time.Now().Add(-5 * time.Minute)) {
		m5stickCStatus = "データ受信不安定"
	}

	const format1 = "2006/01/02 15:04:05"
	return fmt.Sprintf("%s\n取得時刻: %s\nバッテリー残量: %.2f", m5stickCStatus, envData.CreatedAt.Format(format1), envData.Battery)
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
	if data.TimeSetting == 1 {
		t := time.Now().UTC()
		jst := time.FixedZone("JST", +9*60*60)
		data.EnvData.CreatedAt = t.In(jst)
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
