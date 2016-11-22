# Inherit Python 2.7 and gcloud
FROM gcr.io/cloud-builders/gcloud

# Install Docker
RUN apk update && \
	cp -v /usr/local/bin/python /usr/bin/python && \
	apk add gettext && \
	cp -v /usr/bin/envsubst /usr/local/bin/ && \
    apk add curl && \
    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    rm -rf /var/cache/apk/*

COPY app.yaml /
COPY auth.json /

COPY testsuite /testsuite

ENTRYPOINT ["/testsuite/driver.py"]
