FROM gcr.io/cloud-builders/gcloud

COPY main.py /

# This is the file containing the auth information for your service account. 
# Necessary inside the container until gcloud supports cross container authentication when using cloudbuild.
COPY auth.json /
RUN /builder/google-cloud-sdk/bin/gcloud auth activate-service-account --key-file /auth.json

# Since we're inheriting from the gcloud container, overwrite the entrypoint so we don't start in a gcloud command.
ENTRYPOINT []
