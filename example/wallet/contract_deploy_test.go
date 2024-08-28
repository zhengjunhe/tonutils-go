package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

type JettonDeployPara struct {
	Name        string
	Description string
	Symbol      string
	Image       string
}

func Test_DeployJetton(t *testing.T) {
	client := liteclient.NewConnectionPool()
	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	words := strings.Split("kick shoot ghost lounge toward grass custom cabin yard walk powder silly boil maid post fuel tennis silk draft coin minute winner forward fruit", " ")
	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("WalletAddress address:", w.WalletAddress())
	log.Println("Wallet old DEPRECATED address:", w.Address())

	msgBody := cell.BeginCell().EndCell()
	maxSupply := int64(2100e4)
	jettonPara := &JettonDeployPara{
		Name:        "jetton of BTC ",
		Description: "This is a jetton version of BTC",
		Symbol:      "WBTC",
		Image:       "https://idhub-cdn.litentry.io/figma/image 1860.png",
	}

	fmt.Println("Deploying jetton coin contract to testnet...")
	addr, _, _, err := w.DeployContractWaitTransaction(context.Background(), tlb.MustFromTON("0.02"),
		msgBody, getJettonCollectionCode(), getContractData(w.WalletAddress(), maxSupply, jettonPara))
	if err != nil {
		panic(err)
	}

	fmt.Println("The new Deployed Jetton contract addr:", addr.String())
}

//2024/08/26 09:59:50 WalletAddress address: UQDSLbwk6S37jKAqmHrNg6OHo6xBzEgm-FLALjNdQui92Ivo
//2024/08/26 09:59:50 Wallet old DEPRECATED address: EQDSLbwk6S37jKAqmHrNg6OHo6xBzEgm-FLALjNdQui92NYt
//Deploying NFT collection contract to mainnet...
//The new Deployed Jetton contract addr: EQBk190EVRg7u_Sk1py3UiefF4EdIEMMBXQBMj4XooVuMVW0
//The new Deployed Jetton contract addr: EQBGdprTu8Jelkhe-3YmhJWyNSJ6MRndrfpSXS0op5RQpk4N
//The new Deployed Jetton contract addr: EQBpw9TQO3UX8VTRAA8GiB7zL0rsvrywL8T1si4AcTEyP9Sy

func getWallet(api ton.APIClientWrapped) *wallet.Wallet {
	words := strings.Split("cement secret mad fatal tip credit thank year toddler arrange good version melt truth embark debris execute answer please narrow fiber school achieve client", " ")
	w, err := wallet.FromSeed(api, words, wallet.V3)
	if err != nil {
		panic(err)
	}
	return w
}

// 使用w.Send来发送交易，没有交易hash返回
func Test_Mint100Jetton(t *testing.T) {
	client := liteclient.NewConnectionPool()
	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	words := strings.Split("drink slab govern current elephant member remain human large wife flavor black grow blue skirt picture auto exact dry reject ancient genuine pair lesson", " ")

	//words := strings.Split("kick shoot ghost lounge toward grass custom cabin yard walk powder silly boil maid post fuel tennis silk draft coin minute winner forward fruit", " ")
	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("WalletAddress address:", w.WalletAddress())
	log.Println("Wallet old DEPRECATED address:", w.Address())

	//kQBk190EVRg7u_Sk1py3UiefF4EdIEMMBXQBMj4XooVuMe4
	//jettonMasterContractAddr := address.MustParseAddr("EQBGdprTu8Jelkhe-3YmhJWyNSJ6MRndrfpSXS0op5RQpk4N")
	jettonMasterContractAddr := address.MustParseAddr("EQBpw9TQO3UX8VTRAA8GiB7zL0rsvrywL8T1si4AcTEyP9Sy")

	//cli := jetton.NewJettonMasterClient(api, jettonMasterContractAddr)
	ctx := api.Client().StickyContext(context.Background())
	master := jetton.NewJettonMasterClient(api, jettonMasterContractAddr)

	tokenWallet, err := master.GetJettonWallet(ctx, w.WalletAddress())
	//tokenWallet, err := master.GetJettonWallet(ctx, address.MustParseAddr("0QBLtV_tFghDJmuGG8bMd62Do33AOEFfDRhZ-ibZMWadkks_"))
	if err != nil {
		log.Fatal(err)
	}
	tokenBalanceBefore, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	amtTon := tlb.MustFromTON("1.15")

	mintPayload := buildMint100Data()

	mintMsg := wallet.SimpleMessage(jettonMasterContractAddr, amtTon, mintPayload)

	//w2 := getWallet(api)
	err = w.Send(context.Background(), mintMsg, true)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		tokenBalanceAfter, err := tokenWallet.GetBalance(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if tokenBalanceAfter.Int64() == tokenBalanceBefore.Int64() {
			fmt.Printf("\nThe balane has not changed, and do the %dth sleep with duration 1s ", i)
			time.Sleep(time.Second)
			continue
		}

		fmt.Printf("\n^-^ The balane has changed form %d to %d ", tokenBalanceBefore.Int64(), tokenBalanceAfter.Int64())
		assert.Equal(t, tokenBalanceAfter.Int64()-int64(100), tokenBalanceBefore.Int64())
		break
	}
	tokenBalanceAfter, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\njetton balance:", tlb.MustFromNano(tokenBalanceAfter, 9))
}

