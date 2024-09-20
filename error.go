package routeros

import (
	"errors"
	"fmt"

	"github.com/tealok-tech/routeros/v3/proto"
)

var (
	errAlreadyAsync   = errors.New("method Async() has already been called")
	errAsyncLoopEnded = errors.New("method Async(): loop has ended - probably read error")
)

// A list of sentences we recognize from the router
const (
	MsgInvalidUserNameOrPassword = "invalid user name or password (6)"
)

// UnknownReplyError records the sentence whose Word is unknown.
type UnknownReplyError struct {
	Sentence *proto.Sentence
}

func (err *UnknownReplyError) Error() string {
	return "unknown RouterOS reply word: " + err.Sentence.Word
}

// DeviceError records the sentence containing the error received from the device.
// The sentence may have Word !trap or !fatal.
type DeviceError struct {
	Sentence *proto.Sentence
}

// Well-known errors. We could just dynamically create these out of the various
// strings we get back from RouterOS, but by listing them here we capture knowledge
// and experience recognizing the different messages. This makes it easier for clients
// to write good error-handling code.
var (
	ErrInvalidAuthentication =  errors.New(MsgInvalidUserNameOrPassword)
)

// A mapping between the recognized sentences and error types
var ErrorsByMessage = map[string]error{
	MsgInvalidUserNameOrPassword: ErrInvalidAuthentication,
}

func (err *DeviceError) fetchMessage() string {
	if m := err.Sentence.Map["message"]; m != "" {
		return m
	}

	return "unknown error: " + err.Sentence.String()
}

func (err *DeviceError) Error() string {
	return fmt.Sprintf("from RouterOS device: %s", err.fetchMessage())
}

// Given a sentence see if it's a well-known error. If it is, return that error.
// Otherwise produce a new DeviceError with the details of the sentence.
// Users are encouraged to send pull-requests with new error types to make error
// handling by type easier.
func DeviceErrorFromSentence(sen* proto.Sentence) error {
	m := sen.Map["message"]
	if err, ok := ErrorsByMessage[m]; ok {
		return err
	}
	return &DeviceError{sen}
}
