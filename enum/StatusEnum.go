package enum

type Status struct {
	StatusType string
}

var(
	NotFinish = Status{"未完成"}
	Doing = Status{"正在进行中"}
	Finish = Status{"已完成"}
	Cancel = Status{"已取消"}
)