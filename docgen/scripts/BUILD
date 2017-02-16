load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
)

go_binary(
    name = "docgen",
    srcs = ["docgen/docgen.go"],
    deps = [
        "//lib/proto:go_default_library",
        "//lib/render:go_default_library",
        "//lib/spec:go_default_library",
        "@com_github_ghodss_yaml//:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@in_gopkg_yaml_v2//:go_default_library",
    ],
)
