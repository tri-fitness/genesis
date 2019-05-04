FROM golang:1.11

# labels.
LABEL maintainer="Jon Freer"
LABEL maintainer="freerjm@miamioh.edu"

# transfer binary and configuration file.
COPY genesis /go
COPY configuration.yaml /go
COPY db/dbconf.yml /go/db/dbconf.yml
COPY db/migrations/* /go/db/migrations/

	
# execute the binary.
CMD ["./genesis"]

