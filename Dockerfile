FROM ubuntu:16.04

ENV DEBIAN_FRONTEND noninteractive

RUN \
    apt-get update -q && \
    apt-get install -q -y wget && \
    wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add - && \
    echo "deb http://apt.postgresql.org/pub/repos/apt/ xenial-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    wget https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz && \
    tar -xvf go1.9.2.linux-amd64.tar.gz &&\
    mv go /usr/local/ && \
    \
    apt-get update -q && \
    apt-get install -q -y git postgresql-10 postgresql-contrib-10

ENV GOROOT /usr/local/go
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

RUN ln -s /var/run/postgresql/10-main.pid /var/run/postgresql/.s.PGSQL.5432
RUN echo "local all all trust" > /etc/postgresql/10/main/pg_hba.conf
RUN echo "host all all 0.0.0.0/0 trust" >> /etc/postgresql/10/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/10/main/postgresql.conf

RUN echo "synchronous_commit = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "autovacuum = off" >> /etc/postgresql/10/main/postgresql.conf

WORKDIR /go/src/github.com/reo7sp/technopark-db
COPY . .
RUN go get ./...
RUN go build

ENV PGHOST /var/run/postgresql
ENV PGDATABASE technopark
ENV PGUSER technopark
ENV PGPASSWORD technopark
ENV KILL_POSTGRES 1
ENV DEBUG 0
EXPOSE 5000
CMD /etc/init.d/postgresql start && \
    sleep 10 && \
    ./technopark-db
    
