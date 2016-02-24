package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olpainless/gorecess/backend/app"
	"github.com/op/go-logging"
)

var logFormat = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s} %{color:reset} %{message}`,
)

var logger = logging.MustGetLogger("main")

func main() {
	_InitLogging()
	_InitApp()
	_StartGin()
}

func _InitLogging() {

	//	var err error
	//	var logfile *os.File
	//	// Create logging file
	//	logfile, err = os.Open("logs/gorecess.log")

	//	if err != nil {
	//		panic(err)
	//	}

	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	logging.SetBackend(logFormatter)
	logger.Info("Logging initialized")

}

func _InitApp() {
	// TODO Load configuration
	dbFilename := "resources/recess.sqlite"
	sqlFilename := "resources/sql_statements.ini"

	// Create db instance
	db, err := app.NewDatabaseImpl(dbFilename, sqlFilename)
	if err != nil {
		logger.Panic(err)
	}

	facade := app.NewFacadeImpl(db)
	app.SetFacade(facade)
}

func _StartGin() {
	logger.Infof("Starting Gin on port %s", "8081")
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// Period
	router.POST("/period", app.PeriodPOST)
	router.GET("/period/:periodId", app.PeriodGET)
	//router.GET("/period/current", app.PeriodCurrentGET)
	router.PATCH("/period/:periodId", app.PeriodPATCH)
	router.DELETE("/period/:periodId", app.PeriodDELETE)
	router.POST("/period/:periodId/finalize", app.PeriodFinalizePOST)

	// Reservation
	router.POST("/period/:periodId/reservations", app.ReservationsPOST)
	router.DELETE("/period/:periodId/reservations/:reservationId", app.ReservationsDELETE)
	router.GET("/period/:periodId/reservations", app.ReservationsGET)

	// Schedule
	router.POST("/period/:periodId/schedule", app.SchedulePOST)
	router.GET("/period/:periodId/schedule", app.ScheduleGET)
	router.POST("/period/:periodId/schedule/post", app.ScheduleCreatePOST)

	// Schema
	router.GET("/schemas", app.SchemasAllGET)
	router.GET("/schemas/:schemaId", app.SchemasGET)
	router.POST("/schemas", app.SchemasPOST)

	// Timeslot
	router.POST("/timeslots", app.TimeslotsPOST)
	router.DELETE("/timeslots/:timeslotId", app.TimeslotsDELETE)
	router.PATCH("/timeslots/:timeslotId", app.TimeslotsPATCH)

	// Location
	router.POST("/locations", app.LocationPOST)

	err := router.Run(":8081")
	if err != nil {
		logger.Fatal(err.Error())
	}

}
