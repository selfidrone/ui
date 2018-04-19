FROM alpine

RUN adduser -h /home/dronelive -D dronelive dronelive

COPY ./drone-live /home/dronelive/
COPY ./index.html /home/dronelive/
RUN chmod +x /home/dronelive/drone-live

WORKDIR /home/dronelive

USER dronelive

ENTRYPOINT ["/home/dronelive/drone-live"]
