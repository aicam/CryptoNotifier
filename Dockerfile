FROM golang:1.16-alpine

WORKDIR /app

ADD . .

RUN go mod download
#RUN go get github.com/githubnemo/CompileDaemon

#ENTRYPOINT CompileDaemon -command="./main"
RUN go build -o /main
CMD [ "/main" ]
