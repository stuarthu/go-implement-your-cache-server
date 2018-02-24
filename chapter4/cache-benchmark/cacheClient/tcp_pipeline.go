package cacheClient

func (r *tcpClient) PipelinedRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	for _, c := range cmds {
		if c.Name == "get" {
			r.sendGet(c.Key)
		}
		if c.Name == "set" {
			r.sendSet(c.Key, c.Value)
		}
		if c.Name == "del" {
			r.sendDel(c.Key)
		}
	}
	for _, c := range cmds {
		c.Value, c.Error = r.recvResponse()
	}
}
