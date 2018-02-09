http_archive(
    name = "io_bazel_rules_go",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.9.0/rules_go-0.9.0.tar.gz",
    sha256 = "4d8d6244320dd751590f9100cf39fd7a4b75cd901e1f3ffdfd6f048328883695",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

git_repository(
    name = "io_bazel_rules_docker",
    commit = "8bbe2a8abd382641e65ff7127a3700a8530f02ce",
    remote = "https://github.com/bazelbuild/rules_docker.git",
)

git_repository(
    name = "containerregistry",
    commit = "6b250f0bae8cce028df939010ee3118c8f2977ba",
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

new_http_archive(
    name = "retrying",
    build_file_content = """
# Rename retrying.py to __init__.py
genrule(
    name = "rename",
    srcs = ["retrying.py"],
    outs = ["__init__.py"],
    cmd = "cat $< >$@",
)
py_library(
   name = "retrying",
   srcs = [":__init__.py"],
   visibility = ["//visibility:public"],
)""",
    # sha256 = "b839dd2d9c117c701430c149956918a423a9863b48b09c90e30a6013e7d2f44f",
    strip_prefix = "retrying-1.3.3/",
    type = "tar.gz",
    url = "https://pypi.python.org/packages/44/ef/beae4b4ef80902f22e3af073397f079c96969c69b2c7d52a57ea9ae61c9d/retrying-1.3.3.tar.gz",
)

docker_pull(
    name = "python_base",
    digest = "sha256:163a514abdb54f99ba371125e884c612e30d6944628dd6c73b0feca7d31d2fb3",
    registry = "gcr.io",
    repository = "google-appengine/python",
)

new_http_archive(
    name = "docker_credential_gcr",
    sha256 = "c4f51ff78c25e2bfef38af0f38c6966806e25da7c5e43092c53a4d467fea4743",
    type = "tar.gz",
    build_file_content = """package(default_visibility = ["//visibility:public"])
exports_files(["docker-credential-gcr"])""",
    url = "https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v1.4.1/docker-credential-gcr_linux_amd64-1.4.1.tar.gz",
)

new_http_archive(
    name = "gsutil",
    # sha256 = "",
    type = "tar.gz",
    build_file_content = """package(default_visibility = ["//visibility:public"])
exports_files(["gsutil"])""",
    url = "https://storage.googleapis.com/pub/gsutil.tar.gz",
)

# TODO(aaron-prindle) cleanup circular dep here by pushing ubuntu_base to GCR
# OR by moving structure_test to own repo

git_repository(
    name = "base_images_docker",
    commit = "ac87be384d4e321a14aa9c11b3383a0f374511d3",
    remote = "https://github.com/GoogleCloudPlatform/base-images-docker.git",
)

UBUNTU_MAP = {
    "16_0_4": {
        "sha256": "51a8c466269bdebf232cac689aafad8feacd64804b13318c01096097a186d051",
        "url": "https://storage.googleapis.com/ubuntu_tar/20171028/ubuntu-xenial-core-cloudimg-amd64-root.tar.gz",
    },
}

[http_file(
    name = "ubuntu_%s_tar_download" % version,
    sha256 = map["sha256"],
    url = map["url"],
) for version, map in UBUNTU_MAP.items()]

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
docker_pull(
    name = "php_base",
    digest = "sha256:b4a1f5de8156f30ea1a6e6f84afb7ea79013a57d0cae4a530d4806df4a04a1e3",
    registry = "gcr.io",
    repository = "gae-runtimes/php72_app_builder"
)
