package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// LogFileMaxSize max size per log file: 100M
const LogFileMaxSize = 100

// RotateWriter struct define
type RotateWriter struct {
	lock        sync.Mutex
	filename    string
	maxSize     int64
	rotateDaily bool
	currDate    string
	fp          *os.File
	quit        chan int
}

// NewRotateWriter make a new RotateWriter. Return nil if error occurs during setup.
func NewRotateWriter(filename string, maxSize int64, rotateDaily bool) (*RotateWriter, error) {
	w := &RotateWriter{filename: filename, maxSize: maxSize, rotateDaily: rotateDaily}
	err := w.rotate()
	if err != nil {
		return nil, err
	}
	if w.maxSize == 0 {
		w.maxSize = LogFileMaxSize
	}
	w.currDate = time.Now().Format("2006-01-02")
	w.quit = make(chan int)
	go w.autoRotate(w.quit)
	return w, nil
}

// Close ...
func (w *RotateWriter) Close() error {
	if w.quit != nil {
		close(w.quit)
	}
	if w.fp != nil {
		return w.fp.Close()
	}
	return nil
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(output []byte) (int, error) {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.fp.Write(output)
}

// Perform the actual act of rotating and reopening file.
func (w *RotateWriter) rotate() (err error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// Close existing file if open
	if w.fp != nil {
		err = w.fp.Close()
		w.fp = nil
		if err != nil {
			return
		}
	}

	// Rename dest file if it already exists
	fileinfo, err := os.Stat(w.filename)
	if err == nil {
		if fileinfo.Size() > 0 {
			backupFilename := w.filename + "." + time.Now().Format("2006-01-02_15:04:05")
			err = os.Rename(w.filename, backupFilename)
			if err != nil {
				return
			}
		}
	}

	// Create a file.
	w.fp, err = os.Create(w.filename)
	return
}

func (w *RotateWriter) autoRotate(quit chan int) {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-quit:
			fmt.Printf("quit auto rotate log file\n")
			return
		case <-ticker.C:

			//check log file size
			fileinfo, err := os.Stat(w.filename)
			if err == nil {
				if fileinfo.Size() >= w.maxSize*1024*1024 {
					//rotate log file
					fmt.Printf("start to rotate log file...\n")
					err = w.rotate()
					if err != nil {
						fmt.Printf("rotate log file fail: %v\n", err)
					}
					continue
				}
			}

			//check date
			if w.rotateDaily {
				date := time.Now().Format("2006-01-02")
				if date != w.currDate {
					if fileinfo.Size() > 0 {
						//rotate log file
						fmt.Printf("start to rotate log file...\n")
						err = w.rotate()
						if err != nil {
							fmt.Printf("rotate log file fail: %v\n", err)
						} else {
							w.currDate = date
						}
						continue
					}
				}
			}
		}
	}
}
