package lock

import "os"

var locks *os.File
var lockFile = "./logs/BuilderXLock.pid"
