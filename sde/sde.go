/*
	The sde package is a fully functional library for use with the DUST514
	Static Data Export.  The package automatically can download and manage
	multiple versions of the SDE and has multiple data structures to
	manipulate data.
*/
package sde

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	_ "github.com/mattn/go-sqlite3" // Database driver
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	PrimarySDE *SDE
)

// GiveSDE is used to give the sde package your primary SDE that you've opened
// We need this for fits to pull the data from the correct database.  If you're
// not using fits don't bother.  All of the fit provider functions should warn
// if it's not set.
func GiveSDE(s *SDE) {
	PrimarySDE = s
}

// SDE is a struct containing the database object, the version of the SDE
// and many methods for working with the SDE.
type SDE struct {
	DB      *sql.DB `json:"-"`
	Version string  `json:"version"`
	Cache   bool    `json:"useCache"`
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

func (s *SDE) GetTypesWithTag(tag int) []*SDEType {
	types := make([]*SDEType, 0)
	rows, err := s.DB.Query(fmt.Sprintf("SELECT TypeID FROM CatmaAttributes WHERE catmaValueInt == '%v'", tag))
	if err != nil {
		fmt.Println(err.Error())
	}
	for rows.Next() {
		var typeID int
		rows.Scan(&typeID)
		t, _ := s.GetType(typeID)
		t.GetAttributes()
		types = append(types, &t)
	}
	return types
}

// GetType returns an SDEType of the given TypeID
func (s *SDE) GetType(id int) (SDEType, error) {
	if id <= 0 {
		return SDEType{}, errors.New("Given negative id")
	}
	defer Debug(time.Now())
	if s.Cache {
		t, d := Cache.GetType(id)
		if !d {
			t.FromCache = true
			return t, errors.New("No such type.")
		}
		log.Info("Used cache")
		return t, nil
	}

	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM CatmaTypes WHERE TypeID == '%v'", id))
	if err != nil {
		return SDEType{}, err
	}
	if rows.Next() {
		var nTypeID int
		var nTypeName string

		rows.Scan(&nTypeID, &nTypeName)
		t := SDEType{s, nTypeID, nTypeName, make(map[string]interface{}), false, false}
		t.GetAttributes()
		return t, nil
	}
	return SDEType{}, errors.New("no such type")
}

// GetType returns an SDEType of the given TypeID
func (s *SDE) GetTypeQ(id int) (SDEType, error) {
	if id <= 0 {
		return SDEType{}, errors.New("Given negative id")
	}
	defer Debug(time.Now())
	if s.Cache {
		t, d := Cache.GetType(id)
		if !d {
			t.FromCache = true
			return t, errors.New("No such type.")
		}
		log.Info("Used cache")
		return t, nil
	}

	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM CatmaTypes WHERE TypeID == '%v'", id))
	if err != nil {
		return SDEType{}, err
	}
	if rows.Next() {
		var nTypeID int
		var nTypeName string

		rows.Scan(&nTypeID, &nTypeName)
		t := SDEType{s, nTypeID, nTypeName, make(map[string]interface{}), false, false}
		return t, nil
	}
	return SDEType{}, errors.New("no such type")
}

func (s *SDE) GetTypeByName(name string) (SDEType, error) {
	defer Debug(time.Now())
	if s.Cache {
		vals := Cache.Search(name)
		if len(vals) != 0 {
			vals[0].FromCache = true
			return vals[0], nil
		}
		return SDEType{}, errors.New("No such type.")
	}

	rows, err := s.DB.Query(fmt.Sprintf("SELECT * FROM CatmaTypes WHERE TypeName == '%v'", name))
	if err != nil {
		return SDEType{}, err
	}
	if rows.Next() {
		var nTypeID int
		var nTypeName string

		rows.Scan(&nTypeID, &nTypeName)
		t := SDEType{s, nTypeID, nTypeName, make(map[string]interface{}), false, false}
		t.GetAttributes()
		return t, nil
	}
	return SDEType{}, errors.New("no such type")
}

// GetTypeWhereNameContains should be thought of as a search function that
// checks the display name.
func (s *SDE) GetTypeWhereNameContains(name string) ([]*SDEType, error) {
	log.Info(name)
	defer Debug(time.Now())
	if s.Cache {
		vals := Cache.Search(name)
		t := make([]*SDEType, 0)
		for _, v := range vals {
			v.FromCache = true
			t = append(t, &v)
		}
		return t, nil
	}

	values := make([]*SDEType, 0)
	rows, err := s.DB.Query(fmt.Sprintf("SELECT TypeID FROM CatmaAttributes WHERE catmaValueText LIKE '%%%v%%' AND catmaAttributeName == 'mDisplayName'", name))
	if err != nil {
		return values, err
	}
	for rows.Next() {
		var nTypeID int

		rows.Scan(&nTypeID)
		value := &SDEType{s, nTypeID, "", make(map[string]interface{}), false, false}
		values = append(values, value)
	}
	return values, nil
}

// Search returns a slice of *SDEType where the given string is either in the
// TypeID, TypeName or mDisplayName attribute.  It starts by checking the
// mDisplayName first, than TypeName and ID if all else fails.
func (s *SDE) Search(search string) ([]*SDEType, error) {
	defer Debug(time.Now())

	data, err := s.GetTypeWhereNameContains(search)
	if len(data) != 0 && err == nil {
		return data, nil
	}
	if err != nil {
		log.LogError("Error: ", err.Error())
	}
	log.Info("No data from GetTypeWhereNameContains")

	values := make([]*SDEType, 0)
	var err2 error
	var rows *sql.Rows
	rows, err2 = s.DB.Query(fmt.Sprintf("SELECT typeID, typeName FROM CatmaTypes WHERE typeName like '%%%v%%' ESCAPE '^'", search))

	if err2 != nil {
		log.LogError(err2.Error())
	}

	for rows.Next() {
		var (
			nTypeID   int
			nTypeName string
		)

		err := rows.Scan(&nTypeID, &nTypeName)
		if err != nil {
			log.LogError("Scan error", err.Error())
		}
		values = append(values, &SDEType{s, nTypeID, nTypeName, make(map[string]interface{}), false, false})
	}

	if len(values) != 0 {
		return values, err2
	}

	id, _ := strconv.Atoi(search)
	st, _ := s.GetType(id)
	values = append(values, &st)
	return values, nil
}

type joint struct {
	I int
	D bool
}

// Dump attemps to dump all relevent types to a file.
// Uses lots of memory.  Be careful.
func (s *SDE) Dump() error {
	defer Debug(time.Now())

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

type Progress struct {
	Percent int
	Current int
	Total   int
}

// GobDump attempts to cache all types
func (s *SDE) GobDump() chan Progress {
	defer Debug(time.Now())

	t := make(chan int)
	done := make(chan bool)
	var count int

	percent := make(chan Progress)

	go func() {
		fmt.Println("Starting caching process")
		go func() {
			r, _ := s.DB.Query("SELECT Count(*) FROM CatmaTypes")
			r.Next()

			r.Scan(&count)

			rows, err := s.DB.Query("SELECT TypeID FROM CatmaTypes;")
			if err != nil {
				log.LogError("DBErr", err.Error())
			}
			for rows.Next() {
				var nTypeID int
				rows.Scan(&nTypeID)
				t <- nTypeID
			}
			done <- true
		}()
		var i int
		var ii int
		for {
			select {
			case id := <-t:
				ii = ii + 1
				if i >= 10 {
					log.Info("Saving cache in the case of a crash.")
					SaveCache(fmt.Sprintf("%v.sde", s.Version))
					i = 0
				}
				percent <- Progress{(ii / count), ii, count}

				_, c := Cache.GetType(id)
				if !c {
					i = i + 1
				}
			case <-done:
				close(percent)
				break
			}
		}
	}()
	return percent
}

func (s *SDE) GetTypesByClassName(name string) (map[int]*SDEType, string) {
	out := make(map[int]*SDEType, 0)

	ids := make([]int, 0)

	r, _ := s.DB.Query(fmt.Sprintf("SELECT Count(*) FROM CatmaClasses WHERE className LIKE '%%%v%%'", name))
	var count int
	var done int

	r.Scan(&count)

	fmt.Printf("Found %v entries.  Working on 0/%v", count, done)

	rows, err := s.DB.Query(fmt.Sprintf("SELECT typeID, className FROM CatmaClasses WHERE className LIKE '%%%v%%'", name))

	if err != nil {
		log.LogError(err.Error())
	}

	var className string

	for rows.Next() {
		var id int
		rows.Scan(&id, &className)
		ids = append(ids, id)
	}

	for _, v := range ids {
		t, err := s.GetType(v)
		fmt.Printf("\rFound %v entries.  Working on 0/%v", count, done)
		if err != nil {
			log.LogError(err.Error())
		}
		out[v] = &t
		done++
	}
	return out, className
}

func (s *SDE) IterFunc(f func(TypeID int) bool) error {
	if rows, err := s.DB.Query("SELECT typeID FROM CatmaTypes;"); err == nil {
		for rows.Next() {
			var typeID int
			rows.Scan(&typeID)
			if !f(typeID) {
				break
			}
		}
	} else {
		return err
	}
	return nil
}
