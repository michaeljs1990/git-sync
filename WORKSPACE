load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.1/rules_go-v0.22.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.1/rules_go-v0.22.1.tar.gz",
    ],
    sha256 = "e6a6c016b0663e06fa5fccf1cd8152eab8aa8180c583ec20c872f4f9953a7ac5",
)

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains(go_version = "1.14")

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# gazelle:resolve go github.com/google/go-github @com_github_google_go_github//v21/github:go_default_library
go_repository(
    name = "com_github_google_go_github",
    importpath = "github.com/google/go-github",
    sum = "h1:N0LgJ1j65A7kfXrZnUDaYCs/Sf4rEjNlfyDHW9dolSY=",
    version = "v17.0.0+incompatible",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    importpath = "github.com/prometheus/client_golang",
    sum = "h1:awm861/B8OKDd2I/6o1dy3ra4BamzKhYOiGItCeZ740=",
    version = "v0.9.2",
)

go_repository(
    name = "com_github_santhosh_tekuri_jsonschema",
    importpath = "github.com/santhosh-tekuri/jsonschema",
    sum = "h1:hNhW8e7t+H1vgY+1QeEQpveR6D4+OwKPXCfD2aieJis=",
    version = "v1.2.4",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:hI/7Q+DtNZ2kINb6qt/lS+IyXnHQe9e90POfeewL/ME=",
    version = "v1.3.0",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:u1CmMhe3a44hy8VIgpInORnI01UVaUYheqR7x9BxT3c=",
    version = "v0.0.0-20190122013713-64072686203f",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:pcBdqVcrlT+A3i+tWsOROFONQyey9tisIQHI4xqVGLg=",
    version = "v0.0.0-20190115181402-5dab4167f31c",
)

go_repository(
    name = "in_gopkg_src_d_go_git_v4",
    importpath = "gopkg.in/src-d/go-git.v4",
    sum = "h1:aAyBmkdE1QUUEHcP4YFCGKmsMQRAuRmUcPEQR7lOAa0=",
    version = "v4.8.1",
)

go_repository(
    name = "in_gopkg_src_d_go_billy_v4",
    importpath = "gopkg.in/src-d/go-billy.v4",
    sum = "h1:omN5CrMrMcQ+4I8bJ0wEhOBPanIRWzFC953IiXKdYzo=",
    version = "v4.2.1",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:ZCJp+EgiOT7lHqUV2J862kp8Qj64Jo6az82+3Td9dZw=",
    version = "v2.2.2",
)

go_repository(
    name = "com_github_alcortesm_tgz",
    importpath = "github.com/alcortesm/tgz",
    sum = "h1:uSoVVbwJiQipAclBbw+8quDsfcvFjOpI5iCf4p/cqCs=",
    version = "v0.0.0-20161220082320-9c5fe88206d7",
)

go_repository(
    name = "com_github_anmitsu_go_shlex",
    importpath = "github.com/anmitsu/go-shlex",
    sum = "h1:kFOfPq6dUM1hTo4JG6LR5AXSUEsOjtdm0kw0FtQtMJA=",
    version = "v0.0.0-20161002113705-648efa622239",
)

go_repository(
    name = "com_github_beorn7_perks",
    importpath = "github.com/beorn7/perks",
    sum = "h1:xJ4a3vCFaGF/jqvzLMYoU8P317H5OQ+Via4RmuPwCS0=",
    version = "v0.0.0-20180321164747-3a771d992973",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_emirpasic_gods",
    importpath = "github.com/emirpasic/gods",
    sum = "h1:rUF4PuzEjMChMiNsVjdI+SyLu7rEqpQ5reNFnhC7oFo=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_flynn_go_shlex",
    importpath = "github.com/flynn/go-shlex",
    sum = "h1:BHsljHzVlRcyQhjrss6TZTdY2VfCqZPbv5k3iBFa2ZQ=",
    version = "v0.0.0-20150515145356-3f9db97f8568",
)

go_repository(
    name = "com_github_gliderlabs_ssh",
    importpath = "github.com/gliderlabs/ssh",
    sum = "h1:j3L6gSLQalDETeEg/Jg0mGY0/y/N6zI2xX1978P0Uqw=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    sum = "h1:P3YflyNX/ehuJFLhxviNdFxQPkGK5cDcApsge1SqnvM=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:+dTQ8DZQJz0Mb/HjFlkptS1FeQ4cWSnN941F8aEG4SQ=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_google_go_github_v21",
    importpath = "github.com/google/go-github/v21",
    sum = "h1:tn4/tmCgPAsezJFwZcMnE7U0R9/AtKRBGX4s4LFdDzI=",
    version = "v21.0.0",
)

go_repository(
    name = "com_github_google_go_querystring",
    importpath = "github.com/google/go-querystring",
    sum = "h1:Xkwi/a1rcvNg1PPYe5vI8GbeBY/jrVuDX5ASuANWTrk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jbenet_go_context",
    importpath = "github.com/jbenet/go-context",
    sum = "h1:BQSFePA1RWJOlocH6Fxy8MmwDt+yVQYULKfN0RoTN8A=",
    version = "v0.0.0-20150711004518-d14ea06fba99",
)

