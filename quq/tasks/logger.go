package tasks

import "log"

var Logger *TaskLogger

func init() {
	Logger = &TaskLogger{logger: log.Default()}
}

type TaskLogger struct {
	logger *log.Logger
}

func (tl *TaskLogger) Debug(args ...interface{}) {
	tl.logger.Println(args...)
}

func (tl *TaskLogger) Info(args ...interface{}) {
	tl.logger.Println(args...)
}

func (tl *TaskLogger) Warn(args ...interface{}) {
	tl.logger.Println(args...)
}

func (tl *TaskLogger) Error(args ...interface{}) {
	tl.logger.Println(args...)
}

func (tl *TaskLogger) Fatal(args ...interface{}) {
	tl.logger.Fatal(args...)
}
