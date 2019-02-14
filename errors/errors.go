package errors

import (
	"fmt"
	"sync"
)

// Errors is used to collect errors
type Errors struct {
	errors []error
	lock   *sync.Mutex
}

//Make new Errors object
func NewErrors() Errors {
	return Errors{
		lock: &sync.Mutex{},
	}
}

//Add adds an error if not nil
func (e *Errors) Add(err error) bool {
	if err != nil {
		e.lock.Lock()
		defer e.lock.Unlock()
		e.errors = append(e.errors, err)
		return true
	}
	return false
}

func (e *Errors) Error() string {
	str := ""
	for _, err := range e.errors {
		str = fmt.Sprintf("%s\n%s", str, err.Error())
	}
	return str
}

//Get an error if any have been recorded
func (e *Errors) Get() error {
	if len(e.errors) > 0 {
		return e
	}
	return nil
}
