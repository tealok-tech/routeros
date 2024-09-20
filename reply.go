package routeros

import (
	"strings"

	"github.com/tealok-tech/routeros/v3/proto"
)

// Reply has all the sentences from a reply.
type Reply struct {
	Re   []*proto.Sentence
	Done *proto.Sentence
}

func (r *Reply) String() string {
	var sb strings.Builder
	for _, sen := range r.Re {
		sb.WriteString(sen.String())
		sb.WriteRune('\n')
	}

	if r.Done != nil {
		sb.WriteString(r.Done.String())
	}

	return sb.String()
}

// Return whether or not we are done processing, and any error detected in the sentence
func (r *Reply) processSentence(sen *proto.Sentence) (bool, error) {
	switch sen.Word {
	case reSentence:
		r.Re = append(r.Re, sen)
	case doneSentence:
		r.Done = sen
		return true, nil
	case trapSentence, fatalSentence:
		return sen.Word == fatalSentence, DeviceErrorFromSentence(sen)
	case "":
		// API docs say that empty sentences should be ignored
	default:
		return true, &UnknownReplyError{sen}
	}
	return false, nil
}
