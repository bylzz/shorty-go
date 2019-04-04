package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/color"
	"github.com/speps/go-hashids"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"roboncode.com/go-urlshortener/stores"
	"strconv"
	"time"
)

type Counter struct {
	Value int `bson:"value"`
}

type Link struct {
	ID       interface{} `json:"id,omitempty" bson:"_id,omitempty"`
	Code     string      `json:"code" bson:"code"`
	LongUrl  string      `json:"longUrl" bson:"longUrl"`
	ShortUrl string      `json:"shortUrl,omitempty" bson:"shortUrl,omitempty"`
	Created  time.Time   `json:"created" bson:"created"`
}

var store stores.Store
var h *hashids.HashID

func setupHashIds() *hashids.HashID {
	hd := hashids.NewData()
	hd.Salt = viper.GetString("hashSalt")
	hd.MinLength = viper.GetInt("hashMin")
	h, _ := hashids.NewWithData(hd)
	return h
}

func readConfig() {
	viper.SetDefault("mongoUrl", "mongodb://localhost:27017")
	viper.SetDefault("database", "shorturls")
	viper.SetDefault("hashSalt", "shorturls")
	viper.SetDefault("hashMin", 5)
	viper.SetDefault("address", ":1323")
	viper.SetDefault("baseUrl", "")
	viper.SetDefault("authKey", "")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println(color.Red("No config file found. Using default settings"))
	} else {
		log.Println(color.Green("Config found -- overriding defaults"))
	}

	mongoUrl := os.Getenv("MONGO_DB")
	if mongoUrl != "" {
		viper.Set("mongoUrl", mongoUrl)
	}
}

func setupBadger() {
	opts := badger.DefaultOptions
	opts.Dir = "./data/badger"
	opts.ValueDir = "./data/badger"
	db, err := badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("answer"), []byte(`{"num":6.13,"strs":["a","b"]}`))
		return err
	})

	err = db.View(func(txn *badger.Txn) error {
		var err error
		var item *badger.Item
		item, err = txn.Get([]byte("answer"))
		if err != nil {
			log.Fatal(err)
		}

		var dat = new(struct {
			Num  float32  `json:"num,omitempty"`
			Strs []string `json:"strs,omitempty"`
		})
		valCopy, _ := item.ValueCopy(nil)
		if err := json.Unmarshal(valCopy, &dat); err != nil {
			panic(err)
		}
		fmt.Println(dat.Num, dat.Strs[1])

		//err = db.View(func(txn *badger.Txn) error {
		//	item, _ := txn.Get([]byte("answer"))
		//	valCopy, _ := item.ValueCopy(nil)
		//	fmt.Printf("The answer is: %s\n", valCopy)
		//	return nil
		//})

		return nil
	})
}

func main() {
	readConfig()
	h = setupHashIds()
	store = stores.NewMongoStore()
	//setupBadger()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "query:key",
		Skipper: func(e echo.Context) bool {
			switch e.Path() {
			case "/", "/404", "/:code", "/*":
				return true
			}
			return false
		},
		Validator: func(key string, e echo.Context) (bool, error) {
			return key == viper.GetString("authKey"), nil
		},
	}))

	// Routes
	e.POST("/shorten", CreateLink)
	e.GET("/links", GetLinks)
	e.GET("/links/:code", GetLink)
	e.DELETE("/links/:code", DeleteLink)
	e.File("/", "public/index.html")
	e.File("/404", "public/404.html")
	e.GET("/:code", RedirectToUrl)
	e.File("/*", "public/404.html")

	// Start server
	e.Logger.Fatal(e.Start(viper.GetString("address")))
}

// Handler
func CreateLink(c echo.Context) error {
	var body = new(struct {
		Url string `json:"url"`
	})

	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Url == "" {
		return c.JSON(http.StatusBadRequest, `Missing required property "url"`)
	}

	counter := store.IncCount()
	if code, err := h.Encode([]int{counter}); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else if link, err := store.Create(code, body.Url); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusOK, link)
	}
}

func GetLinks(c echo.Context) error {
	skip, _ := strconv.ParseInt(c.QueryParam("s"), 10, 64)
	limit, _ := strconv.ParseInt(c.QueryParam("l"), 10, 64)
	links := store.List(limit, skip)
	return c.JSON(http.StatusOK, links)
}

func GetLink(c echo.Context) error {
	if link, err := store.Read(c.Param("code")); err != nil {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.JSON(http.StatusOK, link)
	}
}

func DeleteLink(c echo.Context) error {
	if count := store.Delete(c.Param("code")); count == 0 {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.NoContent(http.StatusOK)
	}
}

func RedirectToUrl(c echo.Context) error {
	if link, err := store.Read(c.Param("code")); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/404")
	} else {
		return c.Redirect(http.StatusMovedPermanently, link.LongUrl)
	}
}
