#!/usr/bin/python

import json
import logging
import requests
import time

import test_util

from google.cloud import monitoring as gcloud_monitoring


def _test_exception(base_url):
  logging.info("testing monitoring")
  url = base_url + test_util.MONITORING_ENDPOINT
  
  print ''
