FROM gcr.io/google-appengine/debian8

RUN apt-get update && apt-get install -y --force-yes software-properties-common python-software-properties
RUN add-apt-repository -y ppa:webupd8team/java
RUN apt-get install -y --force-yes sl=3.03-17 figlet
RUN apt-get remove software-properties-common python-software-properties
 && apt-get autoremove && apt-get clean

ENTRYPOINT '/bin/bash'
