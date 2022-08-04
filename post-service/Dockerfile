FROM golang:1.18rc1-alpine3.15
RUN mkdir post
COPY . /post
WORKDIR /post
RUN go mod tidy
RUN go build -o main cmd/main.go
RUN make migrate-up
CMD ./main
EXPOSE 7007