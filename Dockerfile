FROM golang:1.22.1

WORKDIR /app

COPY . /app

WORKDIR /app/cmd

RUN go mod download

RUN go build -o /chatbot

CMD [ "/chatbot" ]