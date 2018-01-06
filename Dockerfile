FROM ubuntu:16.04

RUN apt-get update -q && \
    apt-get install -q -y git golang-go postgresql postgresql-contrib

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin"

ENV TZ Europe/Moscow
RUN echo "$TZ" > /etc/timezone

USER postgres
RUN /etc/init.d/postgresql start && \
    psql --command "CREATE USER technopark WITH SUPERUSER PASSWORD 'technopark';" && \
    createdb -E UTF8 -T template0 -O technopark technopark && \
    /etc/init.d/postgresql stop

USER root
WORKDIR /go/src/github.com/reo7sp/technopark-db
COPY . .
RUN go get ./...
RUN go build

ENV DATABASE_URL postgres://technopark:technopark@localhost/technopark?sslmode=disable
EXPOSE 5000
CMD /etc/init.d/postgresql start && \
    sleep 10 && \
    ./technopark-db