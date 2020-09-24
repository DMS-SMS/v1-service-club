package test

import (
	"club/db"
	"sync"
)

var (
	manager db.AccessorManage
	testGroup sync.WaitGroup
)

const numberOfTestFunc = 1
