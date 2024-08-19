FROM nvidia/cuda:11.2.2-base-ubuntu20.04 AS builder

RUN apt update \
    && apt install  -y wget gcc git \
    && wget -c https://golang.org/dl/go1.18.1.linux-amd64.tar.gz -O - | tar -xz -C /usr/local

ENV PATH /usr/local/go/bin:$PATH

WORKDIR /app
COPY . .
RUN go mod download \
    && go build -o resource-exporter .


FROM nvidia/cuda:11.2.2-base-ubuntu20.04
COPY --from=builder /app/resource-exporter .
COPY self_heal.sh /usr/local/bin/self_heal.sh
RUN chmod +x /usr/local/bin/self_heal.sh
CMD ["sh", "-c", "/usr/local/bin/self_heal.sh & ./resource-exporter"]