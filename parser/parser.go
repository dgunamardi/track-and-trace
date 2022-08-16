package parser

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// 1 Block -> Multiple Transactions
// 1 Transaction -> Multiple Actions
// 1 Action -> Multiple NamespaceReadWriteSets
// NsRWSet -> single asset key value read

type Block struct {
	BlockProto *common.Block

	BlockData BlockData
}

func (b *Block) Init(block *common.Block) {
	b.BlockProto = block
	b.GetBlockData()
}

func (b *Block) GetBlockData() {
	b.BlockData.Init(b.BlockProto.GetData())
}

type BlockData struct {
	BlockDataProto *common.BlockData

	Envelopes []Envelope
}

func (bd *BlockData) Init(blockData *common.BlockData) {
	bd.BlockDataProto = blockData
	bd.ParseEnvelope()
}

func (bd *BlockData) ParseEnvelope() {
	dataArray := bd.BlockDataProto.GetData()
	for _, data := range dataArray {
		envelope := &common.Envelope{}
		err := proto.Unmarshal(data, envelope)
		if err != nil {
			panic(fmt.Errorf("f to parse envelope: %v", err))
		}
		parsedEnvelope := Envelope{}
		parsedEnvelope.Init(envelope)
		bd.Envelopes = append(bd.Envelopes, parsedEnvelope)
	}
}

type Envelope struct {
	EnvelopeProto *common.Envelope

	Payload   Payload
	Signature []byte

	IsTransaction bool
}

func (e *Envelope) Init(envelope *common.Envelope) {
	e.IsTransaction = true
	e.EnvelopeProto = envelope
	//tbm
	//fmt.Println(e.ENvelopeProto)
	e.Signature = envelope.GetSignature()
	e.ParsePayload()
	//fmt.Println(e.IsTransaction)
}

func (e *Envelope) ParsePayload() {
	payload := &common.Payload{}
	err := proto.Unmarshal(e.EnvelopeProto.GetPayload(), payload)
	if err != nil {
		panic(fmt.Errorf("f to parse payload: %v", err))
	}
	e.IsTransaction = e.Payload.Init(payload)
}

type Payload struct {
	PayloadProto *common.Payload

	Header      Header
	Transaction Transaction
}

func (p *Payload) Init(payload *common.Payload) bool {
	p.PayloadProto = payload
	// tbm
	//fmt.Println(p.PayloadProto)
	p.GetHeader()
	//fmt.Println(p.Header.ChannelHeader.ChannelHeaderProto.Type)
	//
	if p.Header.ChannelHeader.ChannelHeaderProto.Type != 3 {
		return false
	}
	p.ParseTransaction()

	return true
}

func (p *Payload) GetHeader() {
	p.Header.Init(p.PayloadProto.GetHeader())
}

func (p *Payload) ParseTransaction() {
	transaction := &peer.Transaction{}
	err := proto.Unmarshal(p.PayloadProto.GetData(), transaction)
	if err != nil {
		panic(fmt.Errorf("f to parse transaction: %v", err))
	}
	p.Transaction.Init(transaction)
}

type Header struct {
	HeaderProto *common.Header

	ChannelHeader   ChannelHeader
	SignatureHeader []byte
}

func (h *Header) Init(header *common.Header) {
	h.HeaderProto = header
	h.ParseChannelHeader()
	h.SignatureHeader = header.SignatureHeader
}

func (h *Header) ParseChannelHeader() {
	channelHeader := &common.ChannelHeader{}
	err := proto.Unmarshal(h.HeaderProto.GetChannelHeader(), channelHeader)
	if err != nil {
		panic(fmt.Errorf("f to parse channel header: %v", err))
	}
	h.ChannelHeader.Init(channelHeader)
}

type ChannelHeader struct {
	ChannelHeaderProto *common.ChannelHeader
}

func (ch *ChannelHeader) Init(channelHeader *common.ChannelHeader) {
	ch.ChannelHeaderProto = channelHeader
	//tbm
	//fmt.Println(ch.ChannelHeaderProto)
}

type Transaction struct {
	TransactionProto *peer.Transaction

	TransactionActions []TransactionAction
}

