package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DanArmor/go-holodex"
	"github.com/DanArmor/vtuber-go/internal/config"
	"github.com/DanArmor/vtuber-go/internal/infrastructure/auth"
	"github.com/DanArmor/vtuber-go/internal/setup"
	"github.com/DanArmor/vtuber-go/pkg/controllers"
	"github.com/DanArmor/vtuber-go/pkg/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jessevdk/go-flags"

	ginzap "github.com/gin-contrib/zap"
	"go.uber.org/zap"
)

type Options struct {
	ConfigPath string `long:"config" description:"Config path (with extension)" required:"true"`
}

const defaultTrustedProxy = "127.0.0.1"

func getBaseRouter(router *gin.Engine, basePath string) *gin.RouterGroup {
	if basePath != "" {
		return router.Group(basePath)
	}
	return router.Group("")
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

	cfg := zap.NewDevelopmentConfig()

	logger, err := cfg.Build()
	if err != nil {
		panic("Can't create new logger: " + err.Error())
	}

	zap.ReplaceGlobals(logger)

	zap.L().Info("Service started")

	holodexConfig := holodex.NewConfiguration()
	holodexConfig.DefaultHeader["X-APIKEY"] = config.HolodexAPIKey
	holodexConfig.UserAgent = "Vtuber-Go"

	jwt, err := auth.NewJwtMaker(config.JwtSecretKey)
	if err != nil {
		panic("Can't create token maker")
	}

	service := controllers.NewService(
		setup.MustDatabaseSetup(config.DriverName, config.SQLURL),
		config.TgBotToken,
		config.ExpirationHours,
		config.TimeNotifyAfter,
		config.TimeStep,
		holodex.NewAPIClient(holodexConfig),
		jwt,
	)

	router := gin.New()
	// Установим zap в качестве логгера для Gin
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	// CORS
	router.Use(middleware.CORSMiddleware)

	if err := router.SetTrustedProxies([]string{defaultTrustedProxy}); err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    config.IP + ":" + config.Port,
		Handler: router,
	}
	base := getBaseRouter(router, config.BasePath)

	api := base.Group("/api")
	api.GET("/auth", service.AuthUser)

	protectedAPI := api.Group("")
	if !config.IsDebug {
		protectedAPI.Use(service.CheckToken)
	}
	protectedAPI.POST("/search", service.SearchVtubers)
	protectedAPI.POST("/select", service.SelectVtuber)
	protectedAPI.POST("/timezone", service.UserChangeTimezone)
	protectedAPI.GET("/timezone", service.UserGetTimezone)
	protectedAPI.GET("/orgs", service.GetOrgs)

	admin := api.Group("/admin")
	if !config.IsDebug {
		admin.Use(middleware.AdminVerify(config.AdminToken))
	}
	admin.POST("/vtubers", service.PostVtubers)

	// TODO change run command for scheduler
	go service.Scheduler.Run(context.Background())

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			zap.L().Info("Listen", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// Stop main server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server forced to shutdown", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
