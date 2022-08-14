FROM golang:1.18rc1-alpine3.15
RUN mkdir user
COPY . /user
WORKDIR /user
RUN go mod tidy
# RUN migrate -source file:./migrations/ -database postgres://hatsker:1@dbpost:5433/postdb up
RUN go build -o main cmd/main.go
# RUN make migrate-up
CMD ./main
EXPOSE 9000