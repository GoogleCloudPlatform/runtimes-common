load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.16.3",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

http_archive(
    name = "bazel_gazelle",
    sha256 = "92a3c59734dad2ef85dc731dbcb2bc23c4568cded79d4b87ebccd787eb89e8d0",
    urls = [("https://github.com/bazelbuild/bazel-gazelle/releases/download/0.11.0/bazel-gazelle-0.11.0.tar.gz")],
)


git_repository(
    name = "subpar",
    remote = "https://github.com/google/subpar",
    tag = "1.0.0",
)

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

git_repository(
    name = "io_bazel_rules_docker",
    remote = "https://github.com/bazelbuild/rules_docker.git",
    commit = "5eb0728594013d746959c4bd21aa4b0c3e3848d8",
)

git_repository(
    name = "io_bazel_rules_python",
    commit = "115e3a0dab4291184fdcb0d4e564a0328364571a",
    remote = "https://github.com/bazelbuild/rules_python.git",
)

load("@io_bazel_rules_python//python:pip.bzl", "pip_repositories", "pip_import")

pip_repositories()

pip_import(
    name = "requests",
    requirements = "//ftl:requirements.txt",
)

load("@requests//:requirements.bzl", "pip_install")

pip_install()

git_repository(
    name = "containerregistry",
    commit = "a6f2d42895f3017e244fa6ff18e1dc1d797aebb8",
    remote = "https://github.com/google/containerregistry",
)

load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_repositories",
    "docker_pull",
)

docker_repositories()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "repositories",
)

repositories()

http_archive(
    name = "mock",
    build_file_content = """
# Rename mock.py to __init__.py
genrule(
    name = "rename",
    srcs = ["mock.py"],
    outs = ["__init__.py"],
    cmd = "cat $< >$@",
)
py_library(
   name = "mock",
   srcs = [":__init__.py"],
   visibility = ["//visibility:public"],
)""",
    sha256 = "b839dd2d9c117c701430c149956918a423a9863b48b09c90e30a6013e7d2f44f",
    strip_prefix = "mock-1.0.1/",
    type = "tar.gz",
    urls = [("https://pypi.python.org/packages/source/m/mock/mock-1.0.1.tar.gz")],
)

docker_pull(
    name = "python_base",
    digest = "sha256:163a514abdb54f99ba371125e884c612e30d6944628dd6c73b0feca7d31d2fb3",
    registry = "gcr.io",
    repository = "google-appengine/python",
)

http_archive(
    name = "docker_credential_gcr",
    build_file_content = """package(default_visibility = ["//visibility:public"])
exports_files(["docker-credential-gcr"])""",
    sha256 = "c4f51ff78c25e2bfef38af0f38c6966806e25da7c5e43092c53a4d467fea4743",
    type = "tar.gz",
    urls = [("https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v1.4.1/docker-credential-gcr_linux_amd64-1.4.1.tar.gz")],
)

# TODO(aaron-prindle) cleanup circular dep here by pushing ubuntu_base to GCR
# OR by moving structure_test to own repo

git_repository(
    name = "base_images_docker",
    commit = "ac87be384d4e321a14aa9c11b3383a0f374511d3",
    remote = "https://github.com/GoogleCloudPlatform/base-images-docker.git",
)

http_file(
    name = "ubuntu_16_0_4_tar_download",
    sha256 = "51a8c466269bdebf232cac689aafad8feacd64804b13318c01096097a186d051",
    urls = [("https://storage.googleapis.com/ubuntu_tar/20171028/ubuntu-xenial-core-cloudimg-amd64-root.tar.gz")],
)

docker_pull(
    name = "node_base",
    digest = "sha256:fdc4906b4253e21663aa78805a3c3c45c6d8bc47596b4af953024f7a4ac64a71",
    registry = "gcr.io",
    repository = "gcp-runtimes/nodejs8_app_builder",
)

docker_pull(
    name = "distroless_base",
    digest = "sha256:4a8979a768c3ef8d0a8ed8d0af43dc5920be45a51749a9c611d178240f136eb4",
    registry = "gcr.io",
    repository = "distroless/base",
)

docker_pull(
    name = "php_base",
    digest = "sha256:194817f8f35b3cf26a706e6825994845c98a4745d61d6dea30fee976ed5aac0d",
    registry = "gcr.io",
    repository = "gae-runtimes/php72_app_builder",
)
