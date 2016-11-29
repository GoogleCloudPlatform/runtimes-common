FROM gcr.io/google-appengine/python

# gcloud first 
ENV PATH=$PATH:/builder/google-cloud-sdk/bin/

RUN apt-get update && \
    apt-get install -y wget unzip ca-certificates git && \
    # Setup Google Cloud SDK (latest)
    wget -nv https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    mkdir -p /builder && \
    unzip -qq google-cloud-sdk.zip -d /builder && \
    rm google-cloud-sdk.zip && \
    /builder/google-cloud-sdk/install.sh --usage-reporting=false \
        --bash-completion=false \
        --disable-installation-options && \
    # Install alpha and beta components
    /builder/google-cloud-sdk/bin/gcloud -q components install alpha beta && \
    /builder/google-cloud-sdk/bin/gcloud -q components update

RUN apt-get install -y python-dev gcc curl && \

    cp -v /usr/local/bin/python /usr/bin/python && \
 	cp -v /usr/local/bin/pip /usr/bin/pip && \
 	cp -v /usr/bin/envsubst /usr/local/bin/ && \

    curl https://get.docker.com/builds/Linux/x86_64/docker-1.9.1 > /usr/bin/docker && chmod +x /usr/bin/docker && \
    # Clean up
    apt-get remove -y wget unzip gcc && \
    apt-get clean

RUN pip install --upgrade google-cloud-logging

COPY testsuite /testsuite

ENTRYPOINT ["/testsuite/driver.py"]
