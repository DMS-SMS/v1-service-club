
.PHONY: build
build: proto
	GOOS=linux GOARCH=amd64 go build -o club-service *.go

.PHONY: image
image:
	docker build . -t dms-sms-service-club:${VERSION}

.PHONY: upload
upload:
	docker tag dms-sms-service-club:${VERSION} jinhong0719/dms-sms-service-club:${VERSION}.RELEASE
	docker push jinhong0719/dms-sms-service-club:${VERSION}.RELEASE

.PHONY: pull
pull:
	docker pull jinhong0719/dms-sms-service-club:${VERSION}.RELEASE

.PHONY: run
run:
	docker-compose -f ./docker-compose.yml up -d
