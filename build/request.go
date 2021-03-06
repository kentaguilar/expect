package build

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type RequestBuilder struct {
	Request *http.Request
}

func Request() *RequestBuilder {
	u, _ := url.Parse("http://local.test/spice")
	return &RequestBuilder{
		Request: &http.Request{
			Method:     "GET",
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Host:       "localtest",
			URL:        u,
			Header:     make(http.Header),
		},
	}
}

func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.Request.Method = method
	return rb
}

func (rb *RequestBuilder) Proto(major, minor int) *RequestBuilder {
	rb.Request.Proto = "HTTP/" + strconv.Itoa(major) + "." + strconv.Itoa(minor)
	rb.Request.ProtoMajor = major
	rb.Request.ProtoMinor = minor
	return rb
}

func (rb *RequestBuilder) URL(u *url.URL) *RequestBuilder {
	rb.Request.URL = u
	return rb
}

func (rb *RequestBuilder) URLString(u string) *RequestBuilder {
	url, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	rb.Request.URL = url
	return rb
}

func (rb *RequestBuilder) Path(path string) *RequestBuilder {
	rb.Request.URL.Path = path
	return rb
}

func (rb *RequestBuilder) RawQuery(query string) *RequestBuilder {
	rb.Request.URL.RawQuery = query
	return rb
}

func (rb *RequestBuilder) Host(host string) *RequestBuilder {
	rb.Request.URL.Host = host
	rb.Request.Host = host
	return rb
}

func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	rb.Request.Header.Set(key, value)
	if key == "Content-Encoding" && value == "gzip" && rb.Request.Body != nil {
		rb.compressBody()
	}
	return rb
}

func (rb *RequestBuilder) Body(b string) *RequestBuilder {
	rb.Request.Body = ioutil.NopCloser(bytes.NewBufferString(b))
	if rb.Request.Header.Get("Content-Encoding") == "gzip" {
		rb.compressBody()
	}
	return rb
}

func (rb *RequestBuilder) compressBody() {
	buffer := new(bytes.Buffer)
	writer := gzip.NewWriter(buffer)
	io.Copy(writer, rb.Request.Body)
	writer.Close()
	rb.Request.Body.Close()
	rb.Request.Body = ioutil.NopCloser(buffer)
}
