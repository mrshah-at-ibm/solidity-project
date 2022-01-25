package config

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/kube"
	"gopkg.in/yaml.v2"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

/*
nodeurl:
addresses:
- 0x1
- 0x2
claims:
  - address:
    claimtime:

*/

type Claim struct {
	Address   string `yaml:"address"`
	Claimtime int64  `yaml:"claimtime"`
}

type Contract struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
	ABIJSON string `yaml:"abijson"`
}

type Config struct {
	NodeUrl         string     `yaml:"nodeurl"`
	Addresses       []string   `yaml:"addresses"`
	Claims          []*Claim   `yaml:"claims"`
	ClaimExpiration string     `yaml:"claimexpiration"`
	ReceiptWaitMin  uint       `yaml:"receiptwaitmin"`
	ReceiptWaitMax  uint       `yaml:"receiptwaitmax"`
	RPCTimeout      uint       `yaml:"rpctimeout"`
	NumWorkers      uint       `yaml:"numworkers"`
	Contracts       []Contract `yaml:"contracts"`
}

func ReadConfig() (*Config, error) {
	if os.Getenv("INCLUSTER") == "" {
		confyaml, err := ioutil.ReadFile("./config.yaml")
		if err != nil {
			return nil, err
		}

		c := &Config{}
		err = yaml.Unmarshal([]byte(confyaml), c)
		if err != nil {
			return nil, err
		}

		return c, nil

	} else {
		cm, err := kube.ReadConfigMap(os.Getenv("NAMESPACE"), "config")
		if err != nil {
			return nil, err
		}

		confyaml, ok := cm["config"]
		if !ok {
			err = errors.New("config object not found")
			return nil, err
		}

		c := &Config{}
		err = yaml.Unmarshal([]byte(confyaml), c)
		if err != nil {
			return nil, err
		}

		return c, nil
	}
}

func WriteConfig(conf *Config) error {
	b, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	if os.Getenv("INCLUSTER") == "" {
		err = ioutil.WriteFile("./config.yaml", b, 0600)
		if err != nil {
			return err
		}
		return nil
	} else {
		data := map[string]string{
			"config": string(b),
		}
		err = kube.WriteConfigMap(os.Getenv("NAMESPACE"), "config", data)
		if err != nil {
			return err
		}

		return nil
	}
}

func ClaimAddress() ([]string, error) {
	conf, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	var addresses []string
	if len(conf.Addresses)-len(conf.Claims) >= int(conf.NumWorkers) {
		for i := 0; i < int(conf.NumWorkers); i++ {
			newaddr := findNewAddress(conf)
			newclaim := &Claim{
				Address:   newaddr,
				Claimtime: time.Now().Unix(),
			}
			conf.Claims = append(conf.Claims, newclaim)
			addresses = append(addresses, newaddr)
		}

		err = WriteConfig(conf)
		if err != nil {
			return nil, err
		}

		return addresses, nil
	} else {
		for i := 0; i < int(conf.NumWorkers); i++ {
			newaddr, err := searchExpired(conf)
			if err != nil {
				return nil, err
			}
			addresses = append(addresses, newaddr)
		}

		err = WriteConfig(conf)
		if err != nil {
			return nil, err
		}

		return addresses, nil
	}
}

func RetainOwnership(address string) error {
	conf, err := ReadConfig()
	if err != nil {
		return err
	}

	for _, claim := range conf.Claims {
		if claim.Address == address {
			claim.Claimtime = time.Now().Unix()
			break
		}
	}

	err = WriteConfig(conf)
	if err != nil {
		return err
	}

	return nil
}

func findNewAddress(c *Config) string {
	for _, address := range c.Addresses {
		found := false
		for _, claim := range c.Claims {
			if address == claim.Address {
				found = true
				break
			}
		}
		if !found {
			return address
		}
	}
	return ""
}