func (t *Transaction) Init(transaction *peer.Transaction) {
	t.TransactionProto = transaction
	// tbm
	//fmt.Println(transaction)
	t.GetTransactionActions()
}

func (t *Transaction) GetTransactionActions() {
	for _, transactionAction := range t.TransactionProto.GetActions() {
		txAction := TransactionAction{}
		txAction.Init(transactionAction)
		t.TransactionActions = append(t.TransactionActions, txAction)
	}
}

type TransactionAction struct {
	TransactionActionProto *peer.TransactionAction

	ChaincodeActionHeader  []byte
	ChaincodeActionPayload ChaincodeActionPayload
}

func (ta *TransactionAction) Init(transactionAction *peer.TransactionAction) {
	ta.TransactionActionProto = transactionAction
	// tbm
	//fmt.Printf("action:\n%v\n", ta.TransactionAction)
	ta.ChaincodeActionHeader = transactionAction.GetHeader()
	ta.ParseChaincodeActionPayload()
}

func (ta *TransactionAction) ParseChaincodeActionPayload() {
	chaincodeActionPayload := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(ta.TransactionActionProto.GetPayload(), chaincodeActionPayload)
	if err != nil {
		panic(fmt.Errorf("f to parse ChaincodeActionPayload: %v", err))
	}
	ta.ChaincodeActionPayload.Init(chaincodeActionPayload)
}

type ChaincodeActionPayload struct {
	ChaincodeActionPayloadProto *peer.ChaincodeActionPayload

	ChaincodeProposalPayload ChaincodeProposalPayload
	ChaincodeEndorsedAction  ChaincodeEndorsedAction
}

func (cap *ChaincodeActionPayload) Init(chaincodeActionPayload *peer.ChaincodeActionPayload) {
	cap.ChaincodeActionPayloadProto = chaincodeActionPayload
	// tbm
	//fmt.Println(cap.ChaincodeActionPayloadProto)
	cap.ParseChaincodeProposalPayload()
	cap.GetChaincodeEndorsedAction()
}

func (cap *ChaincodeActionPayload) ParseChaincodeProposalPayload() {
	chaincodeProposalPayload := &peer.ChaincodeProposalPayload{}
	err := proto.Unmarshal(cap.ChaincodeActionPayloadProto.GetChaincodeProposalPayload(), chaincodeProposalPayload)
	if err != nil {
		panic(fmt.Errorf("f to parse proposal payload: %v", err))
	}
	cap.ChaincodeProposalPayload.Init(chaincodeProposalPayload)
}

func (cap *ChaincodeActionPayload) GetChaincodeEndorsedAction() {
	cap.ChaincodeEndorsedAction.Init(cap.ChaincodeActionPayloadProto.GetAction())
}

type ChaincodeProposalPayload struct {
	ChaincodeProposalPayloadProto *peer.ChaincodeProposalPayload

	// Expand Later
}

func (cpp *ChaincodeProposalPayload) Init(chaincodeProposalPayload *peer.ChaincodeProposalPayload) {
	cpp.ChaincodeProposalPayloadProto = chaincodeProposalPayload
	// tbm
	//fmt.Println(cpp.ChaincodeProposalPayloadProto)
}

type ChaincodeEndorsedAction struct {
	ChaincodeEndorsedActionProto *peer.ChaincodeEndorsedAction

	ProposalResponsePayload ProposalResponsePayload
	Endorsements            []Endorsement
}

func (cea *ChaincodeEndorsedAction) Init(endorsedAction *peer.ChaincodeEndorsedAction) {
	cea.ChaincodeEndorsedActionProto = endorsedAction
	// tbm
	//fmt.Println(cea.ChaincodeEndorsedActionProto)
	cea.ParseProposalResponsePayload()
	cea.GetEndorsements()
}

func (cea *ChaincodeEndorsedAction) ParseProposalResponsePayload() {
	proposalResponsePayload := &peer.ProposalResponsePayload{}
	err := proto.Unmarshal(cea.ChaincodeEndorsedActionProto.GetProposalResponsePayload(), proposalResponsePayload)
	if err != nil {
		panic(fmt.Errorf("f to parse proposal response payload: %v", err))
	}
	cea.ProposalResponsePayload.Init(proposalResponsePayload)
}

