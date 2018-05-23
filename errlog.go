package main

import (
    "os"
	"log"
	"time"
)


const LOGDIR = "./log"

var logger *Logger

type Logger struct{
    dir string
	lr  *log.Logger
}

func NewLogger() (logger *Logger) {
	if _, err := os.Stat(LOGDIR); (err != nil) && os.IsNotExist(err) {
		if err = os.Mkdir(LOGDIR, os.ModePerm); err != nil {
			log.Fatal("failed to create log directory .")
		}
	}
	filename := time.Now().Format("2006-01-02") + ".log"
	file := LOGDIR + "/" + filename
	var f *os.File
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err = os.Create(file)
		if err != nil {
		    log.Fatalln("fail to create log file!")
	    }
	} else {
		f, err = os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
		    log.Fatalln("fail to open log file!")
	    }
	}
	
	logger = &Logger{dir:LOGDIR + "/" + filename, lr:log.New(f, "[sys]", log.Ltime)}
	return 
}

//print log and exit
func (l *Logger)Fatal(msg string){
	log.Println(msg)
	l.lr.Fatal(msg)
}

func (l *Logger)Print(msg string){
	log.Println(msg)
	l.lr.Println(msg)
}
