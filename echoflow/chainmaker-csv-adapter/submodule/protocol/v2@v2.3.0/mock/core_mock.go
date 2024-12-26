// Code generated by MockGen. DO NOT EDIT.
// Source: core_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	msgbus "chainmaker.org/chainmaker/common/v2/msgbus"
	common "chainmaker.org/chainmaker/pb-go/v2/common"
	consensus "chainmaker.org/chainmaker/pb-go/v2/consensus"
	maxbft "chainmaker.org/chainmaker/pb-go/v2/consensus/maxbft"
	txpool "chainmaker.org/chainmaker/pb-go/v2/txpool"
	protocol "chainmaker.org/chainmaker/protocol/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockBlockCommitter is a mock of BlockCommitter interface.
type MockBlockCommitter struct {
	ctrl     *gomock.Controller
	recorder *MockBlockCommitterMockRecorder
}

// MockBlockCommitterMockRecorder is the mock recorder for MockBlockCommitter.
type MockBlockCommitterMockRecorder struct {
	mock *MockBlockCommitter
}

// NewMockBlockCommitter creates a new mock instance.
func NewMockBlockCommitter(ctrl *gomock.Controller) *MockBlockCommitter {
	mock := &MockBlockCommitter{ctrl: ctrl}
	mock.recorder = &MockBlockCommitterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockCommitter) EXPECT() *MockBlockCommitterMockRecorder {
	return m.recorder
}

// AddBlock mocks base method.
func (m *MockBlockCommitter) AddBlock(blk *common.Block) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBlock", blk)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBlock indicates an expected call of AddBlock.
func (mr *MockBlockCommitterMockRecorder) AddBlock(blk interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBlock", reflect.TypeOf((*MockBlockCommitter)(nil).AddBlock), blk)
}

// MockBlockProposer is a mock of BlockProposer interface.
type MockBlockProposer struct {
	ctrl     *gomock.Controller
	recorder *MockBlockProposerMockRecorder
}

// MockBlockProposerMockRecorder is the mock recorder for MockBlockProposer.
type MockBlockProposerMockRecorder struct {
	mock *MockBlockProposer
}

// NewMockBlockProposer creates a new mock instance.
func NewMockBlockProposer(ctrl *gomock.Controller) *MockBlockProposer {
	mock := &MockBlockProposer{ctrl: ctrl}
	mock.recorder = &MockBlockProposerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockProposer) EXPECT() *MockBlockProposerMockRecorder {
	return m.recorder
}

// OnReceiveMaxBFTProposal mocks base method.
func (m *MockBlockProposer) OnReceiveMaxBFTProposal(proposal *maxbft.BuildProposal) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveMaxBFTProposal", proposal)
}

// OnReceiveMaxBFTProposal indicates an expected call of OnReceiveMaxBFTProposal.
func (mr *MockBlockProposerMockRecorder) OnReceiveMaxBFTProposal(proposal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveMaxBFTProposal", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveMaxBFTProposal), proposal)
}

// OnReceiveProposeStatusChange mocks base method.
func (m *MockBlockProposer) OnReceiveProposeStatusChange(proposeStatus bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveProposeStatusChange", proposeStatus)
}

// OnReceiveProposeStatusChange indicates an expected call of OnReceiveProposeStatusChange.
func (mr *MockBlockProposerMockRecorder) OnReceiveProposeStatusChange(proposeStatus interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveProposeStatusChange", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveProposeStatusChange), proposeStatus)
}

// OnReceiveRwSetVerifyFailTxs mocks base method.
func (m *MockBlockProposer) OnReceiveRwSetVerifyFailTxs(rwSetVerifyFailTxs *consensus.RwSetVerifyFailTxs) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveRwSetVerifyFailTxs", rwSetVerifyFailTxs)
}

// OnReceiveRwSetVerifyFailTxs indicates an expected call of OnReceiveRwSetVerifyFailTxs.
func (mr *MockBlockProposerMockRecorder) OnReceiveRwSetVerifyFailTxs(rwSetVerifyFailTxs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveRwSetVerifyFailTxs", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveRwSetVerifyFailTxs), rwSetVerifyFailTxs)
}

// OnReceiveTxPoolSignal mocks base method.
func (m *MockBlockProposer) OnReceiveTxPoolSignal(proposeSignal *txpool.TxPoolSignal) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnReceiveTxPoolSignal", proposeSignal)
}

