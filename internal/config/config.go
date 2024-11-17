package config

import (
	"github.com/y4ney/cloud-native-security-vuln/internal/component"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var (
	Token      = "write-your-github-token"
	CacheDir   = "cloud-native-sec-vuln"
	Components = []*component.Component{
		{
			Owner: "moby",
			Repo:  "buildkit",
		},
		{
			Owner: "opencontainers",
			Repo:  "runc",
		},
		{
			Owner: "containerd",
			Repo:  "containerd",
		},
		{
			Owner: "cri-o",
			Repo:  "cri-o",
		},
		{
			Owner: "google",
			Repo:  "gvisor",
		},
		{
			Owner: "inclavare-containers",
			Repo:  "inclavare-containers",
		},
		{
			Owner: "openeuler-mirror",
			Repo:  "iSulad",
		},
		{
			Owner: "kata-containers",
			Repo:  "kata-containers",
		},
		{
			Owner: "krustlet",
			Repo:  "krustlet",
		},
		{
			Owner: "kuasar-io",
			Repo:  "kuasar",
		},
		{
			Owner: "lima-vm",
			Repo:  "lima",
		},
		{
			Owner: "lxc",
			Repo:  "incus",
		},
		{
			Owner: "rkt",
			Repo:  "rkt",
		},
		{
			Owner: "apptainer",
			Repo:  "singularity",
		},
		{
			Owner: "TritonDataCenter",
			Repo:  "smartos-live",
		},
		{
			Owner: "openeuler-mirror",
			Repo:  "stratovirt",
		},
		{
			Owner: "nestybox",
			Repo:  "sysbox",
		},
		{
			Owner: "virtual-kubelet",
			Repo:  "virtual-kubelet",
		},
		{
			Owner: "WasmEdge",
			Repo:  "WasmEdge",
		},
		{
			Owner: "youki-dev",
			Repo:  "youki",
		},
	}
)

type Config struct {
	Token      string                 `yaml:"token,omitempty"`
	CacheDir   string                 `yaml:"cache_dir,omitempty"`
	Components []*component.Component `yaml:"components,omitempty"`
}
type ConfFunc func(*Config)

// WithToken 配置 Token
func WithToken(token string) ConfFunc {
	return func(c *Config) { c.Token = token }
}
func WithCacheDir(cacheDir string) ConfFunc {
	return func(c *Config) { c.CacheDir = cacheDir }
}

func WithComponent(component []*component.Component) ConfFunc {
	return func(c *Config) { c.Components = component }
}

func NewConfig(opts ...ConfFunc) *Config {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	var conf = &Config{Token: Token, Components: Components, CacheDir: filepath.Join(cacheDir, CacheDir)}
	for _, opt := range opts {
		opt(conf)
	}
	return conf
}

func (c *Config) Generate(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
func LoadConfFile(filename string) (*Config, error) {
	var config *Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, xerrors.New("file is empty")
	}
	return config, yaml.Unmarshal(data, &config)
}
