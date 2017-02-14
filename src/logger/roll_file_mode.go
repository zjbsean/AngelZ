package logger

//定义LOG文件滚动模式

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
	"util"
)

var (
	fileObj *fileStruct  //LOG文件管理对象
	ticker  *time.Ticker //ticker计时器
)

//文件大小

const (
	B int64 = 1 << (iota * 10)
	KB
	MB
	GB
	TB
)

//文件滚动模式
type RollMode int32

const (
	//Daily : 每日
	Daily RollMode = iota + 1
	//Hour : 每小时
	Hour
	//Size : 文件大小
	Size
)

type fileStruct struct {
	dir              string
	fileName         string
	suffix           string
	nextRollTimeUnix int64
	mu               *sync.RWMutex
	logFile          *os.File
	lg               *log.Logger

	rollMode rollCheckInterface
}

func (pfile *fileStruct) setRollData(dir, fileName string, rollMode rollCheckInterface) {
	fileObj.dir = dir
	fileObj.fileName = dir + "/" + fileName
	fileObj.suffix = ""

	fileObj.rollMode = rollMode
}

func (pfile *fileStruct) chackAndRename(t *time.Time) {
	if pfile.rollMode.isMustRename(t) == true {
		pfile.rollMode.rename(t)
		pfile.rollMode.setTimeSuffixAndNextRename(t)
	}
}

type rollCheckInterface interface {
	isMustRename(t *time.Time) bool
	rename(t *time.Time)
	setTimeSuffixAndNextRename(t *time.Time)
}

//dailyRoll:
//  每日滚动
type dailyRoll struct {
	file *fileStruct
}

func (roll *dailyRoll) isMustRename(t *time.Time) bool {
	tUnix := t.Unix()
	if tUnix >= roll.file.nextRollTimeUnix {
		return true
	}
	return false
}

