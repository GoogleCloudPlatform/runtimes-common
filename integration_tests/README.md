Integration Testing Framework
=============

This code builds an image which serves as a framework to run basic integration tests for a language runtime image. It relies on the existence of a sample application which fulfills a spec (detailed below), that will be built and deployed on top of the newly built runtime image. These tests are intended to be run as part of a post-push verification step in a continuous deployment system. The test driver can be run directly via any build host, or as a build step in a [Google Cloud Container Build](https://cloud.google.com/container-builder/docs/overview)).

To run these tests through a cloudbuild, add the following build step to the **end** of your build config (cloudbuild.yaml or cloudbuild.json):

	name: gcr.io/gcp-runtimes/integration_test
	args: [
		'-i', <target_staging_image>,
		'-d', <sample_application_directory>
	]

**It's crucial that this step appears after your image under test has been built**; without a built image, there will be nothing to test, and your build will fail!


The sample application directory should contain the application fulfilling the integration test spec, as well as the build artifacts necessary to deploy the application via gcloud (namely, a templated Dockerfile and an app.yaml).

Alternatively, the application can be manually deployed *before* running the tests; in this scenario, the '--no-deploy' flag can be passed to the build step to opt out of deploying, in tandem with the URL at which the deployed application can be accessed:

	args: [
		...,
		'--no-deploy',
		'--url', '<application_url>'
	]

In addition, each test (detailed below) can be skipped by passing the corresponding flag to the build step. For example, to opt out of the monitoring test in the run, pass the `--no-monitoring` flag to the build step.

## Tests

### Serving (Root)
#####` - GET http://<application_url>/`

*Response*

- If successful, the application should respond with the string 'Hello World!'

This is the most basic integration test. The driver performs a GET request to the root endpoint at the deployed application's URL. It retrieves data from the application, and verifies it matches the expected output (the text ‘Hello World!’).


### Logging
#####` - GET http://<application_url>/logging`

*Request Body*

| Property Name | Value | Description |
| --- | --- | --- |
| log_name      | string; 16 characters; alphabetic (upper/lower case) | The name of the log to write in Stackdriver. |
| token | 16 character (8 byte) hexadecimal string (uppercase) | The token to write in Stackdriver logs. |

*Response*

- If successful, the application should respond with a 200 response code and the string ‘OK’.

This tests the runtime’s integration with Stackdriver Logging. The driver will generate a log name and token, and send this payload to the sample application via a POST request. Once the application receives the payload, it will create a [Cloud Logging Client](https://github.com/GoogleCloudPlatform/google-cloud-python/blob/master/logging/google/cloud/logging/client.py) instance. The application will use this client to write a log entry to Stackdriver with the provided name and token value (through the default channel, which for all silver language runtimes should be stdout), and then signal back to the driver that either the log entry has been written successfully, or that an error was encountered (causing the test to fail). If the write succeeded, the test driver will wait a short period of time for the log entry to propagate through Stackdriver, then retrieve all log entries written to stdout in the previous 2 minutes and verify that the provided name/token pair appears in Stackdriver.


### Monitoring
#####` - GET http://<application_url>/monitoring`

*Request Body*

| Property Name | Value | Description |
| --- | --- | --- |
| name      | string; 16 characters; alphabetic (upper/lower case) | The name of the metric to write in Stackdriver; **always prefixed with ‘custom.googleapis.com/’**. |
| token | integer 64-bit | The metric value to write into Stackdriver. |

*Response*

- If successful, the application should respond with a 200 response code and the string ‘OK’.

This tests the runtime’s integration with Stackdriver Logging. The driver will generate a log name and token, and send this payload to the sample application via a POST request. Once the application receives the payload, it will create a [Cloud Logging Client](https://github.com/GoogleCloudPlatform/google-cloud-python/blob/master/logging/google/cloud/logging/client.py) instance. The application will use this client to write a log entry to Stackdriver with the provided name and token value (through the default channel, which for all silver language runtimes should be stdout), and then signal back to the driver that either the log entry has been written successfully, or that an error was encountered (causing the test to fail). If the write succeeded, the test driver will wait a short period of time for the log entry to propagate through Stackdriver, then retrieve all log entries written to stdout in the previous 2 minutes and verify that the provided name/token pair appears in Stackdriver.


### Error Reporting/Exceptions
#####` - GET http://<application_url>/exception`

*Request Body*

| Property Name | Value | Description |
| --- | --- | --- |
| token | integer 64-bit | The metric value to write into Stackdriver. |

*Response*

- If successful, the application should respond with a 200 response code and the string ‘OK’.

This tests the runtime’s ability to report errors to Stackdriver. The driver will generate a request payload and POST it to the sample application. Upon receiving the payload, the application will create a Cloud Error Reporting Client. It will then use this client to report a generic exception. Additionally, it will use the client to report an exception with the provided token. Finally, it will report back to the test driver indicating that the exceptions were recorded successfully, or that an error was encountered, failing the test.

At this time, there is no support within the google-cloud-python client library for reading exceptions from Stackdriver, though the public API exists. It is potentially planned to add support within this client library to do this, enabling the test driver to actually verify the specific exception was recorded in Stackdriver successfully.