func (cea *ChaincodeEndorsedAction) GetEndorsements() {
	for _, endorsement := range cea.ChaincodeEndorsedActionProto.GetEndorsements() {
		tEndorsement := Endorsement{}
		tEndorsement.Init(endorsement)
		cea.Endorsements = append(cea.Endorsements, tEndorsement)
	}
}

type ProposalResponsePayload struct {
	ProposalResponsePayloadProto *peer.ProposalResponsePayload

	ProposalHash []byte
	Extension    ChaincodeAction
}

func (prp *ProposalResponsePayload) Init(proposalResponsePayload *peer.ProposalResponsePayload) {
	prp.ProposalResponsePayloadProto = proposalResponsePayload
	// tbm
	//fmt.Println(proposalResponsePayload)
	prp.ProposalHash = proposalResponsePayload.GetProposalHash()
	prp.ParseExtension()
}

func (prp *ProposalResponsePayload) ParseExtension() {
	chaincodeAction := &peer.ChaincodeAction{}
	err := proto.Unmarshal(prp.ProposalResponsePayloadProto.GetExtension(), chaincodeAction)
	if err != nil {
		panic(fmt.Errorf("f to parse extension: %v", chaincodeAction))
	}
	prp.Extension.Init(chaincodeAction)
}

type Endorsement struct {
	EndorsementProto *peer.Endorsement

	// Expand Later
}

func (en *Endorsement) Init(endorsement *peer.Endorsement) {
	en.EndorsementProto = endorsement
	// tbm
	//fmt.Println(en.EndorsementProto)
}

type ChaincodeAction struct {
	ChaincodeActionProto *peer.ChaincodeAction

	Results TxReadWriteSet
	Events  ChaincodeEvent

	// Expanded Later
	Response    *peer.Response
	ChaincodeId *peer.ChaincodeID
}

func (ca *ChaincodeAction) Init(chaincodeAction *peer.ChaincodeAction) {
	ca.ChaincodeActionProto = chaincodeAction
	// tbm
	//fmt.Println(ca.ChaincodeActionProto)
	ca.ParseResults()
	ca.ParseEvents()

	ca.Response = chaincodeAction.GetResponse()
	ca.ChaincodeId = chaincodeAction.GetChaincodeId()
}

func (ca *ChaincodeAction) ParseResults() {
	txRWSet := &rwset.TxReadWriteSet{}
	err := proto.Unmarshal(ca.ChaincodeActionProto.GetResults(), txRWSet)
	if err != nil {
		panic(fmt.Errorf("f to parse txrwset: %v", err))
	}
	ca.Results.Init(txRWSet)
}

func (ca *ChaincodeAction) ParseEvents() {
	chaincodeEvent := &peer.ChaincodeEvent{}
	err := proto.Unmarshal(ca.ChaincodeActionProto.GetEvents(), chaincodeEvent)
	if err != nil {
		panic(fmt.Errorf("f to parse chaincodeEvent: %v", err))
	}
	ca.Events.Init(chaincodeEvent)
}

type TxReadWriteSet struct {
	TxReadWriteSetProto *rwset.TxReadWriteSet

	DataModel       rwset.TxReadWriteSet_DataModel
	NsReadWriteSets []NsReadWriteSet
}

func (txrw *TxReadWriteSet) Init(txRWSet *rwset.TxReadWriteSet) {
	txrw.TxReadWriteSetProto = txRWSet
	// tbm
	//fmt.Println(txrw.TxReadWriteSetProto)
	txrw.DataModel = txRWSet.GetDataModel()
	txrw.GetNsReadWriteSets()
}

func (txrw *TxReadWriteSet) GetNsReadWriteSets() {
	for _, nsRWSet := range txrw.TxReadWriteSetProto.GetNsRwset() {
		tNsRWSet := NsReadWriteSet{}
		tNsRWSet.Init(nsRWSet)
		txrw.NsReadWriteSets = append(txrw.NsReadWriteSets, tNsRWSet)
	}
}

type ChaincodeEvent struct {
	ChaincodeEventProto *peer.ChaincodeEvent

	ChaincodeId  string
	TxId         string
	EventName    string
	EventPayload string
}

