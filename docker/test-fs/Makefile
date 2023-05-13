PROXY="http://192.168.2.117:7890"
VERSION=$(shell  cat VERSION)

.PHONY:  base
base:
	docker build -f  base.Dockerfile -t lvsoso/base:${VERSION} .

.PHONY:  base-proxy
base-proxy:
	docker build --build-arg HTTP_PROXY="${PROXY}"   --build-arg HTTPS_PROXY="${PROXY}"  -f  base.Dockerfile -t lvsoso/base:${VERSION} .

.PHONY:  base-run
base-run:
	docker run -it --rm  --env HTTP_PROXY="${PROXY}"   --env HTTPS_PROXY="${PROXY}"   lvsoso/base:${VERSION}  /bin/bash


.PHONY:  test-env
test-env:
	docker build  --build-arg  VERSION="${VERSION}"  --build-arg HTTP_PROXY="${PROXY}"   --build-arg HTTPS_PROXY="${PROXY}"  -f  Dockerfile -t lvsoso/test-env:${VERSION} .


.PHONY:  test-env-run
test-env-run:
	docker run -it --rm  --env HTTP_PROXY="${PROXY}"   --env HTTPS_PROXY="${PROXY}"   lvsoso/test-env:${VERSION}  /bin/bash