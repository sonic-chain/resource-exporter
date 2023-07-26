FROM nvidia/cuda:11.2.2-base-ubuntu20.04 AS builder

RUN apt update \
    && apt install  -y wget gcc \
    && wget -c https://golang.org/dl/go1.18.1.linux-amd64.tar.gz -O - | tar -xz -C /usr/local

ENV PATH /usr/local/go/bin:$PATH

WORKDIR /app
COPY . .
RUN go mod download \
    && go build -o resource-exporter .


FROM nvidia/cuda:11.2.2-base-ubuntu20.04
COPY --from=builder /app/resource-exporter .

CMD ["./resource-exporter"]