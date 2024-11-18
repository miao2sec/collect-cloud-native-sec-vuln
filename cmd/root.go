package cmd

/*
Copyright © 2024 Yaney yangli.yaney@foxmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/miao2sec/cloud-native-security-vuln/internal/config"
	"github.com/miao2sec/cloud-native-security-vuln/internal/github"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
	"os"

	"github.com/spf13/cobra"
)

var (
	genConf  bool
	run      bool
	cacheDir string
	cfgFile  string
)

var rootCmd = &cobra.Command{
	Use:   "collect",
	Short: "收集云原生安全漏洞",
	Long: `
# 编译
go build -o collect

# 生成默认的配置文件
./collect -g

# 使用默认配置进行漏洞收集
./collect -r

## 自定义漏洞数据的缓存路径
./collect -r -c path/to/cacheDir

## 使用自定义的配置文件进行漏洞收集
./collect -r -f <config>.yaml

`,
	Run: Run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	InitLogger()
	rootCmd.Flags().BoolVarP(&genConf, "gen-config", "g", false, "generate default config file")

	rootCmd.Flags().BoolVarP(&run, "run", "r", false, "begin to collect cloud native security vulnerabilities")
	rootCmd.Flags().StringVarP(&cacheDir, "cache-dir", "c", "", "specify the cache directory")
	rootCmd.Flags().StringVarP(&cfgFile, "cfg-file", "f", "", "specify the config file")

}

func InitLogger() {
	var (
		defaultLogger = zerolog.New(os.Stderr)
		logLevel      = zerolog.TraceLevel
	)

	zerolog.SetGlobalLevel(logLevel)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		defaultLogger = zerolog.New(zerolog.NewConsoleWriter())
	}
	log.Logger = defaultLogger.With().Timestamp().Stack().Logger()
}

func Run(cmd *cobra.Command, args []string) {
	if genConf {
		genCfgFile()
		return
	}
	var (
		conf *config.Config
		err  error
	)

	if cacheDir != "" {
		conf = config.NewConfig(config.WithCacheDir(cacheDir))
	} else {
		conf = config.NewConfig()
	}
	if cfgFile != "" {
		conf, err = config.LoadConfFile(cfgFile)
		if err != nil {
			log.Logger.Fatal().Str("config file", cfgFile).Err(err).Msg("failed to load config file")
		}
	}
	collect(conf)
	return
}

func genCfgFile() {
	var (
		conf           = config.NewConfig()
		defaultCfgFile = "config.yaml"
	)

	if err := conf.Generate(defaultCfgFile); err != nil {
		log.Logger.Fatal().Err(err).Str("file", defaultCfgFile).Msg("failed to generate config file")
	}
	log.Logger.Info().Str("file", defaultCfgFile).Msg("success to generate config file")
}

func collect(conf *config.Config) {
	var (
		client = github.NewClient(github.WithToken(conf.Token))
		err    error
	)

	for _, component := range conf.Components {
		component.Advisories, err = client.GetAdvisories(component)
		if err != nil {
			log.Logger.Fatal().Str("component", component.Repo).Err(err).Msg("failed to get advisories")
		}
		if len(component.Advisories) != 0 {
			log.Logger.Info().Str("component", component.Repo).Int("count", len(component.Advisories)).Msg("have vuln(s)")
		} else {
			log.Logger.Debug().Str("Repo", component.Repo).Msg("don't have vuln")
			continue
		}
		if err = component.Save(conf.CacheDir); err != nil {
			log.Logger.Fatal().Str("component", component.Repo).Err(err).Msg("failed to save vuln(s)")
		}
	}
	log.Logger.Info().Str("dir", conf.CacheDir).Msg("success to save all vuln(s)")
}
