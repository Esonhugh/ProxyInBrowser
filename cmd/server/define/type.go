package define

type RelayChan chan RelayCommandResp

func NewChannels() RelayChan {
	return make(chan RelayCommandResp, 1000)
}