func (roll *dailyRoll) rename(t *time.Time) {
	fn := roll.file.fileName + "." + roll.file.suffix
	if isExist(fn) {
		os.Remove(fn)
	}

	closeLogFile()

	if isExist(roll.file.fileName) {
		err := os.Rename(roll.file.fileName, fn)
		if err != nil {
			roll.file.lg.Println("rename err", err.Error())
		}
	}

	roll.file.logFile, _ = os.Create(roll.file.fileName)
	roll.file.lg = log.New(roll.file.logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func (roll *dailyRoll) setTimeSuffixAndNextRename(t *time.Time) {
	dayStartTime := util.GetTimeDayStartTime(t)
	roll.file.suffix = dayStartTime.Format(util.DateFormat)
	roll.file.nextRollTimeUnix = dayStartTime.Unix() + 24*3600
}

//hourRoll:
//  每小时滚动
type hourRoll struct {
	file *fileStruct
}

func (roll *hourRoll) isMustRename(t *time.Time) bool {
	tUnix := t.Unix()
	if tUnix >= roll.file.nextRollTimeUnix {
		return true
	}
	return false
}

func (roll *hourRoll) rename(t *time.Time) {
	fn := roll.file.fileName + "." + roll.file.suffix
	if isExist(fn) {
		os.Remove(fn)
	}

	closeLogFile()

	if isExist(roll.file.fileName) {
		err := os.Rename(roll.file.fileName, fn)
		if err != nil {
			roll.file.lg.Println("rename err", err.Error())
		}
	}
	roll.file.logFile, _ = os.Create(roll.file.fileName)
	roll.file.lg = log.New(roll.file.logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func (roll *hourRoll) setTimeSuffixAndNextRename(t *time.Time) {
	hourStartTime := util.GetTimeHourStartTime(t)
	roll.file.suffix = hourStartTime.Format(util.DateTimeHFormat)
	roll.file.nextRollTimeUnix = hourStartTime.Unix() + 3600
}

//sizeRool:
//  文件大小滚动
type sizeRoll struct {
	file         *fileStruct
	maxSize      int64
	maxFileCount int
}

func (roll *sizeRoll) isMustRename(t *time.Time) bool {
	fileInfo, err := os.Stat(roll.file.fileName)
	if err != nil {
		panic("get file stat fail : FileName=" + roll.file.fileName)
	}

	if fileInfo.Size() >= roll.maxSize {
		return true
	}
	return false
}

func (roll *sizeRoll) rename(t *time.Time) {
	fn := roll.file.fileName + "." + roll.file.suffix
	if isExist(fn) {
		os.Remove(fn)
	}

	closeLogFile()

	if isExist(roll.file.fileName) {
		err := os.Rename(roll.file.fileName, fn)
		if err != nil {
			roll.file.lg.Println("rename err", err.Error())
		}
	}

	roll.file.logFile, _ = os.Create(roll.file.fileName)
	roll.file.lg = log.New(roll.file.logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

func (roll *sizeRoll) initSuffix() {
	suffix := 0
	for ; suffix < roll.maxFileCount; suffix++ {
		fileName := roll.file.fileName + "." + strconv.Itoa(suffix)
		if isExist(fileName) == false {
			break
		}
	}
	if suffix == roll.maxFileCount {
		suffix = 0
	}

	roll.file.suffix = strconv.Itoa(suffix)
}

func (roll *sizeRoll) setTimeSuffixAndNextRename(t *time.Time) {
	suffix, _ := strconv.Atoi(roll.file.suffix)
	suffix = (suffix + 1) % roll.maxFileCount
	roll.file.suffix = strconv.Itoa(suffix)
}

//文件监控
func fileMonitor() {
	if ticker == nil {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case curTime := <-ticker.C:
				if fileObj != nil {
					fileObj.chackAndRename(&curTime)
				}
			}
		}
	}
}

//SetDailyRollFile : 设置每日滚动模式
func SetDailyRollFile(dir, fileName string) {
	err := dirCheckAndInitFileObj(dir)
	if err != nil {
		return
	}

	fileObj.mu.Lock()
	defer fileObj.mu.Unlock()

	curTime := time.Now()
	if fileObj.rollMode == nil {
		droll := dailyRoll{fileObj}
		fileObj.setRollData(dir, fileName, &droll)
		logFileTime := getLogFileDate(fileObj.fileName)
		droll.setTimeSuffixAndNextRename(logFileTime)

		if fileObj.rollMode.isMustRename(&curTime) {
			fileObj.rollMode.rename(&curTime)
		} else {
			fileObj.logFile, _ = os.OpenFile(fileObj.fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			fileObj.lg = log.New(fileObj.logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		}

	} else {
		if _, ok := fileObj.rollMode.(*dailyRoll); ok == false {

			closeLogFile()

			droll := dailyRoll{fileObj}
			fileObj.setRollData(dir, fileName, &droll)
			droll.setTimeSuffixAndNextRename(&curTime)
			fileObj.rollMode.rename(&curTime)
		}
	}

	//go fileMonitor()
}

//SetHourRollFile : 设置每小时滚动模式
func SetHourRollFile(dir, fileName string) {
	err := dirCheckAndInitFileObj(dir)
	if err != nil {
		return
	}

	fileObj.mu.Lock()
	defer fileObj.mu.Unlock()

	curTime := time.Now()
	if fileObj.rollMode == nil {
		droll := hourRoll{fileObj}
		fileObj.setRollData(dir, fileName, &droll)
		logFileTime := getLogFileDateTimeHour(fileObj.fileName)
		droll.setTimeSuffixAndNextRename(logFileTime)
		if fileObj.rollMode.isMustRename(&curTime) {
			fileObj.rollMode.rename(&curTime)
		} else {
			fileObj.logFile, _ = os.OpenFile(fileObj.fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			fileObj.lg = log.New(fileObj.logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		}
	} else {
		if _, ok := fileObj.rollMode.(*dailyRoll); ok == false ||
			fileObj.dir != dir ||
			fileObj.fileName != fileName {

			closeLogFile()

			droll := dailyRoll{fileObj}
			droll.setTimeSuffixAndNextRename(&curTime)
			fileObj.setRollData(dir, fileName, &droll)
			fileObj.rollMode.rename(&curTime)
		}
	}

	//go fileMonitor()
}

//SetSizeRollFile : 设置文件大小滚动模式
func SetSizeRollFile(dir, fileName string, maxSize int64, maxFileCount int) {
	err := dirCheckAndInitFileObj(dir)
	if err != nil {
		return
	}

	if maxFileCount <= 0 {
		maxFileCount = 1000
	}

	fileObj.mu.Lock()
	defer fileObj.mu.Unlock()

	curTime := time.Now()
	if fileObj.rollMode == nil {
		droll := sizeRoll{fileObj, maxSize, maxFileCount}
		fileObj.setRollData(dir, fileName, &droll)
		droll.initSuffix()
		if isExist(fileObj.fileName) && fileObj.rollMode.isMustRename(&curTime) {
			fileObj.rollMode.rename(&curTime)
		} else {
			fileObj.logFile, _ = os.OpenFile(fileObj.fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
			fileObj.lg = log.New(fileObj.logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
		}
	} else {
		if _, ok := fileObj.rollMode.(*dailyRoll); ok == false ||
			fileObj.dir != dir ||
			fileObj.fileName != fileName {

			closeLogFile()

			droll := sizeRoll{fileObj, maxSize, maxFileCount}
			droll.initSuffix()
			fileObj.setRollData(dir, fileName, &droll)
			fileObj.rollMode.rename(&curTime)
		}
	}

	//go fileMonitor()
}

func fileSize(file string) int64 {
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func newFileWithBaseData() *fileStruct {
	pfile := new(fileStruct)
	pfile.mu = new(sync.RWMutex)
	pfile.nextRollTimeUnix = 0
	return pfile
}

func dirCheckAndInitFileObj(dir string) error {
	err := makeDir(dir)
	if err != nil {
		return err
	}

	if fileObj == nil {
		fileObj = newFileWithBaseData()
	}
	return nil
}

func closeLogFile() {
	if fileObj.logFile != nil {
		fileObj.logFile.Close()
		fileObj.logFile = nil
		fileObj.lg = nil
	}
}

func makeDir(dir string) error {
	_, err := os.Stat(dir)
	isExitDir := err == nil || os.IsExist(err)
	if isExitDir == false {
		if err = os.MkdirAll(dir, 0666); err != nil {
			if os.IsPermission(err) {
				fmt.Println("Create dir error : ", err.Error())
				return err
			}
		}
	}
	return nil
}

func getLogFileDate(filename string) *time.Time {
	timeFormat := "2006/01/02"
	curTime := time.Now()
	ts := curTime.Format(timeFormat)
	if isExist(filename) {
		curLogFile, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			log.Println("open file fail : err=", err.Error())
			panic(err)
		}
		defer curLogFile.Close()

		if err == nil {
			var fileContext = make([]byte, 30)
			l, err := curLogFile.Read(fileContext)
			if err == io.EOF || l == 0 {
				fmt.Println("empty file")
			}

			index := bytes.IndexByte(fileContext, ' ')
			ts = string(fileContext[:index])
		}
	}

	t, _ := time.ParseInLocation(timeFormat, ts, curTime.Location())
	return &t
}

func getLogFileDateTimeHour(filename string) *time.Time {

	timeFormat := "2006/01/02 15"
	curTime := time.Now()
	ts := curTime.Format(timeFormat)
	if isExist(filename) {
		curLogFile, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			log.Println("open file fail : err=", err.Error())
			panic(err)
		}
		defer curLogFile.Close()

		if err == nil {
			var fileContext = make([]byte, 30)
			l, err := curLogFile.Read(fileContext)
			if err == io.EOF || l == 0 {
				fmt.Println("empty file")
				return &curTime
			}

			index := bytes.IndexByte(fileContext, ':')
			ts = string(fileContext[:index])
		}
	}

	t, _ := time.ParseInLocation(timeFormat, ts, curTime.Location())
	return &t
}
