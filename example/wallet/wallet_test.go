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
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
	"strings"
	"testing"
)

func Test_CheckTonBalance(t *testing.T) {

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

	// bound all requests to single ton node
	ctx := client.StickyContext(context.Background())

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	//words := strings.Split("diet diet attack autumn expose honey skate lounge holiday opinion village priority major enroll romance famous motor pact hello rubber express warfare rose whisper", " ")

	words := strings.Split("kick shoot ghost lounge toward grass custom cabin yard walk powder silly boil maid post fuel tennis silk draft coin minute winner forward fruit", " ")

	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("wallet address:", w.WalletAddress())

	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}
	log.Println("master proof checks are completed successfully, now communication is 100% safe!")

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	fmt.Printf("Address:%s, Ton balance:%s", w.WalletAddress().String(), balance.String())

	assert.Equal(t, int64(5e9), balance.Nano().Int64())
}

func Test_TonTransfer(t *testing.T) {

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

	// bound all requests to single ton node
	ctx := client.StickyContext(context.Background())

	// seed words of account, you can generate them with any wallet or using wallet.NewSeed() method
	//words := strings.Split("diet diet attack autumn expose honey skate lounge holiday opinion village priority major enroll romance famous motor pact hello rubber express warfare rose whisper", " ")

	words := strings.Split("kick shoot ghost lounge toward grass custom cabin yard walk powder silly boil maid post fuel tennis silk draft coin minute winner forward fruit", " ")

	w, err := wallet.FromSeed(api, words, wallet.V4R2)
	if err != nil {
		log.Fatalln("FromSeed err:", err.Error())
		return
	}

	log.Println("wallet address:", w.WalletAddress())

	log.Println("fetching and checking proofs since config init block, it may take near a minute...")
	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}
	log.Println("master proof checks are completed successfully, now communication is 100% safe!")

	balance, err := w.GetBalance(ctx, block)
	if err != nil {
		log.Fatalln("GetBalance err:", err.Error())
		return
	}

	fmt.Printf("Address:%s, Ton balance:%s", w.WalletAddress().String(), balance.String())

	assert.Greater(t, balance.Nano().Uint64(), uint64(3000000))

	if balance.Nano().Uint64() >= 3000000 {
		addr := address.MustParseAddr("0QA6rzs4xJ4tsitQvUo8GTMHz0_3edTK3ZYNdkwHotd2FerJ")

		log.Println("sending transaction and waiting for confirmation...")

		// if destination wallet is not initialized (or you don't care)
		// you should set bounce to false to not get money back.
		// If bounce is true, money will be returned in case of not initialized destination wallet or smart-contract error
		bounce := false

		transfer, err := w.BuildTransfer(addr, tlb.MustFromTON("0.003"), bounce, "Hello from tonutils-go!")
		if err != nil {
			log.Fatalln("Transfer err:", err.Error())
			return
		}

		tx, block, err := w.SendWaitTransaction(ctx, transfer)
		if err != nil {
			log.Fatalln("SendWaitTransaction err:", err.Error())
			return
		}

		balance, err = w.GetBalance(ctx, block)
		if err != nil {
			log.Fatalln("GetBalance err:", err.Error())
			return
		}

		log.Printf("transaction confirmed at block %d, hash: %s balance left: %s", block.SeqNo,
			base64.StdEncoding.EncodeToString(tx.Hash), balance.String())

		return
	}
}

func Test_AddressParse(t *testing.T) {
	addr := address.MustParseAddr("0QA6rzs4xJ4tsitQvUo8GTMHz0_3edTK3ZYNdkwHotd2FerJ")
	fmt.Println(addr.String())
}

func Test_txHashEncoding(t *testing.T) {
	hash, err := hex.DecodeString("1aeaebd6f52ecb5865930d3a9539089aaa24a430904e6974fe4bbe4d54c530c8")
	assert.Equal(t, nil, err)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)
	fmt.Println(hashBase64)
	//assert.Equal(t, hashBase64, "YcbkP8icpmcd9kxzcgMvoQEbVQ2+nlQlaIL72LX4kBY=")

	hash1, err := hex.DecodeString("61c6e43fc89ca6671df64c7372032fa1011b550dbe9e54256882fbd8b5f89016")
	assert.Equal(t, nil, err)
	hash1Base64 := base64.StdEncoding.EncodeToString(hash1)
	fmt.Println(hash1Base64)
	assert.Equal(t, "YcbkP8icpmcd9kxzcgMvoQEbVQ2+nlQlaIL72LX4kBY=", hash1Base64)
}
