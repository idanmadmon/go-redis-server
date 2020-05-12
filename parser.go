package go_redis_server

import (
	"errors"
	"strconv"
	"strings"
)

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

func parseRequest(r string) ([]string, error) {
	//TODO: convert to regex for performance (check first)
	/*
	*1\r\n
	$4\r\n
	ping\r\n

	if r == "" {
		return nil, errors.New("empty request")
	}

	parts := strings.Split(r, "\r\n")

	if r[0] != '*' {
		return nil, errors.New("bad index")
	}

	strconv.Atoi()*/

}
