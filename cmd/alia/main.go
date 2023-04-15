package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aliadotsh/alia/internal/datastore"
	"github.com/aliadotsh/alia/internal/graph"
	"github.com/aliadotsh/alia/internal/router"
	"github.com/aliadotsh/alia/internal/server"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/config"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("couldn’t find .env file. reading from env vars.")
	}
}

func SubstrateConfig() config.Config {
	return config.Config{
		RPCURL:           viper.GetString("RPC_URL"),
		DialTimeout:      10 * time.Second,
		SubscribeTimeout: 10 * time.Second,
	}
}

func GetChainInfo(api *gsrpc.SubstrateAPI) {
	// Print out the node name, node version and chain name
	chain, err := api.RPC.System.Chain()
	if err != nil {
		log.Fatalln(err)
	}

	nodeName, err := api.RPC.System.Name()
	if err != nil {
		log.Fatalln(err)
	}

	nodeVersion, err := api.RPC.System.Version()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Connected to chain %v using %v v%v\n", chain, nodeName, nodeVersion)
}

func main() {
	port := viper.GetString("PORT")
	host := viper.GetString("HOST")
	redisURL := viper.GetString("REDIS_URL")
	if port == "" || host == "" || redisURL == "" {
		log.Fatalln("unable to find critical .env vars")
	}

	client, err := datastore.NewRedisClient(redisURL)
	if !errors.Is(err, nil) {
		log.Fatalln(err)
	}
	defer client.Close()

	r := graph.NewResolver(client)
	r.SubscribeRedis()
	srv := server.NewGQLServer(r)

	api, err := gsrpc.NewSubstrateAPI(SubstrateConfig().RPCURL)
	if err != nil {
		log.Fatalln(err)
	}

	GetChainInfo(api)

	/*
	 * Listen to balance changes for a specific account.
	 */
	/*
		meta, err := api.RPC.State.GetMetadataLatest()
		if err != nil {
			log.Fatalln(err)
		}

		publicKeyHex := "0x648c1e47e798840e2101c779ab16db8d46fa4a1cbf5d18fe2fa480e00266c15a"
		publicKey := codec.MustHexDecodeString(publicKeyHex)
		acc := signature.KeyringPair{
			URI:       "//5ELYGahhCk4KwHYTbkYAmzo1xYzmDqnpVCqaXiv5BzfYzcCi",
			Address:   "5ELYGahhCk4KwHYTbkYAmzo1xYzmDqnpVCqaXiv5BzfYzcCi",
			PublicKey: publicKey,
		}
		key, err := types.CreateStorageKey(meta, "System", "Account", acc.PublicKey)
		if err != nil {
			log.Fatalln(err)
		}

		var accInfo types.AccountInfo
		ok, err := api.RPC.State.GetStorageLatest(key, &accInfo)
		if err != nil || !ok {
			log.Fatalln(err)
		}

		previous := accInfo.Data.Free
		fmt.Printf("Account %v has nonce %v \n", acc.Address, accInfo.Nonce)
		fmt.Printf("Balance: %v \n", previous)

		// Subscribe to balance changes
		sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
		if err != nil {
			log.Fatalln(err)
		}
		defer sub.Unsubscribe()

		// outer loop to keep the program running
		for {
			// inner loop to wait for the next balance change
			set := <-sub.Chan()
			for _, chng := range set.Changes {
				if !chng.HasStorageData {
					continue
				}

				var accInfo types.AccountInfo
				err := codec.Decode(chng.StorageData, &accInfo)
				if err != nil {
					log.Fatalln(err)
				}

				current := accInfo.Data.Free
				var change = types.U128{Int: big.NewInt(0).Sub(current.Int, previous.Int)}

				// Only display positive value changes (Since we are pulling `previous` above already,
				// the initial balance change will also be zero)
				if change.Cmp(big.NewInt(0)) != 0 {
					fmt.Println("New balance: ", current)
					fmt.Println("Balance change: ", change)
					fmt.Println("Previous balance: ", previous)

					fmt.Println("==========================================================================")

					previous = current
					return
				}
			}
		}
	*/

	/*
	 * subscribes to new block head events, and prints relevant block information
	 * along with the existence status of a specific module metadata ("ContractsCancelled")
	 * for the first 10 blocks received before unsubscribing from the new block head events.
	 */
	// sub, err := api.RPC.Chain.SubscribeNewHeads()
	// if err != nil {
	// 	panic(err)
	// }
	// defer sub.Unsubscribe()

	// count := 0

	// for {
	// 	head := <-sub.Chan()
	// 	fmt.Printf("Chain is at block: #%v\n", head.Number)
	// 	fmt.Printf("Hash: %#x\n", head.ExtrinsicsRoot)
	// 	fmt.Printf("Chain metadata: %#x\n", head.ParentHash)
	// 	m, err := api.RPC.State.GetMetadata(head.ParentHash)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	spew.Dump(m.AsMetadataV14.ExistsModuleMetadata("ContractsCancelled"))

	// 	count++

	// 	if count == 10 {
	// 		sub.Unsubscribe()
	// 		break
	// 	}
	// }

	/*
	 * subscribes to the system events storage, and listens for new storage changes.
	 * When new changes are detected, it decodes and prints various event
	 * records such as transfers, balances, staking, and session information.
	 */
	// // Subscribe to system events via storage
	// meta, err := api.RPC.State.GetMetadataLatest()
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // Subscribe to system events via storage
	// key, err := types.CreateStorageKey(meta, "System", "Events", nil)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// sub, err := api.RPC.State.SubscribeStorageRaw([]types.StorageKey{key})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// defer sub.Unsubscribe()

	// // outer loop for subscription notifications.
	// for {
	// 	set := <-sub.Chan()
	// 	// inner loop for each storage change set.
	// 	for _, chng := range set.Changes {
	// 		if !codec.Eq(chng.StorageKey, key) || !chng.HasStorageData {
	// 			// skip, we are only interested in events with content
	// 			continue
	// 		}

	// 		// Decode the event records
	// 		events := types.EventRecords{}
	// 		err = types.EventRecordsRaw(chng.StorageData).DecodeEventRecords(meta, &events)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}

	// 		// Print out the events one by one
	// 		// Show what we are busy with
	// 		for _, e := range events.Balances_Endowed {
	// 			fmt.Printf("\tBalances:Endowed:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
	// 		}
	// 		for _, e := range events.Balances_DustLost {
	// 			fmt.Printf("\tBalances:DustLost:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#x, %v\n", e.Who, e.Balance)
	// 		}
	// 		for _, e := range events.Balances_Transfer {
	// 			fmt.Printf("\tBalances:Transfer:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v, %v, %v\n", e.From, e.To, e.Value)
	// 		}
	// 		for _, e := range events.Balances_BalanceSet {
	// 			fmt.Printf("\tBalances:BalanceSet:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v, %v, %v\n", e.Who, e.Free, e.Reserved)
	// 		}
	// 		for _, e := range events.Balances_Deposit {
	// 			fmt.Printf("\tBalances:Deposit:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v, %v\n", e.Who, e.Balance)
	// 		}
	// 		for _, e := range events.Grandpa_NewAuthorities {
	// 			fmt.Printf("\tGrandpa:NewAuthorities:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.NewAuthorities)
	// 		}
	// 		for _, e := range events.Grandpa_Paused {
	// 			fmt.Printf("\tGrandpa:Paused:: (phase=%#v)\n", e.Phase)
	// 		}
	// 		for _, e := range events.Grandpa_Resumed {
	// 			fmt.Printf("\tGrandpa:Resumed:: (phase=%#v)\n", e.Phase)
	// 		}
	// 		for _, e := range events.ImOnline_HeartbeatReceived {
	// 			fmt.Printf("\tImOnline:HeartbeatReceived:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#x\n", e.AuthorityID)
	// 		}
	// 		for _, e := range events.ImOnline_AllGood {
	// 			fmt.Printf("\tImOnline:AllGood:: (phase=%#v)\n", e.Phase)
	// 		}
	// 		for _, e := range events.ImOnline_SomeOffline {
	// 			fmt.Printf("\tImOnline:SomeOffline:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.IdentificationTuples)
	// 		}
	// 		for _, e := range events.Indices_IndexAssigned {
	// 			fmt.Printf("\tIndices:IndexAssigned:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#x%v\n", e.AccountID, e.AccountIndex)
	// 		}
	// 		for _, e := range events.Indices_IndexFreed {
	// 			fmt.Printf("\tIndices:IndexFreed:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.AccountIndex)
	// 		}
	// 		for _, e := range events.Offences_Offence {
	// 			fmt.Printf("\tOffences:Offence:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v%v\n", e.Kind, e.OpaqueTimeSlot)
	// 		}
	// 		for _, e := range events.Session_NewSession {
	// 			fmt.Printf("\tSession:NewSession:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.SessionIndex)
	// 		}
	// 		for _, e := range events.Staking_Rewarded {
	// 			fmt.Printf("\tStaking:Reward:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.Amount)
	// 		}
	// 		for _, e := range events.Staking_Slashed {
	// 			fmt.Printf("\tStaking:Slash:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#x%v\n", e.AccountID, e.Balance)
	// 		}
	// 		for _, e := range events.Staking_OldSlashingReportDiscarded {
	// 			fmt.Printf("\tStaking:OldSlashingReportDiscarded:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.SessionIndex)
	// 		}
	// 		for _, e := range events.System_ExtrinsicSuccess {
	// 			fmt.Printf("\tSystem:ExtrinsicSuccess:: (phase=%#v)\n", e.Phase)
	// 		}
	// 		for _, e := range events.System_ExtrinsicFailed {
	// 			fmt.Printf("\tSystem:ExtrinsicFailed:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%v\n", e.DispatchError)
	// 		}
	// 		for _, e := range events.System_CodeUpdated {
	// 			fmt.Printf("\tSystem:CodeUpdated:: (phase=%#v)\n", e.Phase)
	// 		}
	// 		for _, e := range events.System_NewAccount {
	// 			fmt.Printf("\tSystem:NewAccount:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#x\n", e.Who)
	// 		}
	// 		for _, e := range events.System_KilledAccount {
	// 			fmt.Printf("\tSystem:KilledAccount:: (phase=%#v)\n", e.Phase)
	// 			fmt.Printf("\t\t%#X\n", e.Who)
	// 		}
	// 	}
	// }

	e := router.NewRouter(echo.New(), srv)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%v:%v", host, port)))
}
