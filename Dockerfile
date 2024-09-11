FROM golang:1.21.1

COPY --chown=go/go . /home/go/src
WORKDIR /home/go/src
RUN go build main.go

EXPOSE 8080

CMD ["./main"]
