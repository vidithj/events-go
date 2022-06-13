package db

import (
	"context"
	"events/model"

	"github.com/go-kivik/kivik"
)

const (
	documentConflict = "Conflict: Document update conflict."
	documentNotFound = "Not Found: missing"
)

//DB ...
type DB struct {
	eventsDB     *kivik.DB
	RegisterName string
	tag          string
	queryLimit   int
}

//InitDB ...
func InitDB(dbLoc, dbType, dbName, dbRegisterName, tag string, queryLimit int) (DB, error) {
	client, err := kivik.New(dbType, dbLoc)
	if err != nil {
		return DB{}, err
	}
	db1 := client.DB(context.Background(), dbName)
	return DB{
		db1,
		dbRegisterName,
		tag,
		queryLimit,
	}, err
}

func (db DB) GetEvents(username string) (model.Events, error) {
	var doc model.Events
	row := db.eventsDB.Get(context.Background(), username)
	if row == nil {
		return model.Events{}, nil
	}
	err := row.ScanDoc(&doc)
	if err != nil {
		return model.Events{}, err
	}
	return doc, nil
}

func (db DB) UpdateEvents(req model.UpdateEventRequest) error {
	userevents, err := db.GetEvents(req.Username)
	if err != nil {
		if err.Error() != documentNotFound {
			return err
		} else {
			_, err := db.eventsDB.Put(context.Background(), req.Username, model.Events{
				Id:       req.Username,
				Username: req.Username,
				Events:   []map[string]int64{req.Event},
			})
			return err
		}
	}
	userevents.Events = append(userevents.Events, req.Event)
	_, err = db.eventsDB.Put(context.Background(), req.Username, userevents)
	for err != nil {
		if err.Error() != documentConflict {
			return err
		}
		rev, err := db.ResolveConfict(req.Username)
		if err != nil {
			return err //can't fetch resolve
		}
		userevents.Revision = rev
		_, err = db.eventsDB.Put(context.Background(), req.Username, userevents)
	}
	return nil
}

func (db DB) ResolveConfict(username string) (string, error) {
	userinfo, err := db.GetEvents(username)
	if err != nil {
		return "", err
	}
	return userinfo.Revision, nil
}
