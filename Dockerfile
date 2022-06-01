FROM golang:1.18 as builder
MAINTAINER Kirill Kiling

# ARG security: https://bit.ly/2oY3pCn
WORKDIR /build

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
COPY . .

RUN go build -v -o bot cmd/bot/main.go

FROM golang:1.18

WORKDIR /app
COPY --from=builder /build/bot /app

# install python
RUN apt-get update
RUN apt-get -y install python3
RUN apt-get -y install python3-setuptools
RUN apt-get -y install python3-pip
RUN apt-get install -y wkhtmltopdf
COPY ./requirements.txt .
COPY ./html2png.py .
RUN pip install -r requirements.txt

RUN mkdir /configs
RUN mkdir /data
#USER 1000:1000

CMD ./bot --config "/configs/config.yml"