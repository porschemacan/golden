package libs

import (
	"fmt"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/lestrrat-go/file-rotatelogs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	program   = filepath.Base(os.Args[0])
	didi_log  = log.New()
	info_log  = log.New()
)
var once sync.Once

func InitLog(config *LogConfig) {
	once.Do(func() {
		didi_log.SetReportCaller(true)
		info_log.SetReportCaller(true)
		didi_log.AddHook(newLfsHook("didi_log."+program+".log", config))
		info_log.AddHook(newLfsHook(program+".log", config))
		didi_log.Out = ioutil.Discard
		info_log.Out = ioutil.Discard
	})
}

func newLfsHook(filename string, config *LogConfig) log.Hook {
	writer, err := rotatelogs.New(
		config.Path+filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(config.Path+filename),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(3*time.Hour),
	)

	if err != nil {
		writer, err = rotatelogs.New(
			"/tmp/"+filename+".%Y%m%d%H",
			rotatelogs.WithLinkName("/tmp/"+filename),
			rotatelogs.WithRotationTime(time.Hour),
			rotatelogs.WithMaxAge(3*time.Hour),
		)
	}

	if err != nil {
		panic(fmt.Sprintf("init log fail, config:%+v", config))
	}

	log.SetLevel(log.Level(config.Level))

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.TextFormatter{
		DisableColors:          true,
		QuoteEmptyFields:       false,
		DisableLevelTruncation: false,
		TimestampFormat:        "2006-01-02T15:04:05.000Z0700",
		FieldMap: log.FieldMap{
			log.FieldKeyTime:  "timestamp",
			log.FieldKeyLevel: "level",
			log.FieldKeyMsg:   "_msg"},
	})

	return lfsHook
}

// didi log
func GetDLog(tag string) *log.Entry {
	return didi_log.WithField("action",tag)
}

func DLogf(format string, args ...interface{}) {
	didi_log.Infof(format, args...)
}

func DTagf(tag string, format string, args ...interface{}) {
	didi_log.WithField("action",tag).Infof(format, args...)
}

// info log
func Infof(format string, args ...interface{}) {
	info_log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	info_log.Warnf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	info_log.Fatalf(format, args...)
}
