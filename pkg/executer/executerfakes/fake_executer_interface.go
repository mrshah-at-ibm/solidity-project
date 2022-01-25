// Code generated by counterfeiter. DO NOT EDIT.
package executerfakes

import (
	"sync"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mrshah-at-ibm/kaleido-project/pkg/executer"
)

type FakeExecuterInterface struct {
	BalanceOfStub        func(string) (int, error)
	balanceOfMutex       sync.RWMutex
	balanceOfArgsForCall []struct {
		arg1 string
	}
	balanceOfReturns struct {
		result1 int
		result2 error
	}
	balanceOfReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	BurnTokenStub        func(string) (*types.Receipt, error)
	burnTokenMutex       sync.RWMutex
	burnTokenArgsForCall []struct {
		arg1 string
	}
	burnTokenReturns struct {
		result1 *types.Receipt
		result2 error
	}
	burnTokenReturnsOnCall map[int]struct {
		result1 *types.Receipt
		result2 error
	}
	DeployContractStub        func() error
	deployContractMutex       sync.RWMutex
	deployContractArgsForCall []struct {
	}
	deployContractReturns struct {
		result1 error
	}
	deployContractReturnsOnCall map[int]struct {
		result1 error
	}
	GetBalanceStub        func(string) (int, error)
	getBalanceMutex       sync.RWMutex
	getBalanceArgsForCall []struct {
		arg1 string
	}
	getBalanceReturns struct {
		result1 int
		result2 error
	}
	getBalanceReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	MintTokenStub        func(string) (*types.Receipt, error)
	mintTokenMutex       sync.RWMutex
	mintTokenArgsForCall []struct {
		arg1 string
	}
	mintTokenReturns struct {
		result1 *types.Receipt
		result2 error
	}
	mintTokenReturnsOnCall map[int]struct {
		result1 *types.Receipt
		result2 error
	}
	TransferTokenStub        func(string, string, string) (*types.Receipt, error)
	transferTokenMutex       sync.RWMutex
	transferTokenArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	transferTokenReturns struct {
		result1 *types.Receipt
		result2 error
	}
	transferTokenReturnsOnCall map[int]struct {
		result1 *types.Receipt
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeExecuterInterface) BalanceOf(arg1 string) (int, error) {
	fake.balanceOfMutex.Lock()
	ret, specificReturn := fake.balanceOfReturnsOnCall[len(fake.balanceOfArgsForCall)]
	fake.balanceOfArgsForCall = append(fake.balanceOfArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.BalanceOfStub
	fakeReturns := fake.balanceOfReturns
	fake.recordInvocation("BalanceOf", []interface{}{arg1})
	fake.balanceOfMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecuterInterface) BalanceOfCallCount() int {
	fake.balanceOfMutex.RLock()
	defer fake.balanceOfMutex.RUnlock()
	return len(fake.balanceOfArgsForCall)
}

func (fake *FakeExecuterInterface) BalanceOfCalls(stub func(string) (int, error)) {
	fake.balanceOfMutex.Lock()
	defer fake.balanceOfMutex.Unlock()
	fake.BalanceOfStub = stub
}

func (fake *FakeExecuterInterface) BalanceOfArgsForCall(i int) string {
	fake.balanceOfMutex.RLock()
	defer fake.balanceOfMutex.RUnlock()
	argsForCall := fake.balanceOfArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeExecuterInterface) BalanceOfReturns(result1 int, result2 error) {
	fake.balanceOfMutex.Lock()
	defer fake.balanceOfMutex.Unlock()
	fake.BalanceOfStub = nil
	fake.balanceOfReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) BalanceOfReturnsOnCall(i int, result1 int, result2 error) {
	fake.balanceOfMutex.Lock()
	defer fake.balanceOfMutex.Unlock()
	fake.BalanceOfStub = nil
	if fake.balanceOfReturnsOnCall == nil {
		fake.balanceOfReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.balanceOfReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) BurnToken(arg1 string) (*types.Receipt, error) {
	fake.burnTokenMutex.Lock()
	ret, specificReturn := fake.burnTokenReturnsOnCall[len(fake.burnTokenArgsForCall)]
	fake.burnTokenArgsForCall = append(fake.burnTokenArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.BurnTokenStub
	fakeReturns := fake.burnTokenReturns
	fake.recordInvocation("BurnToken", []interface{}{arg1})
	fake.burnTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecuterInterface) BurnTokenCallCount() int {
	fake.burnTokenMutex.RLock()
	defer fake.burnTokenMutex.RUnlock()
	return len(fake.burnTokenArgsForCall)
}

func (fake *FakeExecuterInterface) BurnTokenCalls(stub func(string) (*types.Receipt, error)) {
	fake.burnTokenMutex.Lock()
	defer fake.burnTokenMutex.Unlock()
	fake.BurnTokenStub = stub
}

func (fake *FakeExecuterInterface) BurnTokenArgsForCall(i int) string {
	fake.burnTokenMutex.RLock()
	defer fake.burnTokenMutex.RUnlock()
	argsForCall := fake.burnTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeExecuterInterface) BurnTokenReturns(result1 *types.Receipt, result2 error) {
	fake.burnTokenMutex.Lock()
	defer fake.burnTokenMutex.Unlock()
	fake.BurnTokenStub = nil
	fake.burnTokenReturns = struct {
		result1 *types.Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) BurnTokenReturnsOnCall(i int, result1 *types.Receipt, result2 error) {
	fake.burnTokenMutex.Lock()
	defer fake.burnTokenMutex.Unlock()
	fake.BurnTokenStub = nil
	if fake.burnTokenReturnsOnCall == nil {
		fake.burnTokenReturnsOnCall = make(map[int]struct {
			result1 *types.Receipt
			result2 error
		})
	}
	fake.burnTokenReturnsOnCall[i] = struct {
		result1 *types.Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) DeployContract() error {
	fake.deployContractMutex.Lock()
	ret, specificReturn := fake.deployContractReturnsOnCall[len(fake.deployContractArgsForCall)]
	fake.deployContractArgsForCall = append(fake.deployContractArgsForCall, struct {
	}{})
	stub := fake.DeployContractStub
	fakeReturns := fake.deployContractReturns
	fake.recordInvocation("DeployContract", []interface{}{})
	fake.deployContractMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeExecuterInterface) DeployContractCallCount() int {
	fake.deployContractMutex.RLock()
	defer fake.deployContractMutex.RUnlock()
	return len(fake.deployContractArgsForCall)
}

func (fake *FakeExecuterInterface) DeployContractCalls(stub func() error) {
	fake.deployContractMutex.Lock()
	defer fake.deployContractMutex.Unlock()
	fake.DeployContractStub = stub
}

func (fake *FakeExecuterInterface) DeployContractReturns(result1 error) {
	fake.deployContractMutex.Lock()
	defer fake.deployContractMutex.Unlock()
	fake.DeployContractStub = nil
	fake.deployContractReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeExecuterInterface) DeployContractReturnsOnCall(i int, result1 error) {
	fake.deployContractMutex.Lock()
	defer fake.deployContractMutex.Unlock()
	fake.DeployContractStub = nil
	if fake.deployContractReturnsOnCall == nil {
		fake.deployContractReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deployContractReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeExecuterInterface) GetBalance(arg1 string) (int, error) {
	fake.getBalanceMutex.Lock()
	ret, specificReturn := fake.getBalanceReturnsOnCall[len(fake.getBalanceArgsForCall)]
	fake.getBalanceArgsForCall = append(fake.getBalanceArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetBalanceStub
	fakeReturns := fake.getBalanceReturns
	fake.recordInvocation("GetBalance", []interface{}{arg1})
	fake.getBalanceMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecuterInterface) GetBalanceCallCount() int {
	fake.getBalanceMutex.RLock()
	defer fake.getBalanceMutex.RUnlock()
	return len(fake.getBalanceArgsForCall)
}

func (fake *FakeExecuterInterface) GetBalanceCalls(stub func(string) (int, error)) {
	fake.getBalanceMutex.Lock()
	defer fake.getBalanceMutex.Unlock()
	fake.GetBalanceStub = stub
}

func (fake *FakeExecuterInterface) GetBalanceArgsForCall(i int) string {
	fake.getBalanceMutex.RLock()
	defer fake.getBalanceMutex.RUnlock()
	argsForCall := fake.getBalanceArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeExecuterInterface) GetBalanceReturns(result1 int, result2 error) {
	fake.getBalanceMutex.Lock()
	defer fake.getBalanceMutex.Unlock()
	fake.GetBalanceStub = nil
	fake.getBalanceReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) GetBalanceReturnsOnCall(i int, result1 int, result2 error) {
	fake.getBalanceMutex.Lock()
	defer fake.getBalanceMutex.Unlock()
	fake.GetBalanceStub = nil
	if fake.getBalanceReturnsOnCall == nil {
		fake.getBalanceReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.getBalanceReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) MintToken(arg1 string) (*types.Receipt, error) {
	fake.mintTokenMutex.Lock()
	ret, specificReturn := fake.mintTokenReturnsOnCall[len(fake.mintTokenArgsForCall)]
	fake.mintTokenArgsForCall = append(fake.mintTokenArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.MintTokenStub
	fakeReturns := fake.mintTokenReturns
	fake.recordInvocation("MintToken", []interface{}{arg1})
	fake.mintTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecuterInterface) MintTokenCallCount() int {
	fake.mintTokenMutex.RLock()
	defer fake.mintTokenMutex.RUnlock()
	return len(fake.mintTokenArgsForCall)
}

func (fake *FakeExecuterInterface) MintTokenCalls(stub func(string) (*types.Receipt, error)) {
	fake.mintTokenMutex.Lock()
	defer fake.mintTokenMutex.Unlock()
	fake.MintTokenStub = stub
}

func (fake *FakeExecuterInterface) MintTokenArgsForCall(i int) string {
	fake.mintTokenMutex.RLock()
	defer fake.mintTokenMutex.RUnlock()
	argsForCall := fake.mintTokenArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeExecuterInterface) MintTokenReturns(result1 *types.Receipt, result2 error) {
	fake.mintTokenMutex.Lock()
	defer fake.mintTokenMutex.Unlock()
	fake.MintTokenStub = nil
	fake.mintTokenReturns = struct {
		result1 *types.Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) MintTokenReturnsOnCall(i int, result1 *types.Receipt, result2 error) {
	fake.mintTokenMutex.Lock()
	defer fake.mintTokenMutex.Unlock()
	fake.MintTokenStub = nil
	if fake.mintTokenReturnsOnCall == nil {
		fake.mintTokenReturnsOnCall = make(map[int]struct {
			result1 *types.Receipt
			result2 error
		})
	}
	fake.mintTokenReturnsOnCall[i] = struct {
		result1 *types.Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) TransferToken(arg1 string, arg2 string, arg3 string) (*types.Receipt, error) {
	fake.transferTokenMutex.Lock()
	ret, specificReturn := fake.transferTokenReturnsOnCall[len(fake.transferTokenArgsForCall)]
	fake.transferTokenArgsForCall = append(fake.transferTokenArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.TransferTokenStub
	fakeReturns := fake.transferTokenReturns
	fake.recordInvocation("TransferToken", []interface{}{arg1, arg2, arg3})
	fake.transferTokenMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeExecuterInterface) TransferTokenCallCount() int {
	fake.transferTokenMutex.RLock()
	defer fake.transferTokenMutex.RUnlock()
	return len(fake.transferTokenArgsForCall)
}

func (fake *FakeExecuterInterface) TransferTokenCalls(stub func(string, string, string) (*types.Receipt, error)) {
	fake.transferTokenMutex.Lock()
	defer fake.transferTokenMutex.Unlock()
	fake.TransferTokenStub = stub
}

func (fake *FakeExecuterInterface) TransferTokenArgsForCall(i int) (string, string, string) {
	fake.transferTokenMutex.RLock()
	defer fake.transferTokenMutex.RUnlock()
	argsForCall := fake.transferTokenArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeExecuterInterface) TransferTokenReturns(result1 *types.Receipt, result2 error) {
	fake.transferTokenMutex.Lock()
	defer fake.transferTokenMutex.Unlock()
	fake.TransferTokenStub = nil
	fake.transferTokenReturns = struct {
		result1 *types.Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) TransferTokenReturnsOnCall(i int, result1 *types.Receipt, result2 error) {
	fake.transferTokenMutex.Lock()
	defer fake.transferTokenMutex.Unlock()
	fake.TransferTokenStub = nil
	if fake.transferTokenReturnsOnCall == nil {
		fake.transferTokenReturnsOnCall = make(map[int]struct {
			result1 *types.Receipt
			result2 error
		})
	}
	fake.transferTokenReturnsOnCall[i] = struct {
		result1 *types.Receipt
		result2 error
	}{result1, result2}
}

func (fake *FakeExecuterInterface) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.balanceOfMutex.RLock()
	defer fake.balanceOfMutex.RUnlock()
	fake.burnTokenMutex.RLock()
	defer fake.burnTokenMutex.RUnlock()
	fake.deployContractMutex.RLock()
	defer fake.deployContractMutex.RUnlock()
	fake.getBalanceMutex.RLock()
	defer fake.getBalanceMutex.RUnlock()
	fake.mintTokenMutex.RLock()
	defer fake.mintTokenMutex.RUnlock()
	fake.transferTokenMutex.RLock()
	defer fake.transferTokenMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeExecuterInterface) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ executer.ExecuterInterface = new(FakeExecuterInterface)
