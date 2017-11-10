package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//python:image.bzl",
    "py_image",
)

load(
    "@io_bazel_rules_docker//docker:docker.bzl",
    "docker_build",
)
exports_files(["requirements.txt", "Dockerfile"])

py_library(
    name = "benchmark_lib",
    srcs = glob(["*.py"]),
)

py_binary(
    name = "node_benchmark",
    srcs = ["//ftl:node_lib", ":benchmark_lib"],
    main = "ftl/node/benchmark/main.py",
    deps = [
        "@containerregistry",
    ],
    data = ["//ftl:node_builder.par"]
)

load("@base_images_docker//dockerfile_build:dockerfile_build.bzl", "dockerfile_build")

dockerfile_build(
    name = "ftl_benchmark_base",
    base = "@python_base//image",
    dockerfile = "//ftl/node/benchmark:Dockerfile",
)

docker_build(
    name = "node_benchmark_base",
    base = ":ftl_benchmark_base",
    cmd = [""],
    entrypoint = "",
    tars = [
        "//ftl:cred_helper.tar",
        "//ftl:default_creds.tar",
    ],
)

py_image(
    name = "node_benchmark_image",
    srcs = ["//ftl:node_lib", ":benchmark_lib"],
    base = ":node_benchmark_base",
    main = "ftl/node/benchmark/main.py",
    deps = [
        "@containerregistry",
    ],
    data = ["//ftl:node_builder.par"],
)

