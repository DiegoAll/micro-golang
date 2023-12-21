# Faster than Multistage build for production image

FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD ["/app/authApp"]