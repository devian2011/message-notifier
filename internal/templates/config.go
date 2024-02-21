package templates

type TemplateCfg struct {
	Code    string   `json:"code" yaml:"code"`
	From    string   `json:"from"`
	To      []string `json:"to" yaml:"to"`
	Subject string   `json:"subject" yaml:"subject"`
	Body    struct {
		Text string `json:"text" yaml:"text"`
		Path string `json:"path" yaml:"path"`
	} `json:"body" yaml:"body"`
}
