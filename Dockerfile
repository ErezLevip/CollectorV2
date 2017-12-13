FROM iron/go:dev
WORKDIR /app
ENV SRC_DIR="/go/src/Collector"
COPY . $SRC_DIR
WORKDIR $SRC_DIR
RUN go build -o myapp && cp myapp /app
ENTRYPOINT ["./myapp"]

//docker build -t collector .
//docker run --rm -p 8080:8080 collector