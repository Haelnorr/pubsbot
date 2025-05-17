package main

import "flag"

func setupFlags() map[string]string {
	loglevel := flag.String("loglevel", "", "Set log level")
	logoutput := flag.String("logoutput", "", "Set log destination (file, console or both)")
	flag.Parse()
	args := map[string]string{
		"loglevel":  *loglevel,
		"logoutput": *logoutput,
	}
	return args
}
