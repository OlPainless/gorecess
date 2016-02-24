package app

import (
	"github.com/op/go-logging"
)

type Facade interface {
	// Schema
	CreateNewSchema(schema *SchemaTO) (int, Status)
	GetSchema(id int) (*SchemaTO, Status)

	// Timeslot
	CreateNewTimeslot(timeslot *TimeslotTO) (int, Status)

	// Location
	CreateNewLocation(location *LocationTO) (int, Status)
}

type FacadeImpl struct {
	db  Database
	log *logging.Logger
}

func NewFacadeImpl(db Database) *FacadeImpl {
	f := new(FacadeImpl)
	f.log = logging.MustGetLogger("facade")
	f.db = db
	return f
}

func (this *FacadeImpl) CreateNewSchema(schema *SchemaTO) (int, Status) {
	schemaId, err := this.db.InsertSchema(schema)
	if err != nil {
		this.log.Error("Failed to create new schema: ", err)
		return -1, STATUS_ERROR
	}
	this.log.Infof("Created new schema with id=%d", schemaId)
	return schemaId, STATUS_OK
}

func (this *FacadeImpl) GetSchema(id int) (*SchemaTO, Status) {
	schema, err := this.db.GetSchema(id)
	if err != nil {
		this.log.Error("Failed to get schema:", err)
		return nil, STATUS_ERROR
	}
	return schema, STATUS_OK
}

func (this *FacadeImpl) CreateNewTimeslot(timeslot *TimeslotTO) (int, Status) {
	timeslotId, err := this.db.InsertTimeslot(timeslot)
	if err != nil {
		this.log.Error("Failed to create new timeslot: ", err)
		return -1, STATUS_ERROR
	}
	this.log.Infof("Created new timeslot with id=%d", timeslotId)
	return timeslotId, STATUS_OK
}

func (this *FacadeImpl) CreateNewLocation(location *LocationTO) (int, Status) {

	panic("TODO")
	return -1, -1
}
