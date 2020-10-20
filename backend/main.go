package main

import (
	"net/http"

	elogrus "github.com/dictor/echologrus"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// CustomValidator is struct validator for request input
type CustomValidator struct {
	validator *validator.Validate
}

// Validate is just renamed function of struct validate method
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

var (
	// Logger is global logger reference to Echo object's logger
	Logger       elogrus.EchoLogger
	gitHash      string = "N/A"
	gitTag       string = "N/A"
	buildDate    string = "N/A"
	displayDebug string = "false"
)

func main() {
	e := echo.New()
	Logger = elogrus.Attach(e)
	e.Validator = &CustomValidator{validator: validator.New()}
	Logger.Infof("leave-it-us backend %s (%s) : %s\n", gitTag, gitHash, buildDate)

	if displayDebug == "true" {
		Logger.SetLevel(log.DEBUG)
		Logger.Debugln("run in debug mode")
	}

	e.GET("/api/tag", ReadHandler(Tag{}))
	e.POST("/api/tag", CreateHandler(Tag{}))
	e.DELETE("/api/tag", DeleteHandler(Tag{}))

	//static resource
	e.GET("/*", VueRouterStatic("static"))

	SetDBHander(TestDBHander())
	Migrate()
	e.Logger.Fatal(e.Start(":80"))
}

// ParseSQLErrorToResponse parses sql error message and returns echo response
func ParseSQLErrorToResponse(err error, c echo.Context) error {
	if err == nil {
		return c.NoContent(http.StatusOK)
	}
	msg := err.Error()
	var (
		code    int
		respmsg string
	)
	if msg[0:6] == "UNIQUE" {
		code = http.StatusForbidden
		respmsg = "req_violate_unique"
	} else {
		code = http.StatusInternalServerError
		respmsg = "unexpected_error"
	}
	return c.JSON(code, map[string]string{"msg": respmsg})
}
