FROM ubuntu:16.04

RUN apt-get update -q && \
    apt-get install -q -y wget unzip

WORKDIR /tmp
RUN wget https://bozaro.github.io/tech-db-forum/linux_amd64.zip && \
    unzip linux_amd64.zip && \
    rm linux_amd64.zip && \
    chmod +x ./tech-db-forum && \
    mv ./tech-db-forum /usr/local/bin/tech-db-forum

CMD ["tech-db-forum"]