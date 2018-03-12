package(default_visibility = ["//visibility:public"])

# Converts template files into valid Go files with a single
# exported constant whose value is the template string.
# These Go files are under Go package "templates".
# Each template can be then referenced direclty in Go code
# such as templates.Readme.
genrule(
    name = "gen_readme_template",
    srcs = ["README.md.tmpl"],
    outs = ["readme.go"],
    # sed is used to escape back ticks.
    cmd = """
    echo -e "package templates\n\nconst Readme = \``cat $(SRCS) | sed 's/\`/\` + \"\`\" + \`/g'`\`" > $@""",
)
