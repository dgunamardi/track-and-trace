package main

import (
	"fmt"
	"log"

	cfg "earhart.com/config"
	parser "earhart.com/parser"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/events/deliverclient/seek"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"

	eventClient "github.com/hyperledger/fabric-sdk-go/pkg/client/event"
)

func main() {
	cfg.LoadConfig()
	cfg.InitializeSDK()
	cfg.InitializeUserIdentity()

	session := cfg.Sdk.Context(fabsdk.WithIdentity(cfg.User))
	channelProvider := func() (context.Channel, error) {
		return contextImpl.NewChannel(session, cfg.CVars.ChannelId)
	}
	ListenToBlockEvents(channelProvider)
}

func ListenToBlockEvents(channelProvider context.ChannelProvider) {

	evClient, err := eventClient.New(
		channelProvider,
		eventClient.WithBlockEvents(),
		eventClient.WithSeekType(seek.Newest),
		eventClient.WithBlockNum(0),
	)

	if err != nil {
		panic(err)
	}

	eventRegister, blockEvents, err := evClient.RegisterBlockEvent()
	defer evClient.Unregister(eventRegister)

	for events := range blockEvents {
		parsedBlock := parser.Block{}
		parsedBlock.Init(events.Block)

		block := events.Block
		DissectHeader(block.GetHeader())
		DissectBody(block.GetData())
		DissectMeta(block.GetMetadata())
	}

}

func DissectHeader(bHeader *common.BlockHeader) {
	fmt.Println("===== Header =====")
	log.Println("raw: ", bHeader)

	log.Println("Content List: ")

	blockNumber := bHeader.GetNumber()
	fmt.Println(":: Block Number: ", blockNumber)
	fmt.Println("")

}

// means contains what it contains
const (
	placeHolder = " -- V -- "
)

