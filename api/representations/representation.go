package representations

import (
	"time"
)

type Representation interface {

	// ContentLanguage indicates the language(s) of the
	// consuming audience for this representation.
	ContentLanguage() string

	// ContentType indicates the media type of the representation.
	ContentType() string

	// ContentEncoding indicates the scheme(s) used to encode
	// the representation.
	ContentEncoding() string

	// LastModified indicates when the representation was last modified.
	LastModified() *time.Time

	// ETag represents the entity tag of the representation.
	ETag() string

	// AsBytes provides the representation as raw bytes. This is
	// used when passing the representation onto the wire.
	AsBytes() []byte
}

type representation struct {
	contentLanguage string
	contentType     string
	contentEncoding string
	lastModified    *time.Time
	eTag            string
}

func (r *representation) ContentLanguage() string {
	return r.contentLanguage
}

func (r *representation) ContentType() string {
	return r.contentType
}

func (r *representation) ContentEncoding() string {
	return r.contentEncoding
}

func (r *representation) LastModified() *time.Time {
	return r.lastModified
}

func (r *representation) ETag() string {
	return r.eTag
}
