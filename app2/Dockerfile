FROM golang:1.14-stretch

WORKDIR /app

COPY /app2 .

RUN go mod download && go get github.com/pilu/fresh

COPY ./app2/runner.conf /

ENTRYPOINT ["fresh", "-c", "./runner.conf"]