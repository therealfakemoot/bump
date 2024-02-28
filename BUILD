load("@rules_go//go:def.bzl", "go_binary", "go_library")
load("@gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/therealfakemoot/bump
gazelle(name = "gazelle")

go_library(
    name = "bump_lib",
    srcs = ["main.go"],
    importpath = "github.com/therealfakemoot/bump",
    visibility = ["//visibility:private"],
    deps = ["//cmd"],
)

go_binary(
    name = "bump",
    embed = [":bump_lib"],
    visibility = ["//visibility:public"],
)
