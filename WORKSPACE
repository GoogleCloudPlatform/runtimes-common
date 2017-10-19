git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.5.5",
)

load(
    "@io_bazel_rules_go//go:def.bzl",
    "new_go_repository",
    "go_repositories",
)

go_repositories()

load("@io_bazel_rules_go//proto:go_proto_library.bzl", "go_proto_repositories")

go_proto_repositories()

git_repository(
    name = "io_bazel_rules_docker",
    commit = "8bbe2a8abd382641e65ff7127a3700a8530f02ce",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

git_repository(
    name = "containerregistry",
    commit = "009ca89e9616c2f68155cf9c5fc6cbbb34aff3a0",
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
    "repositories"
)

repositories()

new_http_archive(
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
    url = "https://pypi.python.org/packages/source/m/mock/mock-1.0.1.tar.gz",
)

docker_pull(
    name = "python_base",
    digest = "sha256:163a514abdb54f99ba371125e884c612e30d6944628dd6c73b0feca7d31d2fb3",
    registry = "gcr.io",
    repository = "google-appengine/python",
)

http_file(
    name = "docker_credential_gcr",
    sha256 = "c4f51ff78c25e2bfef38af0f38c6966806e25da7c5e43092c53a4d467fea4743",
    url = "https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v1.4.1/docker-credential-gcr_linux_amd64-1.4.1.tar.gz",
)

# TODO(aaron-prindle) cleanup circular dep here by pushing ubuntu_base to GCR
# OR by moving structure_test to own repo

git_repository(
    name = "debian_docker",
    commit = "72cd6607fb809d8eda3d4b88c8b94f84c2921372",
    remote = "https://github.com/GoogleCloudPlatform/debian-docker.git",
)

http_file(
    name = "ubuntu_tar_download",
    sha256 = "e43f802f876505c0cee1759fbc545f65b2d59c4dd33835e93afea3fa124b5799",
    url = "https://partner-images.canonical.com/core/xenial/20171006/ubuntu-xenial-core-cloudimg-amd64-root.tar.gz",
)

docker_pull(
    name = "node_base",
    digest = "sha256:f98878fe17ac9474f5a4beb9f692272f698a9ce2dc1e6297d449b2003cfec3e9",
    registry = "gcr.io",
    repository = "google-appengine/nodejs",
)

docker_pull(
    name = "distroless_base",
    digest = "sha256:4a8979a768c3ef8d0a8ed8d0af43dc5920be45a51749a9c611d178240f136eb4",
    registry = "gcr.io",
    repository = "distroless/base"
)
