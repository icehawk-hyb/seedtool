package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/33cn/chain33/wallet/bipwallet"
)

var (
	seed         = flag.String("s", "", "seed")
	addr         = flag.String("a", "", "account address")
	count        = flag.Int("c", 100, "count")
	SaveSeedLong = 12
)

func main() {
	flag.Parse()
	fmt.Println("input:seed:addr:count", *seed, *addr, *count)
	if len(*seed) == 0 || len(*addr) == 0 {
		fmt.Println("input seed or addr err")
		return
	}
	findseed, privkey, addr, index, find := checkSeed(*seed, *addr, *count)
	if find {
		fmt.Println("find:seed:privkey:addr:index", findseed, privkey, addr, index)
	} else {
		fmt.Println("not find privkey by seed")
	}
	return
}

//尝试根据输入的seed以及index[0,count]生成公私钥对，比较生成的账户地址是否等于addr
//找到之后需要返回私钥，地址，以及index索引
func checkSeed(seed string, findAddr string, count int) (string, string, string, int, bool) {
	//首先是seed的合法性交易，
	seedarry := strings.Fields(seed)
	curseedlen := len(seedarry)
	if curseedlen < SaveSeedLong {
		return "", "", "", 0, false
	}

	var newseed string
	for index, seedstr := range seedarry {
		if index != curseedlen-1 {
			newseed += seedstr + " "
		} else {
			newseed += seedstr
		}
	}

	//通过seed生成公私钥对,分三种情况处理，
	//1,seed是绝对按照规则取出多余的空格，使用新的方式生成公私钥
	//2,seed是绝对按照规则取出多余的空格，使用旧的方式生成公私钥
	//3,seed直接使用客户传入的即可，不需要取出多余的空格，只能使用旧的方式生成公私钥
	for i := 0; i < count; i++ {
		privkey, addr, err := getPrivkeyBySeed(newseed, uint32(i), true)
		if err == nil {
			if addr == findAddr {
				return newseed, privkey, addr, i, true
			}
		}
		privkey, addr, err = getPrivkeyBySeed(newseed, uint32(i), false)
		if err == nil {
			if addr == findAddr {
				return newseed, privkey, addr, i, true
			}
		}
		privkey, addr, err = getPrivkeyBySeed(seed, uint32(i), false)
		if err == nil {
			if addr == findAddr {
				return seed, privkey, addr, i, true
			}
		}
	}
	return "", "", "", 0, false
}

func getWalletBySeed(seed string, isNewWallet bool) (*bipwallet.HDWallet, error) {
	if isNewWallet {
		return bipwallet.NewWalletFromMnemonic(bipwallet.TypeBty, seed)
	}
	return bipwallet.NewWalletFromSeed(bipwallet.TypeBty, []byte(seed))
}

//isNewWallet=true 使用新的规则生成指定index索引对应的私钥和地址
func getPrivkeyBySeed(seed string, index uint32, isNewWallet bool) (string, string, error) {
	var Hexsubprivkey string
	var addr string
	var err error

	//通过seed使用新旧规则生成wallet
	wallet, err := getWalletBySeed(seed, isNewWallet)
	if err != nil {
		//seedlog.Error("getPrivkeyBySeed getWalletBySeed", "err", err)
		return "", "", errors.New("types.ErrNewWalletFromSeed")
	}

	//通过索引生成Key pair
	priv, pub, err := wallet.NewKeyPair(index)
	if err != nil {
		//seedlog.Error("getPrivkeyBySeed NewKeyPair", "err", err)
		return "", "", errors.New("types.ErrNewKeyPair")
	}

	Hexsubprivkey = hex.EncodeToString(priv)
	//生成addr
	addr, err = bipwallet.PubToAddress(bipwallet.TypeBty, pub)
	if err != nil {
		//seedlog.Error("getPrivkeyBySeed PubToAddress", "err", err)
		return "", "", errors.New("types.ErrPrivkeyToPub")
	}
	return Hexsubprivkey, addr, nil
}
