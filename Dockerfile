FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

RUN ls -l /app

CMD ["make", "run"]