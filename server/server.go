package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kotakanbe/goval-dictionary/config"
	"github.com/kotakanbe/goval-dictionary/db"
	log "github.com/kotakanbe/goval-dictionary/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// Start starts CVE dictionary HTTP Server.
func Start(logDir string) error {
	e := echo.New()
	e.SetDebug(config.Conf.Debug)

	// Middleware
	//  e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// setup access logger
	logPath := filepath.Join(logDir, "access.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if _, err := os.Create(logPath); err != nil {
			return err
		}
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: f,
	}))

	// Routes
	e.Get("/health", health())
	e.Get("/cves/:family/:release/:id", getByCveID())
	e.Get("/packs/:family/:release/:pack", getByPackName())
	//  e.Post("/cpes", getByPackName())

	bindURL := fmt.Sprintf("%s:%s", config.Conf.Bind, config.Conf.Port)
	log.Infof("Listening on %s", bindURL)

	e.Run(standard.New(bindURL))
	return nil
}

// Handler
func health() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	}
}

// Handler
func getByCveID() echo.HandlerFunc {
	return func(c echo.Context) error {
		family := c.Param("family")
		release := c.Param("release")
		cveID := c.Param("id")
		log.Infof("%s %s %s", family, release, cveID)
		defs, err := db.GetByCveID(family, release, cveID)
		if err != nil {
			log.Infof("Failed to get by CveID: %s", err)
		}
		return c.JSON(http.StatusOK, defs)
	}
}

func getByPackName() echo.HandlerFunc {
	return func(c echo.Context) error {
		family := c.Param("family")
		release := c.Param("release")
		pack := c.Param("pack")
		log.Infof("%s %s %s", family, release, pack)
		defs, err := db.GetByPackName(family, release, pack)
		if err != nil {
			log.Infof("Failed to get by CveID: %s", err)
		}
		return c.JSON(http.StatusOK, defs)
	}
}
