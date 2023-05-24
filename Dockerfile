# syntax=docker/dockerfile:1

FROM golang:1.20
 
# Add libtensorflow, note: must match your docker env
ENV FILENAME=libtensorflow-cpu-linux-x86_64-2.11.0.tar.gz
RUN wget -q --no-check-certificate https://storage.googleapis.com/tensorflow/libtensorflow/${FILENAME}
RUN tar -C /usr/local -xzf ${FILENAME}

ENV LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./
COPY ./match/*.go ./match/

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /tenseimage

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD [ "/tenseimage" ]
