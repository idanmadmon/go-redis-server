package go_redis_server

func buildRespNullBulkString() string {
	return string(bulkStringSign) + "-1\r\n"
}