// OnReceiveTxPoolSignal indicates an expected call of OnReceiveTxPoolSignal.
func (mr *MockBlockProposerMockRecorder) OnReceiveTxPoolSignal(proposeSignal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnReceiveTxPoolSignal", reflect.TypeOf((*MockBlockProposer)(nil).OnReceiveTxPoolSignal), proposeSignal)
}

// ProposeBlock mocks base method.
func (m *MockBlockProposer) ProposeBlock(proposal *maxbft.BuildProposal) (*consensus.ProposalBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProposeBlock", proposal)
	ret0, _ := ret[0].(*consensus.ProposalBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProposeBlock indicates an expected call of ProposeBlock.
func (mr *MockBlockProposerMockRecorder) ProposeBlock(proposal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProposeBlock", reflect.TypeOf((*MockBlockProposer)(nil).ProposeBlock), proposal)
}

// Start mocks base method.
func (m *MockBlockProposer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockBlockProposerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockBlockProposer)(nil).Start))
}

// Stop mocks base method.
func (m *MockBlockProposer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockBlockProposerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockBlockProposer)(nil).Stop))
}

// MockBlockVerifier is a mock of BlockVerifier interface.
type MockBlockVerifier struct {
	ctrl     *gomock.Controller
	recorder *MockBlockVerifierMockRecorder
}

// MockBlockVerifierMockRecorder is the mock recorder for MockBlockVerifier.
type MockBlockVerifierMockRecorder struct {
	mock *MockBlockVerifier
}

// NewMockBlockVerifier creates a new mock instance.
func NewMockBlockVerifier(ctrl *gomock.Controller) *MockBlockVerifier {
	mock := &MockBlockVerifier{ctrl: ctrl}
	mock.recorder = &MockBlockVerifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockVerifier) EXPECT() *MockBlockVerifierMockRecorder {
	return m.recorder
}

// VerifyBlock mocks base method.
func (m *MockBlockVerifier) VerifyBlock(block *common.Block, mode protocol.VerifyMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyBlock", block, mode)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyBlock indicates an expected call of VerifyBlock.
func (mr *MockBlockVerifierMockRecorder) VerifyBlock(block, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyBlock", reflect.TypeOf((*MockBlockVerifier)(nil).VerifyBlock), block, mode)
}

// VerifyBlockSync mocks base method.
func (m *MockBlockVerifier) VerifyBlockSync(block *common.Block, mode protocol.VerifyMode) (*consensus.VerifyResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyBlockSync", block, mode)
	ret0, _ := ret[0].(*consensus.VerifyResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyBlockSync indicates an expected call of VerifyBlockSync.
func (mr *MockBlockVerifierMockRecorder) VerifyBlockSync(block, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyBlockSync", reflect.TypeOf((*MockBlockVerifier)(nil).VerifyBlockSync), block, mode)
}

// VerifyBlockWithRwSets mocks base method.
func (m *MockBlockVerifier) VerifyBlockWithRwSets(block *common.Block, rwsets []*common.TxRWSet, mode protocol.VerifyMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyBlockWithRwSets", block, rwsets, mode)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyBlockWithRwSets indicates an expected call of VerifyBlockWithRwSets.
func (mr *MockBlockVerifierMockRecorder) VerifyBlockWithRwSets(block, rwsets, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyBlockWithRwSets", reflect.TypeOf((*MockBlockVerifier)(nil).VerifyBlockWithRwSets), block, rwsets, mode)
}

// MockCoreEngine is a mock of CoreEngine interface.
type MockCoreEngine struct {
	ctrl     *gomock.Controller
	recorder *MockCoreEngineMockRecorder
}

// MockCoreEngineMockRecorder is the mock recorder for MockCoreEngine.
type MockCoreEngineMockRecorder struct {
	mock *MockCoreEngine
}

// NewMockCoreEngine creates a new mock instance.
func NewMockCoreEngine(ctrl *gomock.Controller) *MockCoreEngine {
	mock := &MockCoreEngine{ctrl: ctrl}
	mock.recorder = &MockCoreEngineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoreEngine) EXPECT() *MockCoreEngineMockRecorder {
	return m.recorder
}

// GetBlockCommitter mocks base method.
func (m *MockCoreEngine) GetBlockCommitter() protocol.BlockCommitter {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockCommitter")
	ret0, _ := ret[0].(protocol.BlockCommitter)
	return ret0
}

// GetBlockCommitter indicates an expected call of GetBlockCommitter.
func (mr *MockCoreEngineMockRecorder) GetBlockCommitter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockCommitter", reflect.TypeOf((*MockCoreEngine)(nil).GetBlockCommitter))
}

// GetBlockProposer mocks base method.
func (m *MockCoreEngine) GetBlockProposer() protocol.BlockProposer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockProposer")
	ret0, _ := ret[0].(protocol.BlockProposer)
	return ret0
}

// GetBlockProposer indicates an expected call of GetBlockProposer.
func (mr *MockCoreEngineMockRecorder) GetBlockProposer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockProposer", reflect.TypeOf((*MockCoreEngine)(nil).GetBlockProposer))
}

// GetBlockVerifier mocks base method.
func (m *MockCoreEngine) GetBlockVerifier() protocol.BlockVerifier {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockVerifier")
	ret0, _ := ret[0].(protocol.BlockVerifier)
	return ret0
}

// GetBlockVerifier indicates an expected call of GetBlockVerifier.
func (mr *MockCoreEngineMockRecorder) GetBlockVerifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockVerifier", reflect.TypeOf((*MockCoreEngine)(nil).GetBlockVerifier))
}

