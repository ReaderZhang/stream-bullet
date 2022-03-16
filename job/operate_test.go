package job

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	InitDB()
	bullet := &Bullet{
		Ip:      "123",
		Id:      1,
		User:    "qqz",
		Start:   "2022/03/15 10:26",
		Content: "abc",
	}
	str, _ := json.Marshal(bullet)
	fmt.Printf(string(str))
	Insert(bullet)
	//DB.Find(&b, "id=?", 1)
}
