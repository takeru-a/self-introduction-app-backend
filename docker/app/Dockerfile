FROM golang:1.19-alpine

WORKDIR /go/src
COPY ./src .
RUN apk upgrade --update && \
    apk --no-cache add git
RUN go install github.com/cosmtrek/air@latest  
# RUN go mod download
RUN go mod tidy


# copy  local: ./src   docker: ./go/src

# RUN go get -u github.com/cosmtrek/air && \
#     go build -o /go/bin/air github.com/cosmtrek/air

# air -c .air.toml
CMD ["air", "-c", ".air.toml"]