func (ce *ChaincodeEvent) Init(chaincodeEvent *peer.ChaincodeEvent) {
	ce.ChaincodeEventProto = chaincodeEvent
	// tbm
	//fmt.Println(ce.ChaincodeEventProto)
	ce.ChaincodeId = chaincodeEvent.GetChaincodeId()
	ce.TxId = chaincodeEvent.GetTxId()
	ce.EventName = chaincodeEvent.GetEventName()
	ce.EventPayload = string(chaincodeEvent.GetPayload())
}

type NsReadWriteSet struct {
	NsReadWriteSetProto *rwset.NsReadWriteSet

	Namespace string
	RWSet     KVRWSet

	// Expand Later
	CollectionHashedRWSet []*rwset.CollectionHashedReadWriteSet
}

func (nsrw *NsReadWriteSet) Init(nsRWSet *rwset.NsReadWriteSet) {
	nsrw.NsReadWriteSetProto = nsRWSet
	// tbm
	//fmt.Println(nsrw.NsReadWriteSetProto)
	nsrw.Namespace = nsRWSet.GetNamespace()
	nsrw.ParseRWSet()
	nsrw.CollectionHashedRWSet = nsRWSet.GetCollectionHashedRwset()
}

func (nsrw *NsReadWriteSet) ParseRWSet() {
	kvRWSet := &kvrwset.KVRWSet{}
	err := proto.Unmarshal(nsrw.NsReadWriteSetProto.GetRwset(), kvRWSet)
	if err != nil {
		panic(fmt.Errorf("f to parse kvrwset: %v", err))
	}
	nsrw.RWSet.Init(kvRWSet)
}

type KVRWSet struct {
	KVRWSetProto *kvrwset.KVRWSet

	KVReads          []KVRead
	RangeQueriesInfo []*kvrwset.RangeQueryInfo
	KVWrites         []KVWrite
	MetadataWrites   []*kvrwset.KVMetadataWrite
}

func (kvrw *KVRWSet) Init(kvRWSet *kvrwset.KVRWSet) {
	kvrw.KVRWSetProto = kvRWSet
	// tbm
	//fmt.Println(kvrw.KVRWSetProto)
	kvrw.GetReads()
	kvrw.RangeQueriesInfo = kvRWSet.GetRangeQueriesInfo()
	kvrw.GetWrites()
	kvrw.MetadataWrites = kvRWSet.GetMetadataWrites()
}

func (kvrw *KVRWSet) GetReads() {
	for _, kvRead := range kvrw.KVRWSetProto.GetReads() {
		tKVRead := KVRead{}
		tKVRead.Init(kvRead)
		kvrw.KVReads = append(kvrw.KVReads, tKVRead)
	}
}

func (kvrw *KVRWSet) GetWrites() {
	for _, kvWrite := range kvrw.KVRWSetProto.GetWrites() {
		tKVWrite := KVWrite{}
		tKVWrite.Init(kvWrite)
		kvrw.KVWrites = append(kvrw.KVWrites, tKVWrite)
	}
}

type KVRead struct {
	KVReadProto *kvrwset.KVRead

	Key     string
	Version *kvrwset.Version
}

func (r *KVRead) Init(kvRead *kvrwset.KVRead) {
	r.KVReadProto = kvRead
	// tbm
	//fmt.Printf("key: %v\n", r.KVReadProto.Key)
	r.Key = kvRead.GetKey()
	r.Version = kvRead.GetVersion()
}

type KVWrite struct {
	KVWriteProto *kvrwset.KVWrite

	Key         string
	IsDelete    bool
	Value       []byte
	ValueString string
}

func (w *KVWrite) Init(kvWrite *kvrwset.KVWrite) {
	w.KVWriteProto = kvWrite
	// tbm
	//fmt.Printf("key: %v, value: %v\n", w.KVWriteProto.Key, string(w.KVWriteProto.Value))
	w.Key = kvWrite.GetKey()
	w.IsDelete = kvWrite.GetIsDelete()
	w.Value = kvWrite.GetValue()
	w.ValueString = string(w.Value)
}
