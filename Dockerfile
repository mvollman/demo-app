FROM ubuntu:18.04

COPY app /usr/bin/app

EXPOSE 8080
CMD /usr/bin/app
