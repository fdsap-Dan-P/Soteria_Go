FROM golang:1.20

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /report_webtool

EXPOSE 8083

CMD [ "/report_webtool" ]