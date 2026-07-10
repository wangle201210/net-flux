package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	checkInterval = 2 * time.Minute
)

var (
	// make filePipeline buffered to avoid deadlock when Open() is called
	// from the same goroutine that reads from the channel (see _startTimer).
	filePipeline chan *RollingFile
)

type RollingFile struct {
	sync.Mutex
	fileName    string
	checkFlag   string
	fileHandler *os.File
	checkFun    func() string
}

// injectPidIntoLogPath 在文件名中插入 _pid_<pid>，得到 filename_pid_xxx.log 形式
func injectPidIntoLogPath(path string) string {
	dir, base := filepath.Dir(path), filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	pid := os.Getenv("HOSTNAME") //os.Getpid()
	newBase := name + "_" + pid + ext
	return filepath.Join(dir, newBase)
}

func NewRollingFile(fileName string) *RollingFile {
	fileName = injectPidIntoLogPath(fileName)
	rf := &RollingFile{
		fileName: fileName,
		checkFun: func() string {
			return time.Now().Format("20060102")
		},
	}

	return rf
}

func (rf *RollingFile) SetCheckFun(fn func() string) {
	rf.checkFun = fn
}

func (rf *RollingFile) Open() error {

	rf.Lock()
	defer rf.Unlock()
	if rf.fileHandler != nil {
		_ = rf.fileHandler.Close()
		rf.fileHandler = nil
	}
	rf.checkFlag = rf.checkFun()
	fileName := fmt.Sprintf("%s.%s", rf.fileName, rf.checkFlag)

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open file error:%s", fileName)
		return err
	}
	rf.fileHandler = f

	// try to register this RollingFile with the background timer goroutine
	// use a non-blocking send to avoid potential deadlock when the
	// background goroutine is currently handling a tick and not able to
	// receive from the channel.
	select {
	case filePipeline <- rf:
	default:
		// drop registration if pipeline is full; timer will pick it up
		// on next Open call or it has already been registered
	}
	return nil
}

func (rf *RollingFile) Write(p []byte) (n int, err error) {
	if rf.fileHandler != nil {
		rf.Lock()
		defer rf.Unlock()
		return rf.fileHandler.Write(p)
	}
	return 0, os.ErrClosed
}

func (rf *RollingFile) RedirectSTDOutput() error {
	if rf.fileHandler != nil {
		if err := dup2(int(rf.fileHandler.Fd()), 1); err != nil {
			return err
		}
		if err := dup2(int(rf.fileHandler.Fd()), 2); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("RedirectSTDOutput:Open file first")
	}
}

func _startTimer() {
	tick := time.Tick(checkInterval)
	go func() {
		fileList := make(map[*RollingFile]interface{})
		for {
			select {
			case f := <-filePipeline:
				fileList[f] = new(interface{})
			case <-tick:
				for k, _ := range fileList {
					newName := k.checkFun()
					if !strings.EqualFold(newName, k.fileName) {
						k.checkFlag = newName
						delete(fileList, k)
						if err := k.Open(); err != nil {
							fmt.Printf("RollingFile:Open file failed:%v", err)
						}
					}
				}
			}
		}
	}()
}

func init() {
	// small buffer to decouple producers from the timer consumer and
	// avoid self-deadlock when the timer goroutine calls Open()
	filePipeline = make(chan *RollingFile, 16)
	_startTimer()
}
