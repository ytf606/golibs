package ginx

var (
	defaultHttpStatus int
	isMsgKey          bool
)

func SetDefaultHttpStatus(httpStatus int) {
	defaultHttpStatus = httpStatus
}

func SetMsgKey() {
	isMsgKey = true
}
