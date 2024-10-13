package define

type FetchResponse struct {
	Status   int        `json:"status"`
	Headers  [][]string `json:"headers"`
	Text     string     `json:"text"`
	FinalUrl string     `json:"final_url"`
}
type FetchResp struct {
	Response FetchResponse `json:"response"`
	Error    string        `json:"error_stack"`
}
type RelayCommandResp struct {
	CommandId     string    `json:"command_id"`
	CommandResult FetchResp `json:"command_result"`
}
