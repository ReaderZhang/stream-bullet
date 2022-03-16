package consumer

import (
	"stream-bullet/job"
	"testing"
)

func TestConsumer(t *testing.T) {
	job.InitDB()
	consumer()
}
