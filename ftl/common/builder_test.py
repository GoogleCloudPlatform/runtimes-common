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

import cStringIO
import os
import unittest
import tarfile
import tempfile

import layer_builder


def gen_tmp_dir(dirr):
    tmp_dir = tempfile.mkdtemp()
    dir_name = os.path.join(tmp_dir, dirr)
    os.mkdir(dir_name)
    return dir_name


class JustAppTest(unittest.TestCase):
    def test_build_app_layer(self):
        # All the files in the context should be added to the layer.
        tmp_dir = gen_tmp_dir("justapptest")
        files = {
            'foo': 'foo_contents',
            'bar': 'bar_contents',
            'baz/bat': 'bat_contents'
        }
        print "AGHHHHH"
        print tmp_dir
        for name, contents in files.iteritems():
            print name, contents
            path_lst = name.split("/")
            for i in range(len(path_lst)):
                if i == len(path_lst) - 1:
                    break
                os.mkdir(os.path.join(tmp_dir, path_lst[i]))
            with open(os.path.join(tmp_dir, name), "w") as f:
                f.write(contents)

        app_builder = layer_builder.AppLayerBuilder(tmp_dir)
        app_builder.BuildLayer()
        app_layer = app_builder.GetImage().GetFirstBlob()
        stream = cStringIO.StringIO(app_layer)
        with tarfile.open(fileobj=stream, mode='r:gz') as tf:
            # two additional files in a real directory . and the 'baz' dir
            # ['srv/.', 'srv/./foo', 'srv/./baz', 'srv/./baz/bat', 'srv/./bar']
            self.assertEqual(len(tf.getnames()), len(files) + 2)
            for p, f in files.iteritems():
                tar_path = os.path.join('srv/.', p)
                self.assertEquals(tf.extractfile(tar_path).read(), f)


if __name__ == '__main__':
    unittest.main()