// 使用w.SendWaitTransaction来发送交易，具有交易hash返回
func Test_Mint100Jetton2(t *testing.T) {
	client := liteclient.NewConnectionPool()
	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	words := strings.Split("drink slab govern current elephant member remain human large wife flavor black grow blue skirt picture auto exact dry reject ancient genuine pair lesson", " ")
	//words := strings.Split("kick shoot ghost lounge toward grass custom cabin yard walk powder silly boil maid post fuel tennis silk draft coin minute winner forward fruit", " ")
	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("WalletAddress address:", w.WalletAddress())
	log.Println("Wallet old DEPRECATED address:", w.Address())

	//kQBk190EVRg7u_Sk1py3UiefF4EdIEMMBXQBMj4XooVuMe4
	jettonMasterContractAddr := address.MustParseAddr("EQBpw9TQO3UX8VTRAA8GiB7zL0rsvrywL8T1si4AcTEyP9Sy")

	cli := jetton.NewJettonMasterClient(api, jettonMasterContractAddr)

	ctx := api.Client().StickyContext(context.Background())

	tokenWallet, err := cli.GetJettonWallet(ctx, w.Address())
	if err != nil {
		t.Fatal(err)
	}
	tokenBalanceBefore, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	amtTon := tlb.MustFromTON("1.15")

	mintPayload := buildMint100Data()

	mintMsg := wallet.SimpleMessage(jettonMasterContractAddr, amtTon, mintPayload)

	tx, block, err := w.SendWaitTransaction(ctx, mintMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println("The mint tx hash = ", hex.EncodeToString(tx.Hash))

	// wait next block to be sure everything updated
	block, err = api.WaitForBlock(block.SeqNo + 2).GetMasterchainInfo(ctx)
	if err != nil {
		t.Fatal("Wait master err:", err.Error())
	}

	for i := 0; i < 100; i++ {
		tokenBalanceAfter, err := tokenWallet.GetBalance(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if tokenBalanceAfter.Int64() == tokenBalanceBefore.Int64() {
			fmt.Printf("\nThe balane has not changed, and do the %dth sleep with duration 1s ", i)
			time.Sleep(time.Second)
			continue
		}

		fmt.Printf("\n^-^ The balane has changed form %d to %d ", tokenBalanceBefore.Int64(), tokenBalanceAfter.Int64())
		assert.Equal(t, tokenBalanceAfter.Int64()-int64(100), tokenBalanceBefore.Int64())
		break
	}
	tokenBalanceAfter, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\njetton balance:", tlb.MustFromNano(tokenBalanceAfter, 9))
}

func Test_MintJetton(t *testing.T) {
	client := liteclient.NewConnectionPool()
	// get config
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		log.Fatalln("get config err: ", err.Error())
		return
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	words := strings.Split("kick shoot ghost lounge toward grass custom cabin yard walk powder silly boil maid post fuel tennis silk draft coin minute winner forward fruit", " ")
	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("WalletAddress address:", w.WalletAddress())
	log.Println("Wallet old DEPRECATED address:", w.Address())

	//kQBk190EVRg7u_Sk1py3UiefF4EdIEMMBXQBMj4XooVuMe4
	jettonMasterContractAddr := address.MustParseAddr("EQBpw9TQO3UX8VTRAA8GiB7zL0rsvrywL8T1si4AcTEyP9Sy")

	//cli := jetton.NewJettonMasterClient(api, jettonMasterContractAddr)
	ctx := api.Client().StickyContext(context.Background())
	master := jetton.NewJettonMasterClient(api, jettonMasterContractAddr)
	tokenReceiverWallet, err := master.GetJettonWallet(ctx, address.MustParseAddr("0QBLtV_tFghDJmuGG8bMd62Do33AOEFfDRhZ-ibZMWadkks_"))
	if err != nil {
		log.Fatal(err)
	}
	tokenBalanceBefore, err := tokenReceiverWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	amtTon := tlb.MustFromTON("1.15")
	amtJetton := tlb.MustFromTON("0.0000008")

	to := address.MustParseAddr("0QBLtV_tFghDJmuGG8bMd62Do33AOEFfDRhZ-ibZMWadkks_")

	mintPayload := buildMintDataOpt2(amtJetton.Nano().Int64(), to)

	mintMsg := wallet.SimpleMessage(jettonMasterContractAddr, amtTon, mintPayload)

	//w2 := getWallet(api)
	tx, _, err := w.SendWaitTransaction(ctx, mintMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println("The mint tx hash = ", hex.EncodeToString(tx.Hash))

	for i := 0; i < 100; i++ {
		tokenBalanceAfter, err := tokenReceiverWallet.GetBalance(ctx)
		if err != nil {
			log.Fatal(err)
		}
		if tokenBalanceAfter.Int64() == tokenBalanceBefore.Int64() {
			fmt.Printf("\n The balane has not changed, and do the %dth sleep with duration 1s ", i)
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("\n^-^ The balane has changed form %d to %d ", tokenBalanceBefore.Int64(), tokenBalanceAfter.Int64())
		assert.Equal(t, tokenBalanceAfter.Int64()-amtJetton.Nano().Int64(), tokenBalanceBefore.Int64())
		break
	}
	tokenBalanceAfter, err := tokenReceiverWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\njetton balance:", tlb.MustFromNano(tokenBalanceAfter, 9))
}

func Test_QueryJettonBalance(t *testing.T) {
	client := liteclient.NewConnectionPool()

	// connect to testnet lite server
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		panic(err)
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	tokenContract := address.MustParseAddr("EQBGdprTu8Jelkhe-3YmhJWyNSJ6MRndrfpSXS0op5RQpk4N")
	master := jetton.NewJettonMasterClient(api, tokenContract)

	ctx := client.StickyContext(context.Background())

	tokenWallet, err := master.GetJettonWallet(ctx, address.MustParseAddr("0QDSLbwk6S37jKAqmHrNg6OHo6xBzEgm-FLALjNdQui92DBi"))
	if err != nil {
		log.Fatal(err)
	}

	tokenBalance, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("jetton balance:", tokenBalance.String())
	log.Println("jetton balance with decimal:", tlb.MustFromNano(tokenBalance, 9))
}

func Test_QueryJettonMasterInfo(t *testing.T) {
	client := liteclient.NewConnectionPool()

	// connect to testnet lite server
	cfg, err := liteclient.GetConfigFromUrl(context.Background(), "https://ton.org/testnet-global.config.json")
	if err != nil {
		panic(err)
	}

	// connect to mainnet lite servers
	err = client.AddConnectionsFromConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalln("connection err: ", err.Error())
		return
	}

	// api client with full proof checks
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()
	api.SetTrustedBlockFromConfig(cfg)

	tokenContract := address.MustParseAddr("EQBpw9TQO3UX8VTRAA8GiB7zL0rsvrywL8T1si4AcTEyP9Sy")
	master := jetton.NewJettonMasterClient(api, tokenContract)

	ctx := client.StickyContext(context.Background())
	data, err := master.GetJettonData(ctx)
	if err != nil {
		log.Fatal(err)
	}

	decimals := 9
	content := data.Content.(*nft.ContentOnchain)
	log.Println("total supply:", data.TotalSupply.Uint64())
	log.Println("mintable:", data.Mintable)
	log.Println("admin addr:", data.AdminAddr)
	log.Println("onchain content:")
	log.Println("	name:", content.Name)
	log.Println("	symbol:", content.GetAttribute("symbol"))
	if content.GetAttribute("decimals") != "" {
		decimals, err = strconv.Atoi(content.GetAttribute("decimals"))
		if err != nil {
			log.Fatal("invalid decimals")
		}
	}
	log.Println("	decimals:", decimals)
	log.Println("	description:", content.Description)
	log.Println()

	tokenWallet, err := master.GetJettonWallet(ctx, address.MustParseAddr("UQDSLbwk6S37jKAqmHrNg6OHo6xBzEgm-FLALjNdQui92Ivo"))
	if err != nil {
		log.Fatal(err)
	}

	tokenBalance, err := tokenWallet.GetBalance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("jetton balance:", tlb.MustFromNano(tokenBalance, decimals))
}

func Test_ParseAddr(t *testing.T) {
	addr := address.MustParseAddr("EQBGhqLAZseEqRXz4ByFPTGV7SVMlI4hrbs-Sps_Xzx01x8G")
	fmt.Println("hex ", hex.EncodeToString(addr.Data()))
}

//$ jest
//console.log
//debug jetton mint message: {
//'$$type': 'Mint',
//amount: 100000000000n,
//receiver: EQBGhqLAZseEqRXz4ByFPTGV7SVMlI4hrbs-Sps_Xzx01x8G
//}
//
//at SampleJetton.send (sources/output/SampleJetton_SampleJetton.ts:1344:21)
//
//console.log
//debug jetton body: x{FC708BD20000000000000000000000000000000000000000000000000000000BA43B74004004686A2C066C784A915F3E01C853D3195ED254C948E21ADBB3E4A9B3F5F3C74D7}
//
//FC708BD20000000000000000000000000000000000000000000000000000000BA43B74004004686A2C066C784A915F3E01C853D3195ED254C948E21ADBB3E4A9B3F5F3C74D7
//FC708BD20000000000000000000000000000000000000000000000000000000BA43B74004004686A2C066C784A915F3E01C853D3195ED254C948E21ADBB3E4A9B3F5F3C74D7

func TestJettonMasterClient_Mint(t *testing.T) {
	tt, err := tlb.ToCell(jetton.MintPayload{
		Amount:   tlb.MustFromTON("100"),
		Receiver: address.MustParseAddr("EQBGhqLAZseEqRXz4ByFPTGV7SVMlI4hrbs-Sps_Xzx01x8G"),
	})
	if err != nil {
		t.Fatal(err)
	}

	data := tt.ToBOC()
	fmt.Println(hex.EncodeToString(data))

	fmt.Println(tt.Dump())

	cellMint := buildMintData(100000000000, address.MustParseAddr("EQBGhqLAZseEqRXz4ByFPTGV7SVMlI4hrbs-Sps_Xzx01x8G"))
	fmt.Println(hex.EncodeToString(cellMint.ToBOC()))
	fmt.Println(cellMint.Dump())

	cellMint2 := buildMintDataOpt2(100000000000, address.MustParseAddr("EQBGhqLAZseEqRXz4ByFPTGV7SVMlI4hrbs-Sps_Xzx01x8G"))
	fmt.Println(cellMint2.Dump())

	mint100Cell := buildMint100Data()
	fmt.Println(mint100Cell.Dump())

	amtJetton := tlb.MustFromTON("0.0000008")
	fmt.Println("amtJetton", amtJetton.Nano(), amtJetton.String())
}

type T struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Abi  string `json:"abi"`
	Init struct {
		Kind string `json:"kind"`
		Args []struct {
			Name string `json:"name"`
			Type struct {
				Kind     string `json:"kind"`
				Type     string `json:"type"`
				Optional bool   `json:"optional"`
				Format   int    `json:"format,omitempty"`
			} `json:"type"`
		} `json:"args"`
		Prefix struct {
			Bits  int `json:"bits"`
			Value int `json:"value"`
		} `json:"prefix"`
		Deployment struct {
			Kind   string `json:"kind"`
			System string `json:"system"`
		} `json:"deployment"`
	} `json:"init"`
	Sources struct {
		SourcesJettonTact   string `json:"sources/jetton.tact"`
		SourcesMessagesTact string `json:"sources/messages.tact"`
		SourcesContractTact string `json:"sources/contract.tact"`
	} `json:"sources"`
	Compiler struct {
		Name       string `json:"name"`
		Version    string `json:"version"`
		Parameters string `json:"parameters"`
	} `json:"compiler"`
}

func getJettonCollectionCode() *cell.Cell {
	codeBase64 := "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"
	codeCellBytes, err := base64.StdEncoding.DecodeString(codeBase64)
	if nil != err {
		return nil
	}

	codeCell, err := cell.FromBOC(codeCellBytes)
	if err != nil {
		panic(err)
	}

	return codeCell
}

func getContractData(owner *address.Address, maxSupply int64, jettonPara *JettonDeployPara) *cell.Cell {
	systemBase64 := "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"
	systemData, err := base64.StdEncoding.DecodeString(systemBase64)
	if err != nil {
		panic(err)
	}

	systemCell, err := cell.FromBOC(systemData)
	if err != nil {
		panic(err)
	}

	contentCell := BuildJettonContent(jettonPara.Name, jettonPara.Description, jettonPara.Symbol, jettonPara.Image)

	data := cell.BeginCell().
		MustStoreRef(systemCell).
		MustStoreUInt(0, 1).
		MustStoreAddr(owner).
		MustStoreRef(contentCell).
		MustStoreInt(maxSupply, 257).
		EndCell()

	return data
}

//Name:        "BTC",
//Description: "This is jetton of BTC",

func BuildJettonContent(name, description, symbol, image string) *cell.Cell {
	on := nft.ContentSemichain{
		ContentOnchain: nft.ContentOnchain{
			Name:        name,
			Description: description,
			Image:       image,
		},
	}
	_ = on.SetAttribute("symbol", symbol)

	c, err := on.ContentCell()
	if err != nil {
		panic(err.Error())
	}

	return c
}

func buildMintData(amount int64, receiver *address.Address) *cell.Cell {
	dataBuilder := cell.BeginCell().
		MustStoreUInt(0xfc708bd2, 32).
		MustStoreInt(amount, 257).
		MustStoreAddr(receiver)

	data := cell.BeginCell().
		MustStoreBuilder(dataBuilder).
		EndCell()
	return data
}

func buildMintDataOpt(amount tlb.Coins, receiver *address.Address) *cell.Cell {
	body, err := tlb.ToCell(jetton.MintPayload{
		Amount:   amount,
		Receiver: receiver,
	})
	if err != nil {
		panic(err)
	}
	return body
}

func buildMintDataOpt2(amount int64, receiver *address.Address) *cell.Cell {
	data := cell.BeginCell().
		MustStoreUInt(0xfc708bd2, 32).
		MustStoreInt(amount, 257).
		MustStoreAddr(receiver).
		EndCell()
	return data
}

func buildMint100Data() *cell.Cell {
	data := cell.BeginCell().
		MustStoreUInt(0, 32).
		MustStoreStringSnake("Mint: 100").
		EndCell()
	return data
}
