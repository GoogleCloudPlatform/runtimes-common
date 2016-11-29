
FROM gcr.io/cloud-builders/gcloud


ENV PATH=$PATH:/builder/google-cloud-sdk/bin/

# Setup Google Cloud SDK (latest)
RUN apt-get update && apt-get install -y wget unzip && \
	wget -nv https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip && \
    mkdir -p /builder && \
    unzip -qq google-cloud-sdk.zip -d /builder && \
    rm google-cloud-sdk.zip && \
    /builder/google-cloud-sdk/install.sh --usage-reporting=false \
        --bash-completion=false \
        --disable-installation-options && \
    # Install alpha and beta components
    /builder/google-cloud-sdk/bin/gcloud -q components install alpha beta && \
    /builder/google-cloud-sdk/bin/gcloud -q components update && \

    # TODO(jasonhall): These lines pin gcloud to a particular version.
    google-cloud-sdk/bin/gcloud components update --version 110.0.0 && \
    google-cloud-sdk/bin/gcloud config set component_manager/disable_update_check 1 && \
    google-cloud-sdk/bin/gcloud -q components update && \

    apt-get install -y python-dev gcc curl && \

    cp -v /usr/local/bin/python /usr/bin/python && \
 	cp -v /usr/local/bin/pip /usr/bin/pip && \
 	cp -v /usr/bin/envsubst /usr/local/bin/ && \


RUN pip install --upgrade google-cloud

COPY testsuite /testsuite

ENTRYPOINT ["/testsuite/driver.py"]

