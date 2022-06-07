package main

import (
	"fmt"
	"geth/ewom"
	"geth/womnft"
	"geth/womtx"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

type Config struct {
	Rpc  string  `yaml:"rpc"`
	Ewom Account `yaml:"ewom"`
	Nft  Account `yaml:"nft"`
}

type Account struct {
	Private string `yaml:"private"`
	Address string `yaml:"address"`
}

const EwomCacheKey = "wom:ewom:addr"
const NftCacheKey = "wom:nft:addr"
const TxCacheKey = "wom:tx:addr"

var Client *ethclient.Client
var Redis *redis.Client
var Cfg *Config

func LoadConfig() {
	file, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}
	Cfg = &Config{}
	err = yaml.Unmarshal(file, Cfg)
	if err != nil {
		panic(err)
	}
}

func Connect() {
	if Client == nil {
		var err error
		Client, err = ethclient.Dial(Cfg.Rpc)
		if err != nil {
			panic(err)
		}
		fmt.Println("we have a connection")
	}
}

func redisClient() {
	if Redis == nil {
		Redis = redis.NewClient(&redis.Options{
			Addr:     "192.168.1.168:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		fmt.Println("redis is connection")
	}
}

func main() {
	LoadConfig()
	Connect()
	redisClient()

	app := &cli.App{
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:   "deploy",
				Usage:  "deploy",
				Action: Deploy,
			},
			{
				Name:   "clear",
				Usage:  "clear",
				Action: Clear,
			},
			{
				Name:  "approval",
				Usage: "approval",
				Subcommands: []*cli.Command{
					{
						Name:   "ewom",
						Usage:  "ewom",
						Action: ApprovalEwom,
					},
					{
						Name:   "nft",
						Usage:  "nft",
						Action: ApprovalNFT,
					},
				},
			},
			{
				Name:   "send",
				Usage:  "send",
				Action: Send,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
}

// Send 发送购买交易
func Send(c *cli.Context) error {
	transAddr := Redis.Get(c.Context, TxCacheKey).Val()
	// 创建身份，需要私钥
	auth, err := GetAuth(Cfg.Ewom.Private)
	if err != nil {
		panic(err)
		return err
	}

	womTX, err := womtx.NewWomTransfer(common.HexToAddress(transAddr), Client)
	if err != nil {
		panic(err)
		return err
	}
	log.Println("contract:", transAddr)
	log.Println("seller:", Cfg.Nft.Address)
	log.Println("token:", 0)
	log.Println("number:", 1)
	log.Println("price:", 10)
	coin, err := womTX.MintTransfer(
		auth, common.HexToAddress(Cfg.Nft.Address),
		big.NewInt(0),
		"https://admin.jiacaikeji1899.com/files/test/8888.json",
		big.NewInt(1),
		big.NewInt(10),
		signer2encode(Cfg.Nft.Address, 0, 10))
	if err != nil {
		panic(err)
		return err
	}
	log.Println("seller:", coin)
	return nil
}

// ApprovalEwom 发送EWOM授权
func ApprovalEwom(c *cli.Context) error {
	ewomAddr := Redis.Get(c.Context, EwomCacheKey).Val()
	transAddr := Redis.Get(c.Context, TxCacheKey).Val()

	ewoms, err := ewom.NewEWOMToken(common.HexToAddress(ewomAddr), Client)
	if err != nil {
		panic(err)
		return err
	}

	//创建身份，需要私钥
	auth, err := GetAuth(Cfg.Ewom.Private)
	if err != nil {
		panic(err)
		return err
	}

	approve, err := ewoms.Approve(auth, common.HexToAddress(transAddr), big.NewInt(100000000))
	if err != nil {
		panic(err)
		return err
	}
	log.Println("Approve:", approve)
	return nil
}

// ApprovalNFT 发送NFT授权
func ApprovalNFT(c *cli.Context) error {
	nftAddr := Redis.Get(c.Context, NftCacheKey).Val()
	transAddr := Redis.Get(c.Context, TxCacheKey).Val()

	nft, err := womnft.NewWomNFT(common.HexToAddress(nftAddr), Client)
	if err != nil {
		panic(err)
		return err
	}
	//创建身份，需要私钥
	auth, err := GetAuth(Cfg.Nft.Private)
	if err != nil {
		panic(err)
		return err
	}

	all, err := nft.SetApprovalForAll(auth, common.HexToAddress(transAddr), true)
	if err != nil {
		panic(err)
		return err
	}
	log.Println("SetApprovalForAll:", all)
	return nil
}

func Clear(c *cli.Context) error {
	return Redis.Del(c.Context, EwomCacheKey, NftCacheKey, TxCacheKey).Err()
}

// Deploy 生成`ewom`合约,并产生ewom
func Deploy(c *cli.Context) error {
	//创建身份，需要私钥
	auth, err := GetAuth(Cfg.Ewom.Private)
	if err != nil {
		return err
	}

	// 部署 `ewom`
	if Redis.Exists(c.Context, EwomCacheKey).Val() <= 0 {
		addr, ts, pb, err := ewom.DeployEWOMToken(auth, Client, common.HexToAddress(Cfg.Ewom.Address))
		if err != nil {
			log.Fatal(err)
			return err
		}
		Redis.Set(c.Context, EwomCacheKey, addr.Hex(), -1)
		fmt.Println("ewom deploy success", "addr=", addr.Hex(), ts.Hash().Hex(), pb)
	}

	// 部署 `nft`
	if Redis.Exists(c.Context, NftCacheKey).Val() <= 0 {
		addr, ts, pb, err := womnft.DeployWomNFT(auth, Client)
		if err != nil {
			log.Fatal(err)
			return err
		}
		Redis.Set(c.Context, NftCacheKey, addr.Hex(), -1)
		fmt.Println("nft deploy success", "addr=", addr.Hex(), ts.Hash().Hex(), pb)
	}

	// 部署 `trans`
	if Redis.Exists(c.Context, TxCacheKey).Val() <= 0 {
		womAddr := Redis.Get(c.Context, EwomCacheKey).Val()
		nftAddr := Redis.Get(c.Context, NftCacheKey).Val()

		addr, ts, pb, err := womtx.DeployWomTransfer(auth, Client, common.HexToAddress(womAddr), common.HexToAddress(nftAddr))
		if err != nil {
			log.Fatal(err)
			return err
		}
		Redis.Set(c.Context, TxCacheKey, addr.Hex(), -1)
		fmt.Println("trans deploy success", "addr=", addr.Hex(), ts.Hash().Hex(), pb)
	}
	return nil
}

// GetAuth 获取授权
func GetAuth(key string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, err
	}

	return bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(4))
}

func signer2encode(addr string, tokenID, price int64) []byte {
	uint256, _ := abi.NewType("uint256", "", nil)
	address, _ := abi.NewType("address", "", nil)

	arguments := abi.Arguments{
		{
			Type: address,
		},
		{
			Type: uint256,
		},
		{
			Type: uint256,
		},
	}

	bytes, _ := arguments.Pack(
		common.HexToAddress(addr),
		big.NewInt(tokenID),
		big.NewInt(price),
	)

	msg := crypto.Keccak256(bytes)
	log.Println("Keccak256:", common.Bytes2Hex(msg))

	key, _ := crypto.HexToECDSA(Cfg.Nft.Private)
	sig, err := crypto.Sign(msg, key)
	if err != nil {
		panic(err)
	}
	if len(sig) != 65 {
		panic("sig error")
	}
	switch sig[64] {
	case 0:
		sig[64] = 27
	case 1:
		sig[64] = 28
	default:
	}

	log.Println("sig:", common.Bytes2Hex(sig))
	return sig
}
