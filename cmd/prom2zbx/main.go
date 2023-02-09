package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"prom2zbx.com/internal/alerts"
	"prom2zbx.com/internal/prom"
	// "prom2zbx.com/internal/zbxsender"
)

var (
	// Used for flags
	cfgFile  string
	prom2zbx = &cobra.Command{
		Use:   "prom2zbx",
		Short: "Integration prometheus alerts to zabbix and back",
		Long: `Service/tool for intergation alerts from alert manager 
			   notifications to zabbix and back`,
	}
)

func Execute() error {
	return prom2zbx.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	prom2zbx.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./prom2zbx.yaml)")
	viper.BindPFlag("promURL", prom2zbx.PersistentFlags().Lookup("promURL"))
	viper.BindPFlag("mode", prom2zbx.PersistentFlags().Lookup("mode"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("prom2zbx")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {

	promURL := flag.String("promURL", "", "Prometheus URL")
	mode := flag.String("mode", "targets", "Mode: targets or rules")
	prefix := flag.String("prefix", "TEST", "Prefix for avoid duplicates")
	help := flag.String("help", "help", "show help")
	flag.Parse()
	switch *mode {
	case "targets":
		prom.GetTargetsProm2LLD(*promURL, *prefix)
	case "rules":
		prom.GetRules(*promURL)
	case "listen":
		alerts.ListenAlerts()
	case "test":
		test()
	}
	usage := ""
}

func test() {
	testcase := []byte(`ABCDEFABDSDWDO5`)
	testcase2 := []byte(`ABCDEFABDSDWDO1.web.test.domain.com:9090`)
	testcase3 := []byte(`ABCDEFABDSDWDO2.web.test.domain.com:9090`)
	testcase4 := []byte(`ABCDEFABDSDWDP1.web.test.domain.com:9090`)
	testcase5 := []byte(`ABCDEFABDSDWDP2.web.test.domain.com:9090`)

	re := regexp.MustCompile("([A-Z0-9]+)")
	res := string(re.Find(testcase))
	res2 := string(re.Find(testcase2))
	res3 := string(re.Find(testcase3))
	res4 := string(re.Find(testcase4))
	res5 := string(re.Find(testcase5))

	fmt.Printf("test: %v\n", res[len(res)-7:])
	fmt.Printf("test: %v\n", res2[len(res2)-7:])
	fmt.Printf("test: %v\n", res3[len(res3)-7:])
	fmt.Printf("test: %v\n", res4[len(res4)-7:])
	fmt.Printf("test: %v\n", res5[len(res5)-7:])
}
