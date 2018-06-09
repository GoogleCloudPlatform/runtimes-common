# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""This package defines helpful utilities for FTL ."""
import os
import time
import logging
import subprocess
import tempfile
import datetime
import json

from ftl.common import constants
from ftl.common import ftl_error

from containerregistry.client.v2_2 import append
from containerregistry.transform.v2_2 import metadata


class FTLException(Exception):
    pass


def AppendLayersIntoImage(imgs):
    with Timing('Stitching layers into final image'):
        for i, img in enumerate(imgs):
            if i == 0:
                result_image = img
                continue
            diff_ids = img.diff_ids()
            for diff_id in diff_ids:
                lyr = img.blob(img._diff_id_to_digest(diff_id))
                overrides = CfgDctToOverrides(json.loads(img.config_file()))
                result_image = append.Layer(
                    result_image, lyr, diff_id=diff_id, overrides=overrides)
        return result_image


# This is a 'whitelist' of values to pass from the
# config_file of a DockerImage to an Overrides object
# _OVERRIDES_VALUES = ['created', 'Entrypoint', 'Env']
def CfgDctToOverrides(config_dct):
    """
    Takes a dct of config values and runs them through
    the whitelist
    """
    overrides_dct = {}
    for k, v in config_dct.iteritems():
        if k == 'created':
            # this key change is made as the key is
            # 'creation_time' in an Overrides object
            # but 'created' in the config_file
            overrides_dct['creation_time'] = v
    for k, v in config_dct['config'].iteritems():
        if k == 'Entrypoint':
            # this key change is made as the key is
            # 'entrypoint' in an Overrides object
            # but 'Entrypoint' in the config_file
            overrides_dct['entrypoint'] = v
        elif k == 'Env':
            # this key change is made as the key is
            # 'env' in an Overrides object
            # but 'Env' in the config_file
            overrides_dct['env'] = v
        elif k == 'ExposedPorts':
            # this key change is made as the key is
            # 'ports' in an Overrides object
            # but 'ExposedPorts' in the config_file
            overrides_dct['ports'] = v
    return metadata.Overrides(**overrides_dct)


class Timing(object):
    def __init__(self, descriptor):
        logging.info("starting: %s" % descriptor)
        self.descriptor = descriptor

    def __enter__(self):
        self.start = time.time()
        return self

    def __exit__(self, unused_type, unused_value, unused_traceback):
        end = time.time()
        logging.info('%s took %d seconds', self.descriptor, end - self.start)


def zip_dir_to_layer_sha(app_dir, destination_path, alter_symlinks=True):
    tar_dir = tempfile.mkdtemp()

    tar_path = tempfile.mktemp(suffix='.tar')
    txfrm_regex = 's,^,%s/,' % destination_path
    if alter_symlinks:
        txfrm_regex = 'flags=r;s,^,%s/,' % destination_path
    tar_cmd = [
        'tar', '-pcvf', tar_path, '--transform',
        txfrm_regex, '.'
    ]

    run_command('tar_runtime_package', tar_cmd, cmd_cwd=app_dir)

    u_blob = open(tar_path, 'r').read()
    # We use gzip for performance instead of python's zip.
    gzip_cmd = ['gzip', tar_path, '-1']
    run_command('gzip_tar_runtime_package', gzip_cmd)
    return open(os.path.join(tar_dir, tar_path + '.gz'), 'rb').read(), u_blob


def has_pkg_descriptor(descriptor_files, ctx):
    for f in descriptor_files:
        if ctx.Contains(f):
            return True
    return False


def all_descriptor_contents(descriptor_files, ctx):
    descriptor = None
    descriptor_contents = ""
    for f in descriptor_files:
        if ctx.Contains(f):
            descriptor = f
            descriptor_contents += ctx.GetFile(descriptor)
            break
    if not descriptor:
        logging.info("No package descriptor found. No packages installed.")
        return None
    return descriptor_contents


def descriptor_parser(descriptor_files, ctx):
    descriptor = None
    for f in descriptor_files:
        if ctx.Contains(f):
            descriptor = f
            descriptor_contents = ctx.GetFile(descriptor)
            logging.info("descriptor_contents:\n%s", descriptor_contents)
            break
    if not descriptor:
        logging.info("No package descriptor found. No packages installed.")
        return None
    return descriptor_contents


def descriptor_copy(ctx, descriptor_files, app_dir):
    for f in descriptor_files:
        if ctx.Contains(f):
            with open(os.path.join(app_dir, f), 'w') as w:
                w.write(ctx.GetFile(f))


