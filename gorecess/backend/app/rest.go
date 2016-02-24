package app

import (
	"github.com/gin-gonic/gin"
	"github.com/op/go-logging"

	"strconv"
)

var facade Facade
var logger = logging.MustGetLogger("facade")

func SetFacade(f Facade) {
	facade = f
}

// Error message
func _CreateErrorMessage(c *gin.Context, httpCode int, errorCode int) {
	c.JSON(httpCode, gin.H{"error": errorCode})
}

// Period

// Get a period.
// Path parameter "periodId"
func PeriodGET(c *gin.Context) {

}

// Creates a new period
func PeriodPOST(c *gin.Context) {

}

// Gets the current period
func PeriodCurrentGET(c *gin.Context) {

}

// Updates a period.
// Path parameter "periodId"
func PeriodPATCH(c *gin.Context) {

}

// Deletes a period
// Path parameter "periodId"
func PeriodDELETE(c *gin.Context) {

}

// Finalizes the schedule for period "periodId"
func PeriodFinalizePOST(c *gin.Context) {

}

// Reservations

// Creates a new reservation in period "periodId"
func ReservationsPOST(c *gin.Context) {

}

// Deletes a reservation in period "periodId"
func ReservationsDELETE(c *gin.Context) {

}

// Gets reservations in period "periodid".
// Optional query parameter "username" specifies username
func ReservationsGET(c *gin.Context) {

}

// Schedule

// Generates a new schedule in period "periodId"
func SchedulePOST(c *gin.Context) {

}

// Gets current schedule in period "periodId"
func ScheduleGET(c *gin.Context) {

}

// Set the current schedule for period "periodId"
func ScheduleCreatePOST(c *gin.Context) {

}

// Schema

// Creates a new schema
func SchemasPOST(c *gin.Context) {
	var schema SchemaTO
	var schemaId int
	var status Status
	var err error

	err = c.Bind(&schema)
	if err != nil {
		logger.Error("Failed to bind Schema from payload: ", err)
		_CreateErrorMessage(c, 400, STATUS_INVALID_PAYLOAD)
		return
	}

	schemaId, status = facade.CreateNewSchema(&schema)

	if status != STATUS_OK {
		_CreateErrorMessage(c, 400, STATUS_ERROR)
		return
	}

	c.JSON(200, gin.H{"id": schemaId})

}

// Gets a schema with id "schemaId"
func SchemasGET(c *gin.Context) {

	schemaIdStr := c.Params.ByName("schemaId")
	schemaId, err := strconv.Atoi(schemaIdStr)

	if err != nil {
		logger.Error("Failed to parse schema path variable: ", err)
		_CreateErrorMessage(c, 400, STATUS_INVALID_PATH_VARIABLE)
		return
	}

	schema, status := facade.GetSchema(schemaId)
	if status != STATUS_OK {
		logger.Error("Failed to get schema:", err)
		_CreateErrorMessage(c, 400, int(status))
		return
	}

	c.JSON(200, schema)

}

// Get all schemas
func SchemasAllGET(c *gin.Context) {

}

// Timeslots

// Creates a new timeslot
func TimeslotsPOST(c *gin.Context) {
	/*
		var timeslot TimeslotTO
		var timeslotId int
		var status Status
		var err error

		err = c.Bind(&timeslot)
		if err != nil {
			logger.Error("Failed to bind Timeslot from payload: ", err)
			_CreateErrorMessage(c, 400, STATUS_INVALID_PAYLOAD)
			return
		}

		timeslotId, status = facade.CreateNewTimeslot(&timeslot)

		if status != STATUS_OK {
			_CreateErrorMessage(c, 400, STATUS_ERROR)
			return
		}

		c.JSON(200, gin.H{"id": timeslotId})
	*/
}

// Deletes a timeslot with id "timeslotId"
func TimeslotsDELETE(c *gin.Context) {

}

// Updates a timeslot with id "timeslotId"
func TimeslotsPATCH(c *gin.Context) {

}

// Location
func LocationPOST(c *gin.Context) {
	facade.CreateNewLocation(nil)
}
