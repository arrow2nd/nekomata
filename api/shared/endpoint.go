package shared

type Endpoint string

func (e Endpoint) URL(host string) string {
	return host + "/" + string(e)
}
