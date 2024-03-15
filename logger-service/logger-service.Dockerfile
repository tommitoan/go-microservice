FROM alpine:latest

RUN mkdir /app

COPY loggerServiceApp /app

# Run the server executable
CMD [ "/app/loggerServiceApp" ]