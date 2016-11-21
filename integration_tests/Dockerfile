# Inherit Python 2.7 and gcloud
FROM gcr.io/cloud-builders/gcloud

# Install Docker
RUN apk update && \
	apk add gettext && \
	cp -v /usr/bin/envsubst /usr/local/bin/ && \
    apk add curl && \
    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    rm -rf /var/cache/apk/*

COPY app.yaml /
COPY build_and_deploy.sh /
COPY auth.json /
RUN chmod +x /build_and_deploy.sh

ENTRYPOINT ["/build_and_deploy.sh"]
