package execloop

import (
	"log"
	"os"
)

type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warningf(string, ...interface{})
	Errorf(string, ...interface{})
}

func (o *Options) Debugf(f string, v ...interface{}) {
	if o.Logger == nil {
		return
	}
	o.Logger.Debugf(f, v...)
}

func (o *Options) Infof(f string, v ...interface{}) {
	if o.Logger == nil {
		return
	}
	o.Logger.Infof(f, v...)
}

func (o *Options) Warningf(f string, v ...interface{}) {
	if o.Logger == nil {
		return
	}
	o.Logger.Warningf(f, v...)
}

func (o *Options) Errorf(f string, v ...interface{}) {
	if o.Logger == nil {
		return
	}
	o.Logger.Errorf(f, v...)
}

type defaultLogger struct {
	*log.Logger
}

func (l *defaultLogger) Debugf(f string, v ...interface{}) {
	l.Printf("DEBUG: "+f, v...)
}

func (l *defaultLogger) Infof(f string, v ...interface{}) {
	l.Printf("INFO: "+f, v...)
}

func (l *defaultLogger) Warningf(f string, v ...interface{}) {
	l.Printf("WARN: "+f, v...)
}

func (l *defaultLogger) Errorf(f string, v ...interface{}) {
	l.Printf("ERROR: "+f, v...)
}

var defaultLog = &defaultLogger{log.New(os.Stdout, "execloop", log.LstdFlags)}
