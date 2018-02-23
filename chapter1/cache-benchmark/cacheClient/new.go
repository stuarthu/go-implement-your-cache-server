package cacheClient

func NewCacheClient(typ, server string) Client {
	if typ == "redis" {
		return newRedisClient(server)
	}
	if typ == "http" {
		return newHTTPClient(server)
	}
	panic("unknown client type " + typ)
}
