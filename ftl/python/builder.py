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
"""This package defines the interface for orchestrating image builds."""

import json
import concurrent.futures

from ftl.common import builder
from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import layer_builder as base_builder

from ftl.python import layer_builder as package_builder
from ftl.python import python_util


class Python(builder.RuntimeBase):
    def __init__(self, ctx, args):
        super(Python, self).__init__(
            ctx,
            constants.PYTHON_CACHE_NAMESPACE,
            args,
            [
                constants.PIPFILE_LOCK,
                constants.PIPFILE,  # not supported rn
                constants.REQUIREMENTS_TXT
            ])
        self._venv_dir = args.venv_dir
        self._wheel_dir = ftl_util.gen_tmp_dir(constants.WHEEL_DIR)

        self._python_cmd = args.python_cmd.split(" ")
        self._pip_cmd = args.pip_cmd.split(" ")
        self._venv_cmd = args.venv_cmd.split(" ")

        self._is_phase2 = ctx.Contains(constants.PIPFILE_LOCK)

    def _parse_pipfile_pkgs(self):

        pkg_descriptor = ftl_util.descriptor_parser(self._descriptor_files,
                                                    self._ctx)
        pipfile_json = json.loads(pkg_descriptor)
        pkgs = []
        for pkg, info in pipfile_json['default'].iteritems():
            version = info['version']
            pkgs.append((pkg, version))
        return pkgs

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base_image)

        interpreter_builder = package_builder.InterpreterLayerBuilder(
            venv_dir=self._venv_dir,
            python_cmd=self._python_cmd,
            venv_cmd=self._venv_cmd,
            cache=self._cache)
        interpreter_builder.BuildLayer()
        lyr_imgs.append(interpreter_builder.GetImage())

        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            # build interpreter layer

            if self._is_phase2:
                # do a phase 2 build of the package layers w/ Pipfile.lock
                # iterate over package/version Pipfile.lock
                python_util.setup_venv(self._venv_dir,
                                       self._venv_cmd,
                                       self._python_cmd)
                pkgs = self._parse_pipfile_pkgs()
                with ftl_util.Timing('uploading_all_package_layers'):
                    with concurrent.futures.ThreadPoolExecutor(
                            max_workers=constants.THREADS) as executor:
                        future_to_params = {executor.submit(
                                self._build_pkg, pkg,
                                interpreter_builder, lyr_imgs): pkg
                                for pkg in pkgs
                        }
                        for future in concurrent.futures.as_completed(
                                future_to_params):
                            future.result()
            else:
                # do a phase 1 build of the package layers w/ requirements.txt
                req_txt_builder = package_builder.RequirementsLayerBuilder(
                    ctx=self._ctx,
                    descriptor_files=self._descriptor_files,
                    directory=self._args.directory,
                    pkg_dir=None,
                    wheel_dir=self._wheel_dir,
                    venv_dir=self._venv_dir,
                    python_cmd=self._python_cmd,
                    pip_cmd=self._pip_cmd,
                    venv_cmd=self._venv_cmd,
                    dep_img_lyr=interpreter_builder,
                    cache=self._cache)
                req_txt_builder.BuildLayer()
                if req_txt_builder.GetImage():
                    lyr_imgs.append(req_txt_builder.GetImage())

        app = base_builder.AppLayerBuilder(
            directory=self._args.directory,
            destination_path=self._args.destination_path,
            entrypoint=self._args.entrypoint,
            exposed_ports=self._args.exposed_ports)
        app.BuildLayer()
        lyr_imgs.append(app.GetImage())
        additional_directory = base_builder.AppLayerBuilder(
            directory=self._args.additional_directory,
            destination_path=self._args.additional_directory,
            entrypoint=self._args.entrypoint,
            exposed_ports=self._args.exposed_ports)
        additional_directory.BuildLayer()
        lyr_imgs.append(additional_directory.GetImage())
        ftl_image = ftl_util.AppendLayersIntoImage(lyr_imgs)
        self.StoreImage(ftl_image)

    def _build_pkg(self, pkg, interpreter_builder, lyr_imgs):
        pipfile_builder = package_builder.PipfileLayerBuilder(
            ctx=self._ctx,
            descriptor_files=self._descriptor_files,
            directory=self._args.directory,
            pkg_descriptor=pkg,
            pkg_dir=None,
            wheel_dir=ftl_util.gen_tmp_dir(constants.WHEEL_DIR),
            venv_dir=self._venv_dir,
            pip_cmd=self._pip_cmd,
            venv_cmd=self._venv_cmd,
            dep_img_lyr=interpreter_builder,
            cache=self._cache)
        pipfile_builder.BuildLayer()
        lyr_imgs.append(pipfile_builder.GetImage())
