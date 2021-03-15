package eqm

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type QMgr struct {
	root string
	mtx  sync.Mutex
}

/*
Initializes the queue - eventually creating the dir
*/
func (q *QMgr) Init() error {
	q.root = "./db"
	_, err := os.Stat(q.root)
	if err != nil {
		os.MkdirAll(q.root, os.ModePerm)
	}
	return nil
}

/*
Pushes new message into the queue
*/
func (q *QMgr) Push(b []byte) error {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	now := time.Now()
	fname := fmt.Sprintf("%d.%d", now.Unix(), now.UnixNano())
	return os.WriteFile(path.Join(q.root, fname), b, 0600)
}

/*
Checks the 1st message available without removing from queue
*/
func (q *QMgr) Peek() ([]byte, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	found := false
	var ret []byte
	filepath.Walk(q.root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			ret, err = os.ReadFile(path)
			if err == nil {
				found = true
				return errors.New("FOUND")
			}
		}
		return nil
	})
	if !found {
		return ret, errors.New("Queue Empty")
	}
	return ret, nil
}

/*
Removes first message from queue
*/
func (q *QMgr) Pop() ([]byte, error) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	found := false
	var ret []byte
	filepath.Walk(q.root, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			ret, err = os.ReadFile(path)
			if err == nil {
				found = true
				os.Remove(path)
				return errors.New("FOUND")
			}
		}
		return nil
	})
	if !found {
		return ret, errors.New("Queue Empty")
	}
	return ret, nil
}

/*
Clears all messages in the queue
*/
func (q *QMgr) Clear() error {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	err := filepath.Walk(q.root, func(path string, info fs.FileInfo, err error) error {
		os.Remove(path)
		return nil
	})
	return err
}

/*
Creates a new QMgr instance
*/
func New() (*QMgr, error) {
	var ret QMgr
	err := ret.Init()
	return &ret, err
}
