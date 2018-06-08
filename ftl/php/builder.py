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

import logging

from ftl.common import builder
from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import layer_builder as base_builder
from ftl.php import layer_builder as php_builder


class PHP(builder.RuntimeBase):
    def __init__(self, ctx, args):
        super(PHP, self).__init__(
            ctx, constants.PHP_CACHE_NAMESPACE, args,
            [constants.COMPOSER_LOCK, constants.COMPOSER_JSON])

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base_image)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            logging.info('Building package layer')
            layer_builder = php_builder.PhaseOneLayerBuilder(
                ctx=self._ctx,
                descriptor_files=self._descriptor_files,
                destination_path=self._args.destination_path,
                directory=self._args.directory,
                cache=self._cache)
            layer_builder.BuildLayer()
            lyr_imgs.append(layer_builder.GetImage())

        app = base_builder.AppLayerBuilder(
            directory=self._args.directory,
            destination_path=self._args.destination_path,
            entrypoint=self._args.entrypoint,
            exposed_ports=self._args.exposed_ports)
        app.BuildLayer()
        lyr_imgs.append(app.GetImage())
        ftl_image = ftl_util.AppendLayersIntoImage(lyr_imgs)
        self.StoreImage(ftl_image)
