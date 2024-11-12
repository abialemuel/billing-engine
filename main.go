package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	userAPIhttp "github.com/abialemuel/billing-engine/pkg/user/api/http"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	mainCfg "github.com/abialemuel/billing-engine/config"
	"github.com/abialemuel/billing-engine/pkg/common/http/middleware/authguard"
	userBusiness "github.com/abialemuel/billing-engine/pkg/user/business"
	userRepository "github.com/abialemuel/billing-engine/pkg/user/modules/repository"
	"github.com/abialemuel/poly-kit/infrastructure/apm"
	"github.com/abialemuel/poly-kit/infrastructure/logger"
	loggerGen "github.com/abialemuel/poly-kit/infrastructure/logger"
	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	dd "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
)

var (
	APM *apm.APM
	log loggerGen.Logger
)

func main() {
	// config
	cfg := initializeConfig("config.yaml")
	log = initializeLogger(cfg)

	// initialize apm
	if cfg.Get().APM.Enabled {
		host := ""
		apmType := -1
		if cfg.Get().APM.DDAgentHost != "" {
			host = fmt.Sprintf("%s:%d", cfg.Get().APM.DDAgentHost, cfg.Get().APM.DDAgentPort)
			apmType = apm.DatadogAPMType
		} else {
			host = fmt.Sprintf("%s:%d", cfg.Get().APM.Host, cfg.Get().APM.Port)
			apmType = apm.OpenTelemetryAPMType
		}

		apmPayload := apm.APMPayload{
			ServiceHost:    &host,
			ServiceName:    cfg.Get().App.Name,
			ServiceEnv:     cfg.Get().App.Env,
			ServiceVersion: cfg.Get().App.Version,
			ServiceTribe:   cfg.Get().App.Name,
			SampleRate:     cfg.Get().APM.Rate,
		}
		APM, err := apm.NewAPM(apmType, apmPayload)
		if err != nil {
			log.Get().Error(err)
			panic(err)
		}
		fmt.Println("APM started...")
		defer APM.EndAPM()
	}

	// init postgresql
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Get().Postgres.Host, cfg.Get().Postgres.Port, cfg.Get().Postgres.User, cfg.Get().Postgres.Password, cfg.Get().Postgres.DB)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Get().Error(err)
		panic(err)
	}

	userRepo := userRepository.NewPgDBRepository(db)

	// init userService
	userService := newUserService(cfg, db, userRepo)

	// Init HTTP client
	e := echo.New()
	e.Use()
	e.Use(mw.Recover())
	// opentelemetry echo middleware
	e.Use(otelecho.Middleware(cfg.Get().App.Name))
	// datadog echo middleware
	e.Use(dd.Middleware())
	e.Use(mw.CORSWithConfig(
		mw.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization, "X-Service"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}))

	//health check
	e.GET("/", func(c echo.Context) error {
		return c.NoContent(200)
	})
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// run server
	go func() {
		address := fmt.Sprintf(":%d", cfg.Get().App.Port)

		if err := e.Start(address); err != nil {
			log.Get().Info("shutting down the server")
		}
	}()

	authGuard := authguard.NewAuthGuard(*cfg.Get())

	// Register API
	userHandler := userAPIhttp.NewHandler(userService, cfg.Get())
	userAPIhttp.RegisterPath(e, userHandler, authGuard)

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Get().Fatal(err)
	}
}

func initializeLogger(cfg mainCfg.Config) logger.Logger {
	fmt.Printf("%s started...\n", cfg.Get().App.Name)
	log := logger.New().Init(logger.Config{
		Level:  cfg.Get().Log.Level,
		Format: cfg.Get().Log.Format,
	})
	return log
}

func initializeConfig(path string) mainCfg.Config {
	cfg := mainCfg.New()
	err := cfg.Init(path)
	if err != nil {
		fmt.Errorf("failed to load config: %s", err.Error())
		panic(err)
	}
	return cfg
}

func newUserService(
	cfg mainCfg.Config,
	db *gorm.DB,
	userRepo *userRepository.PgDBRepository,
) userBusiness.UserService {
	userService := userBusiness.NewUserService(userRepo, cfg, db)
	return userService
}
