FROM golang:1.23.0

COPY --chown=go/go . /home/go/src
WORKDIR /home/go/src
RUN go build main.go

EXPOSE 8080

CMD ["./main"]
