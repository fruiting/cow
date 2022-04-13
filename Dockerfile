FROM golang:1.16-alpine
WORKDIR /app
COPY . .

RUN cd cmd/cow/ && \
    go build && \
    chmod 777 cow && \
    mv cow ../.. && \
    cd ../..

CMD [ "./cow" ]
