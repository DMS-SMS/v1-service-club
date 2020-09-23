FROM alpine
ADD club-service /club-service
ENTRYPOINT [ "/club-service" ]
