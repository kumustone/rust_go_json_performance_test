package benchmark

//easyjson:json
type HttpRequest struct {
	Mark          string
	Method        string
	Scheme        string
	Url           string
	Proto         string
	Host          string
	RemoteAddr    string
	ContentLength uint64
	Header        map[string][]string
	Body          []byte
}