go_repository(
    name = "com_github_jessevdk_go_flags",
    importpath = "github.com/jessevdk/go-flags",
    sum = "h1:4IU2WS7AumrZ/40jfhf4QVDMsQwqA7VEHozFRrGARJA=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_kevinburke_ssh_config",
    importpath = "github.com/kevinburke/ssh_config",
    sum = "h1:RgQk53JHp/Cjunrr1WlsXSZpqXn+uREuHvUVcK82CV8=",
    version = "v0.0.0-20180830205328-81db2a75821e",
)

go_repository(
    name = "com_github_konsorten_go_windows_terminal_sequences",
    importpath = "github.com/konsorten/go-windows-terminal-sequences",
    sum = "h1:mweAR1A6xJ3oS2pRaGiHgQ4OO8tzTaLawm8vnODuwDk=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_kr_pretty",
    importpath = "github.com/kr/pretty",
    sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_kr_pty",
    importpath = "github.com/kr/pty",
    sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_kr_text",
    importpath = "github.com/kr/text",
    sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
    sum = "h1:4hp9jkHxhMHkqkrB3Ix0jegS5sx/RkqARlsWZ6pIwiU=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    importpath = "github.com/mitchellh/go-homedir",
    sum = "h1:vKb8ShqSby24Yrqr/yDYkuFz8d0WUjys40rvnGC8aR0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_pelletier_go_buffruneio",
    importpath = "github.com/pelletier/go-buffruneio",
    sum = "h1:U4t4R6YkofJ5xHm3dJzuRpPZ0mr5MMCoAWooScCR7aA=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_pkg_errors",
    importpath = "github.com/pkg/errors",
    sum = "h1:WdK/asTD0HN+q6hsWO3/vpuAkAr+tw6aNJNDFFf0+qw=",
    version = "v0.8.0",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_prometheus_client_model",
    importpath = "github.com/prometheus/client_model",
    sum = "h1:idejC8f05m9MGOsuEi1ATq9shN03HrxNkD/luQvxCv8=",
    version = "v0.0.0-20180712105110-5c3871d89910",
)

go_repository(
    name = "com_github_prometheus_common",
    importpath = "github.com/prometheus/common",
    sum = "h1:PnBWHBf+6L0jOqq0gIVUe6Yk0/QMZ640k6NvkxcBf+8=",
    version = "v0.0.0-20181126121408-4724e9255275",
)

go_repository(
    name = "com_github_prometheus_procfs",
    importpath = "github.com/prometheus/procfs",
    sum = "h1:9a8MnZMP0X2nLJdBg+pBmGgkJlSaKC2KaQmTCk1XDtE=",
    version = "v0.0.0-20181204211112-1dc9a6cbc91a",
)

go_repository(
    name = "com_github_sergi_go_diff",
    importpath = "github.com/sergi/go-diff",
    sum = "h1:Kpca3qRNrduNnOQeazBd0ysaKrUJiIuISHxogkT9RPQ=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_src_d_gcfg",
    importpath = "github.com/src-d/gcfg",
    sum = "h1:xXbNR5AlLSA315x2UO+fTSSAXCDf+Ar38/6oyGbDKQ4=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:2vfRuCMp5sSVIDSqO8oNnWJq7mPa6KVP3iPIwFBuy8A=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:bSDNvY7ZPG5RlJ8otE/7V6gMiyenm9RtJ7IUVIAoJ1w=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_xanzy_ssh_agent",
    importpath = "github.com/xanzy/ssh-agent",
    sum = "h1:Adglfbi5p9Z0BmK2oKU9nTG+zKfniSfnaMYB+ULd+Ro=",
    version = "v0.2.0",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    sum = "h1:eOI3/cP2VTU6uZLDYAoic+eyzzB9YyGmJ7eIjl8rOPg=",
    version = "v0.34.0",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:qIbj1fsPNlZgppZ+VLlY7N33q108Sa+fhmuc+sWQYwY=",
    version = "v1.0.0-20180628173108-788fd7840127",
)

go_repository(
    name = "in_gopkg_src_d_go_git_fixtures_v3",
    importpath = "gopkg.in/src-d/go-git-fixtures.v3",
    sum = "h1:XWW/s5W18RaJpmo1l0IYGqXKuJITWRFuA45iOf1dKJs=",
    version = "v3.1.1",
)

go_repository(
    name = "in_gopkg_warnings_v0",
    importpath = "gopkg.in/warnings.v0",
    sum = "h1:wFXVbFY8DY5/xOe1ECiWdKCzZlxgshcYVNkBHstARME=",
    version = "v0.1.2",
)

go_repository(
    name = "org_golang_google_appengine",
    importpath = "google.golang.org/appengine",
    sum = "h1:/wp5JvzpHIxhs/dumFmF7BXTf3Z+dd4uXta4kVyO508=",
    version = "v1.4.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:bRhVy7zSSasaqNksaRZiA5EEI+Ei4I1nO5Jh72wfHlg=",
    version = "v0.0.0-20190108225652-1e06a53dbb7e",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:YUO/7uOKsKeq9UokNS62b8FYywz3ker1l1vDZRCRefw=",
    version = "v0.0.0-20181221193216-37e7f081c4d4",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:I6FyU15t786LL7oL/hn43zqTuEGr4PN7F4XJ1p4E3Y8=",
    version = "v0.0.0-20180905080454-ebe1bf3edb33",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:g61tztE5qeGQ89tm6NTjjM9VPIm088od1l6aSorWRWg=",
    version = "v0.3.0",
)
