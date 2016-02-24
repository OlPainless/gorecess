package app

import (
	"database/sql"

	sqlite "github.com/mattn/go-sqlite3"
	"github.com/op/go-logging"
	"github.com/vaughan0/go-ini"
)

type Database interface {
	InsertSchema(schema *SchemaTO) (int, error)
	GetSchema(id int) (*SchemaTO, error)
	UpdateSchema(schema *SchemaTO) error

	//GetTimelotBySchema(schema int) ([]TimeslotTO, error)
	InsertTimeslot(timeslot *TimeslotTO) (int, error)
	DeleteSchemaTimeslots(id int) error

	InsertLocation(location *LocationTO) (int, error)
	//DeleteTimeslot(id int)

}

type DatabaseImpl struct {
	db  *sql.DB
	log *logging.Logger

	statements map[string]*sql.Stmt
}

func NewDatabaseImpl(database string, statementFile string) (*DatabaseImpl, error) {

	var err error

	d := new(DatabaseImpl)

	d.log = logging.MustGetLogger("database")
	d.db, err = sql.Open("sqlite3", database)

	if err != nil {
		switch errImpl := err.(type) {
		case sqlite.Error:
			println("SQLITE ERROR", errImpl.Error())
		default:
			println("Other error!", errImpl.Error())
		}

		return nil, err
	}

	d.statements = make(map[string]*sql.Stmt)

	err = d._LoadStatements(statementFile)
	return d, err
}

func (this *DatabaseImpl) _LoadStatements(statementFile string) error {

	var iniFile ini.File
	var err error

	iniFile, err = ini.LoadFile(statementFile)
	if err != nil {
		return err
	}

	for name, section := range iniFile {

		var statement *sql.Stmt
		var query string
		for _, v := range section {
			query = v
		}

		statement, err = this.db.Prepare(query)
		if err != nil {
			this.log.Error("Failed to create statement '", query, "'", err)
			return err
		}

		this.statements[name] = statement
	}

	return nil
}

