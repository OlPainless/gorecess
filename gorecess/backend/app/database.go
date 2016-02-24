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

	//GetTimelotBySchema(schema int) ([]TimeslotTO, error)
	InsertTimeslot(timeslot *TimeslotTO) (int, error)

	InsertLocation(location *LocationTO) (int, error)
	//DeleteTimeslot(id int)

}

type DatabaseImpl struct {
	db  *sql.DB
	log *logging.Logger

	statements map[string]*sql.Stmt

	// Statements
	insertSchemaStmt    *sql.Stmt
	getSingleSchemaStmt *sql.Stmt
	getAllSchemaStmt    *sql.Stmt

	// Locations
	insertLocationStmt    *sql.Stmt
	getSingleLocationStmt *sql.Stmt
	//getSchemaLocationStmt *sql.Stmt

	// Timeslots
	insertTimeslotStmt              *sql.Stmt
	insertTimeslotLocationArrayStmt *sql.Stmt
	deleteTimeslotStmt              *sql.Stmt
	deleteTimeslotLocationArrayStmt *sql.Stmt
	getSingleTimeslotStmt           *sql.Stmt
	getSchemaTimeslotStmt           *sql.Stmt

	testStmt *sql.Stmt
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

	//	err = d._CreateStatements(statements)
	//	return d, err

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

//func (this *DatabaseImpl) _CreateStatements() error {
//	var err error

//	this.insertSchemaStmt, err = this.db.Prepare("INSERT INTO Schema(name, state) VALUES(?, ?);")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.getSingleSchemaStmt, err = this.db.Prepare("SELECT * FROM Schema WHERE id=?;")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.getSingleLocationStmt, err = this.db.Prepare("SELECT * FROM Location WHERE id=?;")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	//	this.getSchemaLocationStmt, err = this.db.Prepare("SELECT * FROM Location WHERE schema=?;")
//	//	if err != nil {
//	//		this.log.Error("Failed to create statement", err)
//	//		return err
//	//	}

//	this.getSingleTimeslotStmt, err = this.db.Prepare("SELECT * FROM Timeslot WHERE id=?;")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.getSchemaTimeslotStmt, err = this.db.Prepare("SELECT * FROM Timeslot WHERE schema=?;")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.insertLocationStmt, err = this.db.Prepare("INSERT INTO Location(name) VALUES(?);")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.insertTimeslotStmt, err = this.db.Prepare("INSERT INTO Timeslot(schema, start, end) VALUES(?,?,?);")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.insertTimeslotLocationArrayStmt, err = this.db.Prepare("INSERT INTO Timeslot_location_array(timeslot, location) VALUES(?,?);")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	this.deleteTimeslotLocationArrayStmt, err = this.db.Prepare("DELETE FROM Timeslot_location_array WHERE timeslot=?;")
//	if err != nil {
//		this.log.Error("Failed to create statement", err)
//		return err
//	}

//	return nil
//}

func (this *DatabaseImpl) InsertSchema(schema *SchemaTO) (int, error) {

	var result sql.Result
	var err error
	var schemaId int64

	result, err = this.statements["INSERTSCHEMA"].Exec(schema.Name, schema.State)
	if err != nil {
		return -1, err
	}

	schemaId, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(schemaId), nil
}

func (this *DatabaseImpl) GetSchema(id int) (*SchemaTO, error) {

	var err error
	var response SchemaTO
	//var locations []LocationTO
	var timeslots []TimeslotTO
	//	var rows *sql.Rows

	//	rows, err = this.getSingleSchemaStmt.Query(id)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if rows == nil {
	//		return nil, err
	//	}

	row := this.statements["GETSINGLESCHEMA"].QueryRow(id)
	err = row.Scan(&response.Id, &response.Name, &response.State)

	if err != nil {
		return nil, err
	}

	// Locations
	//	locations, err = this.GetLocationBySchema(id)
	//	if err != nil {
	//		return nil, err
	//	}
	//	locationIds := make([]int, len(locations))
	//	for idx, loc := range locations {
	//		locationIds[idx] = loc.Id
	//	}

	// Timeslots
	timeslots, err = this.GetTimelotBySchema(id)

	//	if err != nil {
	//		return nil, err
	//	}
	//	timeslotIds := make([]int, len(timeslots))
	//	for idx, slot := range timeslots {
	//		timeslotIds[idx] = slot.Id
	//	}

	response.Timeslots = timeslots

	return &response, err

}

//func (this *DatabaseImpl) GetLocationBySchema(schema int) ([]LocationTO, error) {

//	var err error
//	var rows *sql.Rows
//	var locations []LocationTO = make([]LocationTO, 0, 4)

//	// Fill locations
//	rows, err = this.getSchemaLocationStmt.Query(schema)
//	if err != nil {
//		return nil, err
//	}

//	// Fetch each timeslot result
//	for rows.Next() {
//		var locationTO LocationTO

//		err = rows.Scan(&locationTO.Id, &locationTO.Name, &locationTO.Schema)
//		if err != nil {
//			return nil, err
//		}

//		locations = append(locations, locationTO)
//	}
//	err = rows.Err()
//	if err != nil {
//		return nil, err
//	}
//	rows.Close()

//	return locations, err

//}

func (this *DatabaseImpl) GetTimelotBySchema(schema int) ([]TimeslotTO, error) {

	var err error
	var rows *sql.Rows
	var timeslots []TimeslotTO = make([]TimeslotTO, 0, 4)

	// Fill timeslots
	rows, err = this.statements["GETSCHEMATIMESLOT"].Query(schema)
	if err != nil {
		return nil, err
	}

	// Fetch each timeslot result
	for rows.Next() {
		var timeslotTO TimeslotTO

		err = rows.Scan(&timeslotTO.Id, &timeslotTO.Schema, &timeslotTO.Start, &timeslotTO.End)
		if err != nil {
			return nil, err
		}

		timeslots = append(timeslots, timeslotTO)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	rows.Close()

	return timeslots, err
}

func (this *DatabaseImpl) InsertTimeslot(timeslot *TimeslotTO) (int, error) {
	var result sql.Result
	var err error
	var timeslotId int64

	result, err = this.statements["INSERTTIMESLOT"].Exec(timeslot.Schema, timeslot.Start, timeslot.End)
	if err != nil {
		return -1, err
	}

	timeslotId, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(timeslotId), nil
}

func (this *DatabaseImpl) InsertLocation(location *LocationTO) (int, error) {
	var result sql.Result
	var err error
	var locationId int64

	result, err = this.statements["INSERTLOCATION"].Exec(location.Name)
	if err != nil {
		return -1, err
	}

	locationId, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(locationId), nil
}
