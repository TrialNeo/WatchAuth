package initialize

import (
	"Diggpher/global"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func LoadConfigs() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.Unmarshal(global.CONFIG)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return
	}

	//载入rsa_pem
	pkey, err := os.ReadFile("./configs/rsa_pri.pem")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	block, _ := pem.Decode(pkey)
	if block == nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	global.RsaPriPem = privateKey.(*rsa.PrivateKey)

}
