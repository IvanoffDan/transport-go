FROM golang

ARG app_env
ENV APP_ENV $app_env
ENV PATH /go/src/app/:$PATH

COPY ./go-server/src /go/src
WORKDIR ${GOPATH}/src/app/transit_realtime
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN protoc --go_out=. *.proto

WORKDIR ${GOPATH}

RUN go get ./src/app
RUN go build -o ./bin/main ./src/app/main.go

# -c ${GOPATH}/src/app/runner.conf;

CMD if [ ${APP_ENV} = production ]; \
	then \
	${GOPATH}/bin/main; \
	else \
	go get github.com/pilu/fresh && \
	fresh -c ${GOPATH}/src/app/runner.conf; \
	fi
	
EXPOSE 8080