if [ ${APP_ENV} = production ]; \
  then \
	${GOPATH}/bin/main; \
	else \
	go get github.com/pilu/fresh && \
	fresh -c ; \
	fi