#  Return minimum ttl if the descriptor file has unspecified deps
def get_ttl(descriptor_files, ctx):
    for f in descriptor_files:
        if ctx.Contains(f):
            if f in constants.UNSPECIFIED_DEPS_FILES:
                return constants.MINIMUM_TTL_WEEKS
            return constants.DEFAULT_TTL_WEEKS
    return constants.DEFAULT_TTL_WEEKS


def gen_tmp_dir(dirr):
    tmp_dir = tempfile.mkdtemp()
    dir_name = os.path.join(tmp_dir, dirr)
    os.mkdir(dir_name)
    return dir_name


def creation_time(image):
    logging.info(image.config_file())
    cfg = json.loads(image.config_file())
    return cfg.get('created')


def timestamp_to_time(dt_str):
    dt = dt_str.rstrip('Z')
    return datetime.datetime.strptime(dt, "%Y-%m-%dT%H:%M:%S")


def generate_overrides(set_env, venv_dir=constants.VENV_DIR):
    created_time = datetime.datetime.now().strftime('%Y-%m-%dT%H:') + '00:00Z'
    overrides_dct = {
        'created': created_time,
    }
    if set_env:
        env = {
            'VIRTUAL_ENV': venv_dir,
        }
        path_dir = os.path.join(venv_dir, "bin")
        env['PATH'] = '%s:$PATH' % path_dir
        overrides_dct['env'] = venv_dir
    return overrides_dct


def parseCacheLogEntry(entry):
    """
    This takes an FTL log entry and parses out relevant caching information
    It returns a map with the information parsed from the entry

    Example entry (truncated for line length):
        INFO     [CACHE][MISS] v1:PYTHON:click:==6.7->f1ea...

    Return value for this entry:
        {
            "key_version": "v1",
            "language": "python",
            "phase": 2,
            "package": "click",
            "version": "6.7",
            "key": "f1ea...",
            "hit": True
        }
    """
    if "->" not in entry or "[CACHE]" not in entry:
        logging.warn("cannot parse non-cache log entry %s" % entry)
        return None

    entry = entry.rstrip("\n").lstrip("INFO").lstrip(" ").lstrip("[CACHE]")
    hit = True if entry.startswith("[HIT]") else False
    entry = entry.lstrip("[HIT]").lstrip("[MISS]").lstrip(" ")

    parts = entry.split("->")[0]
    key = entry.split("->")[1]
    parts = parts.split(":")
    if len(parts) == 2:
        # phase 1 entry
        return {
            "key_version": parts[0],
            "language": parts[1],
            "phase": 1,
            "key": key,
            "hit": hit
        }
    else:
        # phase 2 entry
        return {
            "key_version": parts[0],
            "language": parts[1],
            "phase": 2,
            "package": parts[2],
            "version": parts[3],
            "key": key,
            "hit": hit
        }


def run_command(cmd_name,
                cmd_args,
                cmd_cwd=None,
                cmd_env=None,
                cmd_input=None,
                err_type=ftl_error.FTLErrors.INTERNAL()):
    with Timing(cmd_name):
        cmd = "%s %s" % (cmd_name, " ".join(cmd_args))
        logging.info(cmd)
        proc_pipe = None
        try:
            proc_pipe = subprocess.Popen(
                cmd_args,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                cwd=cmd_cwd,
                env=cmd_env,
            )
        except OSError as e:
            raise ftl_error.InternalError(
                "%s\nexited with error %s\n%s is likely not on the path" %
                (cmd, e, cmd_name))
        stdout, stderr = proc_pipe.communicate(input=cmd_input)
        logging.info("`%s` stdout:\n%s", cmd_name, stdout)
        err_txt = ""
        if stderr:
            err_txt = "`%s` had stderr output:\n%s" % (cmd_name, stderr)
            logging.info(err_txt)
        if proc_pipe.returncode:
            ret_txt = "error: `%s` returned code: %d" % (cmd_name,
                                                         proc_pipe.returncode)
            logging.error(ret_txt)
            if err_type == ftl_error.FTLErrors.USER():
                raise ftl_error.UserError("%s\n%s" % (err_txt, ret_txt))
            elif err_type == ftl_error.FTLErrors.INTERNAL():
                raise ftl_error.InternalError("%s\n%s" % (err_txt, ret_txt))
            else:
                raise Exception("Unknown error type passed to run_command")
