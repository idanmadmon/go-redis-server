package go_redis_server

func buildRespSimpleString(str string) string {
	return "+" + str + "\r\n"
}

func buildRespError(err error) string {
	return "-" + err.Error() + "\r\n"
}

func buildRespInteger(i int64) string {
	return ":" + string(i) + "\r\n"
}

func buildRespBulkString(str string) string {
	return "$" + string(len(str)) + "\r\n" + str + "\r\n"
}

func buildRespNullBulkString() string {
	return "$-1\r\n"
}