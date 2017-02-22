#!/usr/bin/python

import logging
import os
import subprocess
import tempfile

logging.getLogger().setLevel(logging.DEBUG)


RUNTIME_BUCKET = 'runtime-builders'
RUNTIME_BUCKET_PREFIX = 'gs://{0}/'.format(RUNTIME_BUCKET)


def write_to_gcs(gcs_path, file_contents):
    try:
        fd, f_name = tempfile.mkstemp(text=True)
        os.write(fd, file_contents)

        command = ['gsutil', 'cp', f_name, gcs_path]
        try:
            output = subprocess.check_output(command)
        except subprocess.CalledProcessError as e:
            logging.error('Error encountered when writing to GCS!: {0}'
                          .format(output))
            logging.error(e)
    finally:
        os.remove(f_name)


def _get_file_from_gcs(gcs_file, temp_file):
    command = ['gsutil', 'cp', gcs_file, temp_file]
    try:
        subprocess.check_output(command, stderr=subprocess.STDOUT)
        return True
    except subprocess.CalledProcessError as e:
        logging.error('Error when retrieving file from GCS! {0}'
                      .format(e.output))
        return False
