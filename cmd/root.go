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
	"github.com/miao2sec/cloud-native-security-vuln/internal"
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
collect -g

# 使用默认配置进行漏洞收集
collect -r

## 自定义漏洞数据的缓存路径
collect -r -c path/to/cacheDir

## 使用自定义的配置文件进行漏洞收集
collect -r -f <config>.yaml

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
	rootCmd.Flags().StringVarP(&cacheDir, "cache-dir", "c", internal.CacheDir(), "specify the cache directory")
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
		conf *internal.Config
		err  error
	)

	if cacheDir != "" {
		conf = internal.NewConfig(internal.WithCacheDir(cacheDir))
	} else {
		conf = internal.NewConfig()
	}
	if cfgFile != "" {
		conf, err = internal.LoadConfFile(cfgFile)
		if err != nil {
			log.Logger.Fatal().Str("config file", cfgFile).Err(err).Msg("failed to load config file")
		}
	}
	collect(conf)
	return
}

func genCfgFile() {
	var (
		conf           = internal.NewConfig()
		defaultCfgFile = "config.yaml"
	)

	if err := conf.Generate(defaultCfgFile); err != nil {
		log.Logger.Fatal().Err(err).Str("file", defaultCfgFile).Msg("failed to generate config file")
	}
	log.Logger.Info().Str("file", defaultCfgFile).Msg("success to generate config file")
}

func collect(conf *internal.Config) {
	var (
		client = internal.NewClient(internal.WithToken(conf.Token))
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

	// update vuln of k8s
	k8s, err := internal.NewKubernetes()
	if err != nil {
		log.Logger.Fatal().Str("component", "kubernetes").Err(err).Msg("failed to update vuln(s)")
	}
	if err = k8s.Save(conf.CacheDir); err != nil {
		log.Logger.Fatal().Str("component", "kubernetes").Err(err).Msg("failed to save vuln(s)")
	}
	log.Logger.Info().Str("component", "kubernetes").Int("count", len(k8s.Items)).Msg("have vuln(s)")
	log.Logger.Info().Str("dir", conf.CacheDir).Msg("success to save all vuln(s)")
}
