FROM gcr.io/cloud-builders/gcloud

COPY dedupe.py /

# This is the file containing the auth information for your service account. 
# Necessary inside the container until gcloud supports cross container authentication when using cloudbuild.
# Once this is addressed, remove this and update dedupe.py accordingly.
COPY auth.json /

# Since we're inheriting from the gcloud container, overwrite the entrypoint so we don't start in a gcloud command.
ENTRYPOINT []

