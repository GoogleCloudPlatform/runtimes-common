load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
  name = "go_default_library",
  srcs = glob(["*.go"]),
  deps = [
    "@com_github_spf13_cobra//:go_default_library",
    "//ctc_lib:go_default_library",
  ],
  visibility = ["//visibility:public"],
)
