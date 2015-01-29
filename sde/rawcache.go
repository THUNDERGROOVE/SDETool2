package sde

import (
	"encoding/gob"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"os"
	"strings"
	"time"
	"unsafe"
)

// Raw cache tries caching things in using gobs

type SDECache struct {
	Types   map[int]*SDEType
	Version string
}

func (s *SDECache) GetType(id int) (SDEType, bool) {
	defer Debug(time.Now())
	if v, ok := s.Types[id]; ok {
		return *v, true
	}
	log.Info("Type not in cache.  Adding.")
	t, _ := PrimarySDE.GetType(id)
	s.Types[t.TypeID] = &t
	return t, false
}

func (s *SDECache) Search(name string) []SDEType {
	values := make([]SDEType, 0)
	for _, v := range s.Types {
		if strings.Contains(v.GetName(), name) {
			values = append(values, *v)
			continue
		}
		if strings.Contains(v.TypeName, name) {
			values = append(values, *v)
			continue
		}
	}
	return values
}

var Cache SDECache

func init() {
	gob.Register(SDE{})
	gob.Register(SDEType{})
	gob.Register(SDECache{})
	Cache = SDECache{make(map[int]*SDEType, 0), "uninitialzed"}
}

func SaveCache(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}

	enc := gob.NewEncoder(f)
	eerr := enc.Encode(&Cache)
	if eerr != nil {
		log.LogError("Saving to cache file", eerr.Error(), filename)
		return eerr
	}
	return nil
}

func LoadCache(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	dec := gob.NewDecoder(f)
	derr := dec.Decode(&Cache)
	if derr != nil {
		log.LogError("Error reading cache file", derr.Error(), filename)
		return derr
	}
	log.Info("Cache successfully loaded with a size of", unsafe.Sizeof(&Cache), "bytes")
	return nil
}