func searchExpired(c *Config) (string, error) {
	for _, claim := range c.Claims {
		claimtime := time.Unix(claim.Claimtime, 0)

		duration, err := time.ParseDuration(c.ClaimExpiration)
		if err != nil {
			duration = 30 * time.Second
		}

		if time.Since(claimtime) > duration {
			// Check that address is not removed from list of addresses
			for _, address := range c.Addresses {
				found := false
				if claim.Address == address {
					found = true
				}
				if found {
					claim.Claimtime = time.Now().Unix()
					return address, nil
				}
			}
		}
	}
	return "", errors.New("not enough unclaimed address found")
}

func ReadAllPrivateKeys() (map[string]string, error) {
	if os.Getenv("INCLUSTER") == "" {
		b, err := ioutil.ReadFile("./privatekeys.yaml")
		if err != nil {
			if os.IsNotExist(err) {
				data := map[string]string{}
				return data, nil
			}
			return nil, err
		}

		data := map[string]string{}
		err = yaml.Unmarshal(b, data)
		if err != nil {
			return nil, err
		}

		return data, nil
	} else {
		ns := os.Getenv("NAMESPACE")
		data, err := kube.ReadSecret(ns, "privatekey")
		if err != nil {
			if k8serrors.IsNotFound(err) {
				data = map[string]string{}
			} else {
				return nil, err
			}
		}

		return data, nil
	}
}

func WriteAllPrivateKeys(data map[string]string) error {

	if os.Getenv("INCLUSTER") == "" {
		b, err := yaml.Marshal(data)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile("./privatekeys.yaml", b, 0600)
		if err != nil {
			return err
		}
		return nil
	} else {

		ns := os.Getenv("NAMESPACE")

		err := kube.WriteSecret(ns, "privatekey", data)
		return err
	}
}

func SavePrivateKey(address string, key *ecdsa.PrivateKey) error {
	fmt.Println("Provided private key to save:", key)
	// x509Encoded, err := x509.MarshalECPrivateKey(key)
	// if err != nil {
	// 	return err
	// }
	// pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	b := ecrypto.FromECDSA(key)
	fmt.Println("Saving private key:", string(b))
	data, err := ReadAllPrivateKeys()
	if err != nil {
		return err
	}

	if data == nil {
		data = map[string]string{}
	}
	data[address] = string(b)
	err = WriteAllPrivateKeys(data)
	if err != nil {
		return err
	}
	return nil
}

func ReadPrivateKey(address string) (*ecdsa.PrivateKey, error) {

	data, err := ReadAllPrivateKeys()
	if err != nil {
		return nil, err
	}

	pemEncoded := data[address]

	if pemEncoded == "" {
		return nil, nil
	}

	privateKey, err := ecrypto.ToECDSA([]byte(pemEncoded))
	if err != nil {
		return nil, err
	}
	// block, _ := pem.Decode([]byte(pemEncoded))
	// x509Encoded := block.Bytes
	// privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	fmt.Println("Read private key:", privateKey)
	return privateKey, nil
}

func SaveContractAddress(name string, address, abijson string) error {
	conf, err := ReadConfig()
	if err != nil {
		return err
	}

	found := false
	for k, v := range conf.Contracts {
		if v.Name == name {
			found = true
			conf.Contracts[k].Address = address
			break
		}
	}

	if !found {
		c := Contract{
			Name:    name,
			Address: address,
			ABIJSON: abijson,
		}
		conf.Contracts = append(conf.Contracts, c)
	}

	err = WriteConfig(conf)
	if err != nil {
		return err
	}

	return nil

}

func ReadContractAddress(name string) (string, error) {
	conf, err := ReadConfig()
	if err != nil {
		return "", err
	}

	address := ""
	found := false
	for _, v := range conf.Contracts {
		if v.Name == name {
			found = true
			address = v.Address
			break
		}
	}

	if !found {
		return "", errors.New(fmt.Sprintf("Contract address not found: %s", name))
	}

	return address, nil
}

func ReadContractABI(address string) (string, error) {
	conf, err := ReadConfig()
	if err != nil {
		return "", err
	}

	abi := ""
	found := false
	for _, v := range conf.Contracts {
		if v.Address == address {
			found = true
			abi = v.ABIJSON
			break
		}
	}

	if !found {
		return "", errors.New(fmt.Sprintf("Contract abi not found: %s", address))
	}

	return abi, nil
}
