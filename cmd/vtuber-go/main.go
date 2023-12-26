package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DanArmor/vtuber-go/pkg/config"
	"github.com/DanArmor/vtuber-go/pkg/controllers"
	"github.com/DanArmor/vtuber-go/pkg/middleware"
	"github.com/DanArmor/vtuber-go/pkg/setup"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"
	"github.com/watsonindustries/go-holodex"
)

type Options struct {
	ConfigPath string `long:"config" description:"Config path (with extension)" required:"true"`
}

const defaultTrustedProxy = "127.0.0.1:80"

func getBaseRouter(router *gin.Engine, basePath string) *gin.RouterGroup {
	if basePath != "" {
		return router.Group(basePath)
	} else {
		return router.Group("")
	}
}

func main() {
	var options Options
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		panic("Can't parse cmd args")
	}

	config, err := config.LoadConfig(options.ConfigPath)
	if err != nil {
		panic("Can't load config")
	}

	holodexConfig := holodex.NewConfiguration()
	holodexConfig.DefaultHeader["X-APIKEY"] = config.HolodexApiKey
	holodexConfig.UserAgent = "Vtuber-Go"

	service := controllers.NewService(
		setup.MustDatabaseSetup(config.DriverName, config.SqlUrl),
		config.TgBotToken,
		config.ExpirationHours,
		holodex.NewAPIClient(holodexConfig),
	)

	router := gin.Default()
	router.Use(middleware.CORSMiddleware)
	router.SetTrustedProxies([]string{defaultTrustedProxy})

	srv := &http.Server{
		Addr:    config.Ip + ":" + config.Port,
		Handler: router,
	}
	base := getBaseRouter(router, config.BasePath)

	api := base.Group("/api")
	api.POST("/search", service.SearchVtubers)
	api.GET("/orgs", service.GetOrgs)

	admin := api.Group("/admin")
	admin.POST("/vtubers", service.PostVtubers)

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// Stop main server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}
	log.Println("Server exiting")
}
