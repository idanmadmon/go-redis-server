package go_redis_server

import "strconv"

func buildRespSimpleString(str string) string {
	return string(simpleStringSign) + str + "\r\n"
}

func buildRespError(err error) string {
	return string(errorSign) + err.Error() + "\r\n"
}

func buildRespInteger(i int64) string {
	return string(integerSign) + string(i) + "\r\n"
}

func buildRespBulkString(str string) string {
	return string(bulkStringSign) + strconv.Itoa(len(str)) + "\r\n" + str + "\r\n"
}

func buildRespNullBulkString() string {
	return string(bulkStringSign) + "-1\r\n"
}