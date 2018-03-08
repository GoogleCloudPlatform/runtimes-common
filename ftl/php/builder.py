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
import json

from ftl.common import builder
from ftl.common import ftl_util
from ftl.common import layer_builder as base_builder
from ftl.php import layer_builder as php_builder

_PHP_NAMESPACE = 'php-package-lock-cache'
_COMPOSER_LOCK = 'composer.lock'
_COMPOSER_JSON = 'composer.json'


class PHP(builder.RuntimeBase):
    def __init__(self, ctx, args):
        super(PHP, self).__init__(ctx, _PHP_NAMESPACE, args,
                                  [_COMPOSER_LOCK, _COMPOSER_JSON])

    def _parse_composer_pkgs(self):
        descriptor_contents = ftl_util.descriptor_parser(
            self._descriptor_files, self._ctx)
        composer_json = json.loads(descriptor_contents)
        pkgs = []
        for k, v in composer_json['require'].iteritems():
            pkgs.append((k, v))
        return pkgs

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base_image)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            pkgs = self._parse_composer_pkgs()
            # if there are 42 or more packages, revert to using phase 1
            if len(pkgs) > 41:
                pkgs = [None]
            for pkg_txt in pkgs:
                logging.info('Building package layer')
                logging.info(pkg_txt)
                cache = self._cache if self._args.cache else None
                layer_builder = php_builder.LayerBuilder(
                    ctx=self._ctx,
                    descriptor_files=self._descriptor_files,
                    pkg_descriptor=pkg_txt,
                    destination_path=self._args.destination_path,
                    cache=cache)
                layer_builder.BuildLayer()
                lyr_imgs.append(layer_builder)

        app = base_builder.AppLayerBuilder(
            ctx=self._ctx,
            destination_path=self._args.destination_path,
            entrypoint=self._args.entrypoint,
            exposed_ports=self._args.exposed_ports)
        app.BuildLayer()
        lyr_imgs.append(app)
        ftl_image = ftl_util.AppendLayersIntoImage(lyr_imgs)
        self.StoreImage(ftl_image)
