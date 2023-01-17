package main

import (
	"apiserver/config"
	"apiserver/model"
	"apiserver/router"
	"errors"

	"net/http"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "apiserver config file path")
)

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// init db
	model.DB.Init()
	defer model.DB.Close()

	// set gin mode
	// gin.SetMode(viper.GetString("runmode"))
	// Create the Gin engine
	g := gin.New()

	// gin middlewares
	middlewares := []gin.HandlerFunc{}

	// Routes
	router.Load(
		// Cores.
		g,
		// Middlewares.
		middlewares...,
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	err := http.ListenAndServe(viper.GetString("addr"), g)
	if err != nil {
		log.Printf("the err is %v", err)
	}
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		// Ping the server by sending a GET request to `/health`.
		log.Printf("the url is %v", viper.GetString("url")+"/sd/health")
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("cannot connect to the router")
}
