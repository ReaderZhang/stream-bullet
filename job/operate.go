package job

type Bullet struct {
	Ip      string `json:"ip"`
	Id      int64  `json:"id"`
	User    string `json:"user"`
	Content string `json:"content"`
	Start   string `json:"start"`
}

func Insert(b *Bullet) {
	DB.Create(&b)
}
