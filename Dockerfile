FROM ubuntu:17.04

RUN \
    apt-get update -q && \
    apt-get install -q -y wget && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ zesty-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    \
    apt-get update -q && \
    apt-get install -q -y git golang-go postgresql-10 postgresql-contrib-10

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

ENV PGHOST /tmp/.s.PGSQL.5432
ENV PGDATABASE technopark
EXPOSE 5000
CMD /etc/init.d/postgresql start && \
    sleep 10 && \
    ./technopark-db