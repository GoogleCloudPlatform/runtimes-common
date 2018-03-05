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

from ftl.common import builder
from ftl.common import constants
from ftl.common import ftl_util
from ftl.common import layer_builder as base_builder
from ftl.node import layer_builder as node_builder


class Node(builder.RuntimeBase):
    def __init__(self, ctx, args, cache_version_str):
        super(
            Node, self).__init__(
            ctx, constants.NODE_NAMESPACE, args, cache_version_str, [
                constants.PACKAGE_LOCK, constants.PACKAGE_JSON,
                constants.NPMRC])

    def Build(self):
        lyr_imgs = []
        lyr_imgs.append(self._base_image)
        if ftl_util.has_pkg_descriptor(self._descriptor_files, self._ctx):
            layer_builder = node_builder.LayerBuilder(
                ctx=self._ctx,
                descriptor_files=self._descriptor_files,
                destination_path=self._args.destination_path)
            cached_pkg_img = None
            with ftl_util.Timing("checking cached pkg layer"):
                key = layer_builder.GetCacheKey()
                cached_pkg_img = self._cache.Get(key)
            if cached_pkg_img is not None:
                layer_builder.SetImage(cached_pkg_img)
            else:
                with ftl_util.Timing("building pkg layer"):
                    layer_builder.BuildLayer()
                with ftl_util.Timing("uploading pkg layer"):
                    self._cache.Set(layer_builder.GetCacheKey(),
                                    layer_builder.GetImage())
            lyr_imgs.append(layer_builder)

        app = base_builder.AppLayerBuilder(self._ctx,
                                           self._args.destination_path,
                                           self._args.entrypoint,
                                           self._args.exposed_ports)
        with ftl_util.Timing("builder app layer"):
            app.BuildLayer()
        lyr_imgs.append(app)
        with ftl_util.Timing("stitching lyrs into final image"):
            ftl_image = self.AppendLayersIntoImage(lyr_imgs)
        with ftl_util.Timing("uploading final image"):
            self.StoreImage(ftl_image)
