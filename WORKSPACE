git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.4.2",
)

load(
    "@io_bazel_rules_go//go:def.bzl",
    "new_go_repository",
    "go_repositories",
)

go_repositories()

load("@io_bazel_rules_go//proto:go_proto_library.bzl", "go_proto_repositories")

go_proto_repositories()

new_go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    remote = "https://github.com/go-yaml/yaml",
    tag = "v2",
    vcs = "git",
)

new_go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    tag = "master",
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "65df68f4f64e9c59eb571290eb86bf07766393b6",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

git_repository(
    name = "containerregistry",
    commit = "b0278a1544238d03648861b6d9395414d4c958e5",
    remote = "https://github.com/google/containerregistry",
)

load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_repositories",
    "docker_pull",
)

docker_repositories()

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

docker_pull(
    name = "node_base",
    digest = "sha256:f98878fe17ac9474f5a4beb9f692272f698a9ce2dc1e6297d449b2003cfec3e9",
    registry = "gcr.io",
    repository = "google-appengine/nodejs",
)
