/*
	The sde package is a fully functional library for use with the DUST514
	Static Data Export.  The package automatically can download and manage
	multiple versions of the SDE.
*/
package sde

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Database driver
	"os"
	"strings"
	"time"
)

// SDE is a struct containing the database object, the version of the SDE
// and many methods for working with the SDE.
type SDE struct {
	DB      *sql.DB
	Version string
}

// Open will open our SDE of the version specified.
func Open(Version string) (SDE, error) {
	for k := range Versions {
		if k == Version {
			s := getsde(k)
			return s, nil
		}
	}
	return SDE{}, errors.New("No such version:" + Version)
}

// GetType returns an SDEType of the given TypeID
func (s *SDE) GetType(id int) (SDEType, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM CatmaTypes WHERE TypeID == '%v'", id))
	if err != nil {
		return SDEType{}, err
	}
	if rows.Next() {
		var nTypeID int
		var nTypeName string

		rows.Scan(&nTypeID, &nTypeName)
		return SDEType{s, nTypeID, nTypeName, make(map[string]interface{})}, nil
	}
	return SDEType{}, errors.New("no such type")
}

// GetTypeWhereNameContains should be thought of as a search function that
// checks the display name.
func (s *SDE) GetTypeWhereNameContains(name string) ([]SDEType, error) {
	values := make([]SDEType, 0)
	rows, err := s.DB.Query(fmt.Sprintf("SELECT TypeID FROM CatmaAttributes WHERE catmaValueText LIKE '%%%v%%' AND catmaAttributeName == 'mDisplayName'", name))
	if err != nil {
		return values, err
	}
	for rows.Next() {
		var nTypeID int

		rows.Scan(&nTypeID)
		value := SDEType{s, nTypeID, "", make(map[string]interface{})}
		values = append(values, value)
	}
	return values, nil
}

type joint struct {
	I int
	D bool
}

// Dump attemps to dump all relevent types to a file.
func (s *SDE) Dump() error {
	fmt.Println("Begining relevant type dump")
	rows, err := s.DB.Query("SELECT TypeID FROM CatmaTypes;")
	if err != nil {
		return err
	}
	TypeIDs := make([]*joint, 0)
	for rows.Next() {
		var nTypeID int
		rows.Scan(&nTypeID)
		TypeIDs = append(TypeIDs, &joint{nTypeID, false})
	}
	fmt.Println("Collected all typeIDs.  Total of:", len(TypeIDs))
	fmt.Println("Begining filtering.  This may take awhile.")
	go func() {
		for {
			select {
			case <-time.Tick(time.Second):
				var tDone int
				for _, v := range TypeIDs {
					if v.D {
						tDone++
					}
				}
				if tDone >= len(TypeIDs) {
					break
				}
				fmt.Printf("\r%v/%v", tDone, len(TypeIDs))
			}
		}
	}()
	file, err := os.Create("out.txt")
	defer file.Close()
	if err != nil {
		return err
	}

	for _, v := range TypeIDs {
		t, _ := s.GetType(v.I)
		t.GetAttributes()
		if t.IsObtainable() {
			name := t.GetName()
			id := t.TypeID
			name = strings.Replace(name, " ", "_", -1)
			name = strings.Replace(name, "'", "_", -1)
			name = strings.Replace(name, "-", "_", -1)
			fmt.Fprintf(file, "%v := %v\n", name, id)
		}
		v.D = true
	}
	return nil
}
