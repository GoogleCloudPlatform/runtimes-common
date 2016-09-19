#Tag Check Intermediate Cloudbuild Step
This code serves as an intermediate container in a `gcloud container builds
create` command to check for existence of the specified tag on the target
container, and exit the build if it exists. The motivation for this is to
prevent us from unintentionally overwriting a tag in a remote repository.

*Note:* This step was originally meant to be a "plug-and-play" container,
meaning there would be zero overhead to integrating this with any cloud build
job; however, since gcloud auth does not currently transfer between intermediate
containers in a build, we need to build this image with the correct credentials
baked into the container. In the future, this will not be necessary, and it will
be possible to simply stick this step into your build config file without
prebuilding the image.

##Steps to Build Container Image
1. Clone this repository.
2. Navigate to your GCP project dashboard. Select "IAM & Admin", and then "Service Accounts".
3. Create a new service account for your project. Make sure and grant this account read access to your project.
4. Create a private key for this service account, *in JSON form (NOT P12 form)*. Save this key as `auth.json` in the same directory as this repository (this should overwrite the placeholder `auth.json` file).
5. Ensure you're logged into an account with write permissions to your project's repository.
6. Issue the following command from this repository's root directory: `gcloud alpha container builds create -t {your_repository}/{your_project}/check_if_tag_exists:latest .`
7. In your target project's cloudbuild.yaml file, add the following build step *before* your build is actually executed:
```- name: {your_repository}/{your_project}/check_if_tag_exists:latest```

If everything is successful, you will see a `check_if_tag_exists` container in the *Container Engine* view of your project on the GCP dashboard. This is the intermediate image that will be run in your cloudbuilds to ensure that you don't overwrite existing tags for your project.

