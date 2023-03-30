package info

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

/* Base Info Variables */

// A variable that indicates if service is in test mode
var TestMode bool

// A variable that indicates if service should virtualize external integration calls
var Virtualize bool

// A variable that contains the running diary instance
var Diary diary.IDiary

// A variable that contains the running diary instance's sample trace rate value
var TraceRate int

// A variable that contains the running uniform event bus connection instance
var Conn uniform.IConn
