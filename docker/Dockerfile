FROM golang:1.11

# labels.
LABEL maintainer="Jon Freer"
LABEL maintainer="freerjm@miamioh.edu"

# transfer binary and configuration file.
COPY genesis /go
COPY configuration.yaml /go
COPY infrastructure/dbconf.yml /go/infrastructure/dbconf.yml
COPY infrastructure/migrations/* /go/infrastructure/migrations/

	
# execute the binary.
CMD ["./genesis"]

