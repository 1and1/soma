package main

import "net/url"

//var Slog *log.Logger

/*
func initLogFile() {
	f, err := os.OpenFile(path.Join(Cfg.Run.PathLogs, "somaadm.log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing logging: %s\n", err.Error())
		os.Exit(1)
	}

	utl.SetLog(log.New(f, "", log.Ldate|log.Ltime|log.LUTC))
	// XXX COMPAT
	Cfg.Run.Logger = utl.Log
	Slog = utl.Log
}
*/

func getApiUrl() *url.URL {
	url, err := url.Parse(Cfg.Api)
	if err != nil {
		utl.Log.Printf("Error parsing API address from config file")
		utl.Log.Fatal(err)
	}
	return url
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
