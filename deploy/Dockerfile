FROM gcr.io/eoscanada-shared-services/dubuntu:latest

RUN apt update && apt-get -y install curl ca-certificates && rm -rf /var/cache/apt/*

COPY /exchange /app/exchange
COPY /motd /etc/motd

RUN echo cat /etc/motd >> /root/.bashrc

ENV PATH="/app:${PATH}"
