package simplelb

import "github.com/valyala/fasthttp"

// HeaderSaver provides a structure to save headers.
type HeaderSaver struct {
	headers map[string]string
}

// Setter provides a set key:value  func.
type Setter interface {
	// set sets the given 'key: value' header.
	Set(key, value string)
}

func saveHeaders(res *fasthttp.Response) HeaderSaver {
	resHeaders := make(map[string]string)

	res.Header.VisitAll(func(k, v []byte) {
		key := string(k)
		value := string(v)
		if val, ok := resHeaders[key]; ok {
			resHeaders[key] = val + "," + value
		} else {
			resHeaders[key] = value
		}
	})

	resHeaders["Server"] = "simplelb"

	return HeaderSaver{headers: resHeaders}
}

func (s HeaderSaver) set(setter Setter) {
	for k, v := range s.headers {
		setter.Set(k, v)
	}
}

// Deleter provides the ability to delete by key.
type Deleter interface {
	// Del deletes header with the given key.
	Del(key string)
}

// HeaderHop provides the structure for headers to be hopped.
type HeaderHop struct {
	headers []string
}

// Del deletes the hopped headers.
func (h *HeaderHop) Del(reqHeader Deleter) {
	for _, header := range h.headers {
		reqHeader.Del(header)
	}
}

// MakeHeaderHop makes a HeaderHop structure.
func MakeHeaderHop() *HeaderHop {
	// Hop-by-hop headers. These are removed when sent to the backend.
	// As of RFC 7230, hop-by-hop headers are required to appear in the
	// Connection header field. These are the headers defined by the
	// obsoleted RFC 2616 (section 13.5.1) and are used for backward
	// compatibility.
	return &HeaderHop{
		headers: []string{
			"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
			"Connection", "Keep-Alive", "Proxy-Authenticate", "Proxy-Authorization",
			"Te",      // canonicalized version of "TE"
			"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
			"Transfer-Encoding", "Upgrade",
		},
	}
}