func DissectBody(bData *common.BlockData) {
	fmt.Println("===== Data =====")
	log.Println("raw:", placeHolder)

	log.Println("Content List:")

	fmt.Println("For Each:")

	// PREEMPTIVE CHECKS:
	// - CHECK TXSUCCESS CODE (META), ERR -> SKIP TX
	// - CHECK IF TIMESTAMP (CHANHEADER) is defined, ERR -> REJECT
	// - CHECK IF TXACTIONS ARE DEFINED, ERR -> SKIP TX

	for i, d := range bData.GetData() {
		// == ENVELOPE ==
		envelope := &common.Envelope{}
		err := proto.Unmarshal(d, envelope)
		if err != nil {
			panic(err)
		}
		fmt.Printf(":: Envelope#%v:%v\n", i, placeHolder)

		// === E.PAYLOAD ===
		payload := &common.Payload{}
		err = proto.Unmarshal(envelope.GetPayload(), payload)
		if err != nil {
			panic(err)
		}
		fmt.Println("::: E.Payload:", placeHolder)

		// ==== PAYLOAD HEADER ====
		fmt.Println(":::: P.Header == Proposal Header(?) == Transaction Header {common.Header}:", placeHolder)

		chanHeader := &common.ChannelHeader{}
		err = proto.Unmarshal(payload.GetHeader().GetChannelHeader(), chanHeader)
		if err != nil {
			panic(err)
		}
		fmt.Println("::::: Channel Header:", placeHolder)
		fmt.Println(":::::: Channel ID:", chanHeader.GetChannelId())
		fmt.Println(":::::: Type:", common.HeaderType(chanHeader.GetType()))
		fmt.Println(":::::: Timestamp:", chanHeader.GetTimestamp())
		fmt.Println(":::::: Transaction ID:", chanHeader.GetTxId())
		fmt.Println(":::::: Version:", chanHeader.GetVersion())
		fmt.Println(":::::: Extension:", chanHeader.GetExtension(), ", to string:", string(chanHeader.GetExtension()))

		signHeader := &common.SignatureHeader{}
		err = proto.Unmarshal(payload.GetHeader().GetSignatureHeader(), signHeader)
		if err != nil {
			panic(err)
		}
		fmt.Println("::::: Signature Header:", "signHeader", "<- ~marshaled msp.SerializedIdentity + Nonce")

		// ==== PAYLOAD DATA // TRANSACTION ====
		transaction := &peer.Transaction{}
		err = proto.Unmarshal(payload.GetData(), transaction)
		if err != nil {
			panic(err)
		}
		fmt.Println(":::: P.Data == Transaction:", placeHolder)

		// ===== TX ACTION =====
		fmt.Println("::::: Transaction Actions:", placeHolder)
		fmt.Println("::::: For Each:")
		for i, action := range transaction.Actions {
			fmt.Printf("::::: TxAction#%v: %v\n", i, placeHolder)

			txHeader := &common.Header{}
			err := proto.Unmarshal(action.GetHeader(), txHeader)
			if err != nil {
				panic(err)
			}
			fmt.Println(":::::: TA.Header == Proposal Header == Transaction Header {common.Header}:", "<-- P.Header BUT SignatureHeader is empty")

			// ====== CC ACTION PAYLOAD ======
			ccActionPayload := &peer.ChaincodeActionPayload{}
			err = proto.Unmarshal(action.GetPayload(), ccActionPayload)
			if err != nil {
				panic(err)
			}
			fmt.Println(":::::: TA.Payload == Chaincode Action Payload:", placeHolder)

			// if header type is CHAINCODE

			// ======= CC PROPOSAL PAYLOAD =======
			ccProposalPayload := &peer.ChaincodeProposalPayload{}
			err = proto.Unmarshal(ccActionPayload.GetChaincodeProposalPayload(), ccProposalPayload)
			if err != nil {
				panic(err)
			}
			fmt.Println("::::::: Chaincode Proposal Payload:", placeHolder)

			ccInvocationSpec := &peer.ChaincodeInvocationSpec{}
			err = proto.Unmarshal(ccProposalPayload.GetInput(), ccInvocationSpec)
			if err != nil {
				panic(err)
			}
			fmt.Println(":::::::: Chaincode Invocation Spec:", placeHolder)
			ccSpec := ccInvocationSpec.ChaincodeSpec
			fmt.Println("::::::::: Chaincode Spec:", placeHolder)
			fmt.Println(":::::::::: Type:", peer.ChaincodeSpec_Type(ccSpec.GetType()))
			fmt.Println(":::::::::: Id:", ccSpec.GetChaincodeId())
			fmt.Println("::::::::::: Path:", ccSpec.GetChaincodeId().GetPath())
			fmt.Println("::::::::::: Name:", ccSpec.GetChaincodeId().GetName())
			fmt.Println("::::::::::: Version:", ccSpec.GetChaincodeId().GetPath())
			fmt.Println(":::::::::: Input:", ccSpec.GetInput(), "<-- Conversion from string input to current byte structure by UnmarshalJSON in transaction.go")
			fmt.Println(":::::::::: Timeout:", ccSpec.GetTimeout())

			fmt.Println(":::::::: Transient Map:", ccProposalPayload.GetTransientMap())

			// ======= CC ENDORSED ACTION =======
			ccEndorsedAction := ccActionPayload.GetAction()
			fmt.Println("::::::: Chaincode Action:", placeHolder)

			proposalResponsePayload := &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(ccEndorsedAction.GetProposalResponsePayload(), proposalResponsePayload)
			if err != nil {
				panic(err)
			}
			fmt.Println(":::::::: Proposal Response Payload:", placeHolder)

			fmt.Println("::::::::: Proposal Hash:", proposalResponsePayload.GetProposalHash())

			// ======== CC ACTION ========
			// unmarshal extension by type specified in header
			ccAction := &peer.ChaincodeAction{}
			err = proto.Unmarshal(proposalResponsePayload.GetExtension(), ccAction)
			if err != nil {
				panic(err)
			}
			fmt.Println("::::::::: Extension:", placeHolder, "<-- because type is chaincode, unmarshal to chaincodeAction message")

			// ========= TXRWSET ==========
			ccResult := &rwset.TxReadWriteSet{}
			err = proto.Unmarshal(ccAction.GetResults(), ccResult)
			if err != nil {
				panic(err)
			}
			fmt.Println("::::::::: Results == TxRWSet:", placeHolder)
			fmt.Println(":::::::::: DataModel:", ccResult.GetDataModel().String())

			// ========== NSRWSET ===========
			fmt.Println(":::::::::: NSRWSets:", placeHolder)
			fmt.Println(":::::::::: For each:")
			for i, nsrwset := range ccResult.GetNsRwset() {
				fmt.Printf(":::::::::: NSRWSet#%v:%v\n", i, placeHolder)
				fmt.Println("::::::::::: Namespace:", nsrwset.GetNamespace())

				kvrwset := &kvrwset.KVRWSet{}
				err := proto.Unmarshal(nsrwset.GetRwset(), kvrwset)
				if err != nil {
					panic(err)
				}
				fmt.Println("::::::::::: RWSets == KVRWSets:", placeHolder)
				fmt.Println("::::::::::: For Each: ")

				if setCount := len(kvrwset.GetReads()); setCount == len(kvrwset.GetWrites()) {
					for i := 0; i < setCount; i++ {
						fmt.Printf(":::::::::::: RWSet#%v\n", i)
						fmt.Println("::::::::::::: Read:", placeHolder)
						fmt.Println(":::::::::::::: Key:", kvrwset.GetReads()[i].GetKey())
						fmt.Println(":::::::::::::: Version:", kvrwset.GetReads()[i].GetVersion())
						fmt.Println("::::::::::::: Write:", placeHolder)
						fmt.Println(":::::::::::::: Key:", kvrwset.GetWrites()[i].GetKey())
						fmt.Println(":::::::::::::: Value:", string(kvrwset.GetWrites()[i].GetValue()), "<-- to string")
						fmt.Println(":::::::::::::: IsDelete:", kvrwset.GetWrites()[i].GetIsDelete())
					}
				}

				fmt.Println(":::::::::::: RangeQuiresInfo:", kvrwset.GetRangeQueriesInfo(), "<-- is empty?")
				fmt.Println(":::::::::::: MetadataWrites:", kvrwset.GetMetadataWrites(), "<-- is empty?")

				fmt.Println("::::::::::: CollectionHashedRWSet:", nsrwset.GetCollectionHashedRwset(), "<-- is empty?")
			}

			ccEvent := &peer.ChaincodeEvent{}
			err = proto.Unmarshal(ccAction.GetEvents(), ccEvent)
			if err != nil {
				panic(err)
			}
			fmt.Println("::::::::: Events:", ccEvent, "<-- is empty?")

			fmt.Println("::::::::: Response:", placeHolder)
			fmt.Println(":::::::::: Status:", ccAction.GetResponse().GetStatus(), "<-- should follow the HTTPS status codes")
			fmt.Println(":::::::::: Message:", ccAction.GetResponse().GetMessage())
			fmt.Println(":::::::::: Payload:", string(ccAction.GetResponse().GetPayload()), "<-- to string")
			fmt.Println("::::::::: Chaincode ID:", ccAction.GetChaincodeId(), "<-- same type with one from invocationSpec, but more complete")

			fmt.Println(":::::::: Endorsements:", "ccEndorsedAction.GetEndorsements()", "<-- Array of Endorser identity + signature / approval")

		}

		// === E.SIGNATURE ===
		fmt.Println("::: E.Signature:", envelope.Signature)

	}

	fmt.Println("")

}

func DissectMeta(bMeta *common.BlockMetadata) {
	fmt.Println("===== Meta =====")
	log.Println("raw: ", bMeta)

	log.Println("Content List: ")

	txSuccess := bMeta.GetMetadata()[common.BlockMetadataIndex_TRANSACTIONS_FILTER]
	fmt.Println("For Every Transaction: ")

	for t := range txSuccess {
		successCode := peer.TxValidationCode(txSuccess[t])
		fmt.Printf(":: tx#%v valid code: %v\n", t, successCode)
	}

	fmt.Println("")
}
