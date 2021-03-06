package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/rs/cors"
	"gopkg.in/gographics/imagick.v3/imagick"
)

var (
	appconfig   = AppConfig{}
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	posterBin   []byte
	starBin     []byte
)

const (
	imagesPatch          = "https://image.tmdb.org/t/p/"
	storyHeight     uint = 1920
	storyWidth      uint = 1080
	posterMaxHeight uint = 880
	posterMaxWidth  uint = 880
)

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type AppConfig struct {
	Debug           bool   `json:"debug"`
	UseSSL          bool   `json:"use_ssl"`
	SSLCertFile     string `json:"ssl_cert"`
	SSLKeyFile      string `json:"ssl_key"`
	Upload          bool   `json:"upload"`
	MaxReviewLength int    `json:"max_length"`
}

func GetAppConfig(fname string) AppConfig {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	config := AppConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func init() {
	appconfig = GetAppConfig("config.json")
	imagick.Initialize()
	defer imagick.Terminate()

	starBin, _ = ioutil.ReadFile("images/assets/star.png")

	posterBin, _ = ioutil.ReadFile("images/assets/poster.png")
}

func main() {
	log.Println("Starting RUUARU...")
	log.Printf("Debug: %t, Max review length: %d, Upload: %t", appconfig.Debug, appconfig.MaxReviewLength, appconfig.Upload)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleAPI)
	srv := &http.Server{
		Addr:    ":3333",
		Handler: cors.Default().Handler(mux),
	}
    log.Fatal(srv.ListenAndServe())
}
