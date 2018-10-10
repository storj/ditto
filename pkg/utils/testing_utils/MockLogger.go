package testing_utils

import (
	"errors"
)

var indexOutOfBoundError = errors.New("Index out of bound")
var unitializedSliceError = errors.New("No calls were made to method")

type MockLogger struct {
	logCount, logECount int
	logParams []string
	logEParams []error
}

func (l *MockLogger) Log(msg string) {
	l.logCount++
	l.logParams = append(l.logParams, msg)
}

func (l *MockLogger) LogE(err error) {
	l.logECount++
	l.logEParams = append(l.logEParams, err)
}

func (l *MockLogger) GetLogParam(i int) (string, error) {
	if l.logParams == nil {
		return "", unitializedSliceError
	}

	length := len(l.logParams)
	if i < 0 || i > length - 1 {
		return "", indexOutOfBoundError
	}

	return l.logParams[i], nil
}

func (l *MockLogger) GetLogEParam(i int) (error, error) {
	if l.logEParams == nil {
		return nil, unitializedSliceError
	}

	length := len(l.logEParams)
	if i < 0 || i > length - 1 {
		return nil, indexOutOfBoundError
	}

	return l.logEParams[i], nil
}

func (l *MockLogger) GetLastLogParam() (string, error) {
	return l.GetLogParam(len(l.logParams) - 1)
}

func (l *MockLogger) GetLastLogEParam() (error, error) {
	return l.GetLogEParam(len(l.logEParams) - 1)
}

func (l *MockLogger) LogCount() int {
	return l.logCount
}

func (l *MockLogger) LogECount() int {
	return l.logECount
}

func (l *MockLogger) Clear() {
	l.logEParams = []error{}
	l.logParams = []string{}
}



