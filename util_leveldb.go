package main

import (
	"bytes"
	"encoding/gob"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func jobDbOpen() {
	// we already have an open LevelDB
	if Cfg.Run.LevelDB != nil {
		return
	}
	var ldbOpt opt.Options
	ldbOpt.ErrorIfMissing = true

	if Cfg.Run.PathLevelDB == "" {
		Slog.Fatal("No path information to JobDB available")
	}

	var err error
	Cfg.Run.LevelDB, err = leveldb.OpenFile(Cfg.Run.PathLevelDB, &ldbOpt)
	if err != nil {
		Slog.Fatal(err)
	}
}

func jobDbGetOutstandingJobs() []string {
	jobDbOpen()
	result := make([]string, 0)

	data, err := Cfg.Run.LevelDB.Get([]byte("oustanding_jobs"), nil)
	if err != nil {
		// key does not exist, return empty slice
		return result
	}

	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&result)
	if err != nil {
		Slog.Fatal("Error decoding outstanding jobs: ", err)
	}
	return result
}

func jobDbPutOutstandingJobs(s []string) {
	jobDbOpen()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		Slog.Fatal("Error encoding outstanding jobs: ", err)
	}

	err = Cfg.Run.LevelDB.Put([]byte("oustanding_jobs"), buf.Bytes(), nil)
	if err != nil {
		Slog.Fatal("LevelDB.Put(outstanding_jobs): ", err)
	}
}

func jobDbAddOutstandingJob(s string) {
	jobs := jobDbGetOutstandingJobs()
	jobs = append(jobs, s)
	jobDbPutOutstandingJobs(jobs)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