// GetMaxbftHelper mocks base method.
func (m *MockCoreEngine) GetMaxbftHelper() protocol.MaxbftHelper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaxbftHelper")
	ret0, _ := ret[0].(protocol.MaxbftHelper)
	return ret0
}

// GetMaxbftHelper indicates an expected call of GetMaxbftHelper.
func (mr *MockCoreEngineMockRecorder) GetMaxbftHelper() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaxbftHelper", reflect.TypeOf((*MockCoreEngine)(nil).GetMaxbftHelper))
}

// OnMessage mocks base method.
func (m *MockCoreEngine) OnMessage(arg0 *msgbus.Message) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnMessage", arg0)
}

// OnMessage indicates an expected call of OnMessage.
func (mr *MockCoreEngineMockRecorder) OnMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnMessage", reflect.TypeOf((*MockCoreEngine)(nil).OnMessage), arg0)
}

// OnQuit mocks base method.
func (m *MockCoreEngine) OnQuit() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "OnQuit")
}

// OnQuit indicates an expected call of OnQuit.
func (mr *MockCoreEngineMockRecorder) OnQuit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnQuit", reflect.TypeOf((*MockCoreEngine)(nil).OnQuit))
}

// Start mocks base method.
func (m *MockCoreEngine) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockCoreEngineMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockCoreEngine)(nil).Start))
}

// Stop mocks base method.
func (m *MockCoreEngine) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockCoreEngineMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockCoreEngine)(nil).Stop))
}

// MockStoreHelper is a mock of StoreHelper interface.
type MockStoreHelper struct {
	ctrl     *gomock.Controller
	recorder *MockStoreHelperMockRecorder
}

// MockStoreHelperMockRecorder is the mock recorder for MockStoreHelper.
type MockStoreHelperMockRecorder struct {
	mock *MockStoreHelper
}

// NewMockStoreHelper creates a new mock instance.
func NewMockStoreHelper(ctrl *gomock.Controller) *MockStoreHelper {
	mock := &MockStoreHelper{ctrl: ctrl}
	mock.recorder = &MockStoreHelperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStoreHelper) EXPECT() *MockStoreHelperMockRecorder {
	return m.recorder
}

// BeginDbTransaction mocks base method.
func (m *MockStoreHelper) BeginDbTransaction(arg0 protocol.BlockchainStore, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "BeginDbTransaction", arg0, arg1)
}

// BeginDbTransaction indicates an expected call of BeginDbTransaction.
func (mr *MockStoreHelperMockRecorder) BeginDbTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeginDbTransaction", reflect.TypeOf((*MockStoreHelper)(nil).BeginDbTransaction), arg0, arg1)
}

// GetPoolCapacity mocks base method.
func (m *MockStoreHelper) GetPoolCapacity() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPoolCapacity")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetPoolCapacity indicates an expected call of GetPoolCapacity.
func (mr *MockStoreHelperMockRecorder) GetPoolCapacity() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPoolCapacity", reflect.TypeOf((*MockStoreHelper)(nil).GetPoolCapacity))
}

// RollBack mocks base method.
func (m *MockStoreHelper) RollBack(arg0 *common.Block, arg1 protocol.BlockchainStore) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollBack", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RollBack indicates an expected call of RollBack.
func (mr *MockStoreHelperMockRecorder) RollBack(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollBack", reflect.TypeOf((*MockStoreHelper)(nil).RollBack), arg0, arg1)
}