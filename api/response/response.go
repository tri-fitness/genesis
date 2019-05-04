package response

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type ResponseBuilder interface {
	Created(location url.URL) ResponseBuilder
	NotFound() ResponseBuilder
	BadRequest() ResponseBuilder
	InternalServerError() ResponseBuilder
	NotAcceptable() ResponseBuilder
	OK() ResponseBuilder
	NoContent() ResponseBuilder
	Header(string, ...string) ResponseBuilder
	Body([]byte) ResponseBuilder
	StatusCode(int) ResponseBuilder
	Respond() error
	WithError(error) ResponseBuilder
	Unauthorized() ResponseBuilder
	Challenge(string, string, string) ResponseBuilder
}

type builder struct {
	statusCode int
	headers    http.Header
	body       []byte
	responded  bool
	writer     http.ResponseWriter
}

func Builder(rw http.ResponseWriter) ResponseBuilder {
	return &builder{writer: rw, headers: make(map[string][]string)}
}

func (b *builder) Created(location url.URL) ResponseBuilder {
	bb := b.Header("Location", location.String())
	return bb.StatusCode(http.StatusCreated)
}

func (b *builder) NotAcceptable() ResponseBuilder {
	return b.StatusCode(http.StatusNotAcceptable)
}

func (b *builder) BadRequest() ResponseBuilder {
	return b.StatusCode(http.StatusBadRequest)
}

func (b *builder) InternalServerError() ResponseBuilder {
	return b.StatusCode(http.StatusInternalServerError)
}

func (b *builder) NotFound() ResponseBuilder {
	return b.StatusCode(http.StatusNotFound)
}

func (b *builder) Unauthorized() ResponseBuilder {
	return b.StatusCode(http.StatusUnauthorized)
}

func (b *builder) Challenge(scheme, realm, charset string) ResponseBuilder {
	challenge := fmt.Sprintf("%s realm=%q charset=%q", scheme, realm, charset)
	b.Header("WWW-Authenticate", challenge)
	return b
}

func (b *builder) OK() ResponseBuilder {
	return b.StatusCode(http.StatusOK)
}

func (b *builder) NoContent() ResponseBuilder {
	return b.StatusCode(http.StatusNoContent)
}

func (b *builder) StatusCode(statusCode int) ResponseBuilder {
	b.statusCode = statusCode
	return b
}

func (b *builder) Body(bytes []byte) ResponseBuilder {
	b.body = bytes
	return b
}

func (b *builder) Header(name string, values ...string) ResponseBuilder {
	for _, value := range values {
		b.headers.Add(name, value)
	}
	return b
}

func (b *builder) WithError(err error) ResponseBuilder {
	return b.Body([]byte(err.Error()))
}

func (b *builder) Respond() error {

	if b.responded {
		return errors.New("the resposne has already been sent")
	}

	for header, values := range b.headers {
		for _, value := range values {
			b.writer.Header().Add(header, value)
		}
	}
	b.writer.WriteHeader(b.statusCode)
	b.writer.Write(b.body)
	return nil
}
