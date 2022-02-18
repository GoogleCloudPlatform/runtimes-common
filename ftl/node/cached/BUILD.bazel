package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_docker//python:image.bzl",
    "py_image",
)

py_binary(
    name = "node_cached",
    srcs = [
        "main.py",
        "//ftl/cached:cached_lib",
    ],
    data = ["//ftl:node_builder.par"],
    main = "main.py",
    deps = [
        "//ftl:ftl_lib",
        "@containerregistry",
    ],
)

py_image(
    name = "node_cached_image",
    srcs = [
        "main.py",
        "//ftl/cached:cached_lib",
    ],
    base = "@node_base//image",
    main = "main.py",
    deps = [
        "//ftl:ftl_lib",
        "@containerregistry",
    ],
)
