package restype

import "github.com/nexgus/restype/pbuf"

type ReqContext struct {
	Index   int32
	Method  string
	Headers map[string][]string
	Params  Params
	Queries map[string][]string
	Body    []byte
	Remote  string
}

func (rc *ReqContext) Header(key string) string {
	if values, ok := rc.Headers[key]; !ok {
		return ""
	} else if len(values) > 0 {
		return values[0]
	} else {
		return ""
	}
}

func (rc *ReqContext) Query(key string) string {
	if values, ok := rc.Queries[key]; !ok {
		return ""
	} else if len(values) > 0 {
		return values[0]
	} else {
		return ""
	}
}

func (rc *ReqContext) Param(key string) string {
	return rc.Params.ByName(key)
}

func (rc *ReqContext) ToPbuf() *pbuf.ReqContext {
	req := pbuf.ReqContext{
		Index:  rc.Index,
		Body:   rc.Body,
		Remote: rc.Remote,
	}

	switch rc.Method {
	case "GET":
		req.Method = pbuf.ReqContext_GET
	case "POST":
		req.Method = pbuf.ReqContext_POST
	case "PUT":
		req.Method = pbuf.ReqContext_PUT
	case "DELETE":
		req.Method = pbuf.ReqContext_DELETE
	}

	headers := make(map[string]*pbuf.Strings)
	for key, vals := range rc.Headers {
		headers[key] = &pbuf.Strings{Values: vals}
	}
	req.Headers = headers

	var params []*pbuf.Param
	for _, param := range rc.Params {
		params = append(params, &pbuf.Param{
			Key:   param.Key,
			Value: param.Value,
		})
	}
	req.Params = params

	queries := make(map[string]*pbuf.Strings)
	for key, vals := range rc.Queries {
		queries[key] = &pbuf.Strings{Values: vals}
	}
	req.Queries = queries

	return &req
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// Get returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) Get(name string) (string, bool) {
	for _, entry := range ps {
		if entry.Key == name {
			return entry.Value, true
		}
	}
	return "", false
}

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) (va string) {
	va, _ = ps.Get(name)
	return
}

type Param struct {
	Key   string
	Value string
}

type RespContext struct {
	Code    int
	Type    string
	Payload []byte
}

func FromPbuf(rc *pbuf.RespContext) *RespContext {
	resp := RespContext{
		Code:    int(rc.Code),
		Payload: rc.Payload,
	}

	switch rc.Type {
	case pbuf.RespContext_JSON:
		resp.Type = "JSON"
	case pbuf.RespContext_RAW:
		resp.Type = "RAW"
	case pbuf.RespContext_XML:
		resp.Type = "XML"
	}

	return &resp
}
