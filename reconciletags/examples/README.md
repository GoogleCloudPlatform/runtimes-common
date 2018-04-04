Running Retagger with Kubernetes
================================
The retagger is a tool used to tag images. You will need a git repository which includes a json file in [this](https://github.com/GoogleCloudPlatform/runtimes-common/blob/master/reconciletags/sample.json) format, which contains all the images and associated tags that you want to retag. Then, with Kubernetes, you can create a Pod to retag those images once, or a CronJob which will run the retagger at scheduled intervals. The benefit of this is that if a change is made to the json file in your git repository, the CronJob will automatically retag your images at the next scheduled time.

Getting Started
---------------

To run the retagger through Kubernetes as either a Pod or a CronJob on a Google Container Engine cluster, please follow these steps:

__Creating a Service Account__

First, create a service account, which will later be used for authentication purposes. Name the account, and select `Storage > Storage Admin` as the role. Download as a JSON file, and rename this file to **retagger_secret.json**

__Create a Secret__

A Kubernetes secret will be used to authenticate to GCR, allowing the retagger to pull and push images. To create a secret named retagger-secret, run the following kubectl command:

`kubectl create secret generic retagger-secret  --from-file=[PATH TO retagger_secret.json]`

To ensure the secret was created, run

`kubectl get secrets`

retagger-secret should be in the output list of secrets.


Running Retagger as a Pod
-------------------------

To run the retagger as a Pod: include a link to your git repository and update the path in the container args to point to your json file. A sample is available in samples/retagger\_pod\_sample.yaml, which links to [this](https://github.com/priyawadhwa/retagger-example) git repo and passes in the appropriate path to sample.json in the container args.

**retagger_pod.yaml**

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: retagger
spec:
  containers:
  - name: retagger 
    image: gcr.io/gcp-runtimes/retagger:latest
    volumeMounts:
    - name: git-volume
      mountPath: /data
    - name: retagger-secret
      mountPath: /secret
    args: ["/data/[PATH TO JSON FILE IN YOUR GIT REPOSITORY]"]
    env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /secret/retagger_secret.json
  volumes:
  - name: git-volume
    gitRepo:
      repository: [LINK TO YOUR GIT REPOSITORY]
  - name: retagger-secret
    secret:
      secretName: retagger-secret
```

Create the pod with the following command, which will run the retagger on your json file:

`kubectl create -f retagger_pod.yaml`

Running Retagger as a CronJob 
-----------------------------

To run the retagger as a CronJob: include a link to your git repository, update the path in the container args to point to your json file, and select the desired schedule. A sample can be found in samples/retagger\_cronjob\_sample.yaml, which links to [this](https://github.com/priyawadhwa/retagger-example) git repo and passes in the appropriate path to sample.json in the container args. It is scheduled to run every minute.

**retagger_cronjob.yaml**

```yaml
apiVersion: batch/v2alpha1
kind: CronJob
metadata:
  name: retagger-cronjob
spec:
  schedule: [SCHEDULE]
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: retagger 
            image: gcr.io/gcp-runtimes/retagger:latest
            volumeMounts:
            - name: git-volume
              mountPath: /data
            - name: retagger-secret
              mountPath: /secret
            args: ["/data/[PATH TO JSON FILE IN YOUR GIT REPOSITORY]"]
            env:
            - name: GOOGLE_APPLICATION_CREDENTIALS
              value: /secret/retagger_secret.json
          volumes:
          - name: git-volume
            gitRepo:
              repository: [LINK TO YOUR GIT REPOSITORY]
          - name: retagger-secret
            secret:
              secretName: retagger-secret
          restartPolicy: OnFailure
```

Create the CronJob with the following command, which will run the retagger on the desired schedule. 

`kubectl create -f retagger_cronjob.yaml`