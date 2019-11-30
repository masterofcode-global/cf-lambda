package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		event  events.SQSEvent
		expect string
		err    error
	}{{
		event: events.SQSEvent{
			Records: []events.SQSMessage{{
				MessageId: "1",
				Body:      "Test message",
			},
			},
		},
		expect: "Procesed messages",
		err:    nil,
	},
	}

	for _, test := range tests {
		response, err := Handler(nil, test.event)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response)
	}

}
