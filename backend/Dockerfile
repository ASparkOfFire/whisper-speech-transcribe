FROM debian:bullseye-slim AS builder

RUN DEBIAN_FRONTEND=noninteractive apt -y update && apt -y upgrade && \
    apt -y install build-essential git cmake

WORKDIR /src

RUN git clone https://github.com/ggerganov/whisper.cpp.git

WORKDIR /src/whisper.cpp/build

RUN cmake .. -DCMAKE_INSTALL_PREFIX=/app/whisper && \
    make -j$(nproc) && \
    make install

# Final stage
FROM golang:1.24-bullseye

WORKDIR /app

COPY . .

# Copy libraries and headers from the builder stage
COPY --from=builder /app/whisper/lib /app/lib
COPY --from=builder /app/whisper/include /app/include

ENV LD_LIBRARY_PATH=/app/lib

ENV CGO_CFLAGS="-I/app/include"
ENV CGO_LDFLAGS="-L/app/lib"

RUN go build -o bin/whisper main.go

CMD ["./bin/whisper"]
