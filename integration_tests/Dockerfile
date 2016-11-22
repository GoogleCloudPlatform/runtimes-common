# Inherit Python 2.7 and gcloud
FROM gcr.io/cloud-builders/gcloud

# Install Docker
RUN apk update && \
<<<<<<< HEAD
<<<<<<< HEAD
	cp -v /usr/local/bin/python /usr/bin/python && \
=======
>>>>>>> wip; adding infrastructure for deploying integration test containers
=======
	cp -v /usr/local/bin/python /usr/bin/python && \
>>>>>>> fix python installation in container. add support for multiple sample apps
	apk add gettext && \
	cp -v /usr/bin/envsubst /usr/local/bin/ && \
    apk add curl && \
    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    rm -rf /var/cache/apk/*

COPY app.yaml /
<<<<<<< HEAD
COPY auth.json /

COPY testsuite /testsuite

ENTRYPOINT ["/testsuite/driver.py"]
=======
COPY build_and_deploy.sh /
COPY auth.json /
RUN chmod +x /build_and_deploy.sh

ENTRYPOINT ["/build_and_deploy.sh"]
>>>>>>> wip; adding infrastructure for deploying integration test containers
