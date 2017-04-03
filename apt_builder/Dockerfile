FROM gcr.io/google-appengine/python

RUN pip install ruamel.yaml==0.13.13

COPY install_packages.py /

ENTRYPOINT ["/install_packages.py"]
