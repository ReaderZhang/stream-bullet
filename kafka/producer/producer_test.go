package producer

import (
	"stream-bullet/job"
	"testing"
)

func TestProducer(t *testing.T) {
	job.InitDB()
	bullet := &job.Bullet{
		Ip:      "123",
		Id:      1,
		User:    "www",
		Start:   "2022/03/15 10:26",
		Content: "abc",
	}
	produer(bullet)
}
