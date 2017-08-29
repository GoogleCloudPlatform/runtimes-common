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
    vcs = "git",
    tag = "v2",
)

new_go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    tag = "master",
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "db1b348dfdf161a784bc1efc5a1020395572b996",
    remote = "https://github.com/bazelbuild/rules_docker.git",

)

git_repository(
    name = "containerregistry",
    commit = "b0278a1544238d03648861b6d9395414d4c958e5",
    remote = "https://github.com/google/containerregistry",
)

load(
  "@io_bazel_rules_docker//docker:docker.bzl",
  "docker_repositories"
)
docker_repositories()