func (this *DatabaseImpl) InsertSchema(schema *SchemaTO) (int, error) {

	var result sql.Result
	var err error
	var schemaId int64
	var tx *sql.Tx

	// FIXME Transactions do not work correctly

	// Start transaction
	tx, err = this.db.Begin()
	if err != nil {
		return -1, err
	}

	result, err = tx.Stmt(this.statements["INSERTSCHEMA"]).Exec(schema.Name, schema.State)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	schemaId, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Insert timeslots
	for _, ts := range schema.Timeslots {
		ts.Id = int(schemaId)
		_, err = this.InsertTimeslot(&ts)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	tx.Commit()

	return int(schemaId), nil
}

func (this *DatabaseImpl) UpdateSchema(schema *SchemaTO) error {
	//var result sql.Result
	var err error
	var tx *sql.Tx

	tx, err = this.db.Begin()
	if err != nil {
		this.log.Errorf("Failed to update schema. id=%s. Error=", schema.Id, err.Error())
		return err
	}

	// Update schema
	_, err = tx.Stmt(this.statements["UPDATESCHEMA"]).Exec(schema.Name, schema.Id)
	if err != nil {
		this.log.Errorf("Failed to update schema. id=%s. Error=", schema.Id, err.Error())
		tx.Rollback()
		return err
	}

	// Delete old timeslot
	err = this.__DeleteSchemaTimeslots(tx, schema.Id)
	if err != nil {
		this.log.Errorf("Failed to update schema. id=%s. Error=", schema.Id, err.Error())
		tx.Rollback()
		return err
	}

	// Insert new timeslots
	for _, ts := range schema.Timeslots {
		ts.Schema = schema.Id
		_, err = this.__InsertTimeslot(tx, &ts)
		if err != nil {
			this.log.Errorf("Failed to update schema. id=%s. Error=", schema.Id, err.Error())
			return err
		}
	}
	return nil

}

func (this *DatabaseImpl) GetSchema(id int) (*SchemaTO, error) {

	var err error
	var response SchemaTO
	var timeslots []TimeslotTO

	row := this.statements["GETSINGLESCHEMA"].QueryRow(id)
	err = row.Scan(&response.Id, &response.Name, &response.State)

	if err != nil {
		this.log.Errorf("Failed to get schema. id=%s. Error=", id, err.Error())
		return nil, err
	}

	// Timeslots
	timeslots, err = this.GetTimelotBySchema(id)
	response.Timeslots = timeslots

	return &response, err

}

func (this *DatabaseImpl) GetTimelotBySchema(schema int) ([]TimeslotTO, error) {

	var err error
	var rows *sql.Rows
	var timeslots []TimeslotTO = make([]TimeslotTO, 0, 4)

	// Fill timeslots
	rows, err = this.statements["GETSCHEMATIMESLOT"].Query(schema)
	defer rows.Close()
	if err != nil {
		this.log.Errorf("Failed to get timeslot by schema. id=%s. Error=", schema, err.Error())
		return nil, err
	}

	// Fetch each timeslot result
	for rows.Next() {
		var timeslotTO TimeslotTO
		var locationRows *sql.Rows

		err = rows.Scan(&timeslotTO.Id, &timeslotTO.Schema, &timeslotTO.Start, &timeslotTO.End)
		if err != nil {
			this.log.Errorf("Failed to get timeslot by schema. id=%s. Error=", schema, err.Error())
			return nil, err
		}

		// Fetch locations
		timeslotTO.Locations = make([]int, 0, 4)
		locationRows, err = this.statements["GETTIMESLOTLOCATIONS"].Query(timeslotTO.Id)
		if err != nil {
			this.log.Errorf("Failed to get timeslot by schema. id=%s. Error=", schema, err.Error())
			return nil, err
		}

		for locationRows.Next() {
			var location LocationTO
			err = locationRows.Scan(&location.Id, &location.Name)
			if err != nil {
				locationRows.Close()
				this.log.Errorf("Failed to get timeslot by schema. id=%s. Error=", schema, err.Error())
				return nil, err
			}
			timeslotTO.Locations = append(timeslotTO.Locations, location.Id)
		}

		timeslots = append(timeslots, timeslotTO)
		locationRows.Close()
	}

	return timeslots, err
}

func (this *DatabaseImpl) InsertTimeslot(timeslot *TimeslotTO) (int, error) {

	var err error
	var tx *sql.Tx
	var timeslotId int

	tx, err = this.db.Begin()
	if err != nil {
		this.log.Errorf("Failed to insert timeslot. id=%s. Error=", timeslot.Schema, err.Error())
		return -1, nil
	}

	timeslotId, err = this.__InsertTimeslot(tx, timeslot)
	if err != nil {
		this.log.Errorf("Failed to insert timeslot. id=%s. Error=", timeslot.Schema, err.Error())
		tx.Rollback()
		return -1, nil
	}

	return timeslotId, err
}

func (this *DatabaseImpl) __InsertTimeslot(tx *sql.Tx, timeslot *TimeslotTO) (int, error) {
	var result sql.Result
	var err error
	var timeslotId int64

	result, err = tx.Stmt(this.statements["INSERTTIMESLOT"]).Exec(timeslot.Schema, timeslot.Start, timeslot.End)
	if err != nil {
		this.log.Errorf("Failed to insert timeslot. id=%s. Error=", timeslot.Schema, err.Error())
		return -1, err
	}

	timeslotId, err = result.LastInsertId()
	if err != nil {
		this.log.Errorf("Failed to insert timeslot. id=%s. Error=", timeslot.Schema, err.Error())
		return -1, err
	}

	// Insert location array
	for _, loc := range timeslot.Locations {
		result, err = tx.Stmt(this.statements["INSERTTIMESLOTLOCATIONARRAY"]).Exec(timeslotId, loc)
		if err != nil {
			// TODO Rollback
			this.log.Errorf("Failed to insert timeslot. id=%s. Error=", timeslot.Schema, err.Error())
			return -1, err
		}
	}

	return int(timeslotId), nil
}

func (this *DatabaseImpl) DeleteSchemaTimeslots(id int) error {

	var tx *sql.Tx
	var err error

	tx, err = this.db.Begin()
	if err != nil {
		this.log.Errorf("Failed to delete schema timeslot. id=%s. Error=", id, err.Error())
		return err
	}

	err = this.__DeleteSchemaTimeslots(tx, id)
	if err != nil {
		this.log.Errorf("Failed to delete schema timeslot. id=%s. Error=", id, err.Error())
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil

}

func (this *DatabaseImpl) __DeleteSchemaTimeslots(tx *sql.Tx, id int) error {
	var rows *sql.Rows
	var err error
	var timeslotIds []int

	// Fetch affected timeslots
	rows, err = tx.Stmt(this.statements["GETSCHEMATIMESLOT"]).Query(id)
	if err != nil {
		this.log.Errorf("Failed to delete schema timeslot. id=%s. Error=", id, err.Error())
		return err
	}
	defer rows.Close()

	timeslotIds = make([]int, 0, 4)

	for rows.Next() {
		var timeslot TimeslotTO
		err = rows.Scan(&timeslot.Id, &timeslot.Schema, &timeslot.Start, &timeslot.End)
		if err != nil {
			this.log.Errorf("Failed to delete schema timeslot. id=%s. Error=", id, err.Error())
			return err
		}
		timeslotIds = append(timeslotIds, timeslot.Id)
	}

	_, err = tx.Stmt(this.statements["DELETESCHEMATIMESLOTS"]).Exec(id)
	if err != nil {
		this.log.Errorf("Failed to delete schema timeslot. id=%s. Error=", id, err.Error())
		return err
	}

	for _, ts := range timeslotIds {
		// Delete location array
		_, err = tx.Stmt(this.statements["DELETETIMESLOTLOCATIONARRAY"]).Exec(ts)
		if err != nil {
			this.log.Errorf("Failed to delete schema timeslot. id=%s. Error=", id, err.Error())
			return err
		}
	}

	return err
}

func (this *DatabaseImpl) InsertLocation(location *LocationTO) (int, error) {
	var result sql.Result
	var err error
	var locationId int64

	result, err = this.statements["INSERTLOCATION"].Exec(location.Name)
	if err != nil {
		//this.log.Errorf("Failed to insert location. id=%s. Error=", id, err.Error())
		return -1, err
	}

	locationId, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(locationId), nil
}
