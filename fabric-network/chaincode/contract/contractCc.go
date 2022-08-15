package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

type ContractCc struct {
}


func (p *ContractCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 7 {
		return shim.Error("参数为订单编号，个数必须为7")
	}
	contractId := args[0]           //合同id
	contractName := args[1]         //合同名称
	buyer:= args[2]         		//买家
	seller := args[3]   			//卖方
	productName :=args[4]			//订单产品名称
	productAmount :=args[5]			//订单产品数量
	contract_hash :=args[6]			//合同哈希
	txid:=stub.GetTxID()
	timestamp := ts
	value := fmt.Sprintf(`{"contractId":"%s","contractName":"%s","buyer":"%s","seller":"%s","productName":"%s","productAmount":"%s","contract_hash":"%s","timestamp":"%s","txid":"%s"}`,contractId,contractName,buyer,seller,productName,productAmount,contract_hash,timestamp,txid)
	err = stub.PutState(contractId, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *ContractCc) setValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	if len(args) != 7 {
		return shim.Error("参数为订单编号，个数必须为7")
	}
	contractId := args[0]           //合同id
	contractName := args[1]         //合同名称
	buyer:= args[2]         		//买家
	seller := args[3]   			//卖方
	productName :=args[4]			//订单产品名称
	productAmount :=args[5]			//订单产品数量
	contract_hash :=args[6]			//合同哈希
	txid:=stub.GetTxID()
	timestamp := ts
	rsp, err := stub.GetState(contractId)
	if string(rsp) != "" {
		return shim.Error("该合同号已存在")
	}
	value := fmt.Sprintf(`{"contractId":"%s","contractName":"%s","buyer":"%s","seller":"%s","productName":"%s","productAmount":"%s","contract_hash":"%s","timestamp":"%s","txid":"%s"}`,contractId,contractName,buyer,seller,productName,productAmount,contract_hash,timestamp,txid)
	err = stub.PutState(contractId, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *ContractCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为订单编号，个数必须为1")
	}
	contractId := args[0]
	contractId_byte, err := stub.GetState(contractId)
	if err != nil {
		return shim.Error("没有该合同编号-期数")
	}
	return shim.Success(contractId_byte)
}

func (p *ContractCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为订单编号，个数必须为1")
	}
	contractId := args[0]
	rsp,err:=stub.GetState(contractId)
	if string(rsp) == "" || err!=nil {
		return shim.Error("该合同号不存在")
	}
	err = stub.DelState(contractId)
	if err != nil {
		return shim.Error("没有该合同编号")
	}
	return shim.Success([]byte("删除交易成功"))
}

func (p *ContractCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	if len(args) != 7 {
		return shim.Error("参数为订单编号，个数必须为7")
	}
	contractId := args[0]           //合同id
	contractName := args[1]         //合同名称
	buyer:= args[2]         		//买家
	seller := args[3]   			//卖方
	productName :=args[4]			//订单产品名称
	productAmount :=args[5]			//订单产品数量
	contract_hash :=args[6]			//合同哈希
	txid:=stub.GetTxID()
	timestamp := ts
	result, err :=stub.GetState(contractId)
	if err!=nil  {
		return shim.Error("查找需要更新的记录出错")
	}
	if  result == nil {
		return shim.Error("未找到需要更新的记录")
	}
	value := fmt.Sprintf(`{"contractId":"%s","contractName":"%s","buyer":"%s","seller":"%s","productName":"%s","productAmount":"%s","contract_hash":"%s","timestamp":"%s","txid":"%s"}`,contractId,contractName,buyer,seller,productName,productAmount,contract_hash,timestamp,txid)
	err = stub.PutState(contractId, []byte(value))
	if err != nil {
		return shim.Error("更新合同信息失败")
	}
	return shim.Success([]byte("更新合同信息成功"))
}

func (p *ContractCc) queryByProductName(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品名称，个数必须为1")
	}
	productName := fmt.Sprintf("{\"selector\":{\"productName\":\"%s\"}}", args[0])
	resultsIterator, err :=stub.GetQueryResult(productName)
	if err != nil {
		return shim.Error("查询有误")
	}
	defer  resultsIterator.Close() //记得结束富查询调用
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("遍历有误")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(";;")
		}
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *ContractCc) queryContractHistoryByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	contractId := args[0]
	// 获取历史变更数据，返回一个迭代数组,
	//the historic value and associated transaction id and timestamp are returned.
	iterator, err := stub.GetHistoryForKey(contractId)
	if err != nil {
		return shim.Error("根据指定的部门-用户名组合查询对应的历史变更数据失败")
	}
	defer iterator.Close()
	// 迭代处理
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for iterator.HasNext() {
		hisData, err := iterator.Next()//获取当前i的的迭代器数据
		if err != nil {
			return shim.Error("获取历史变更数据失败")
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(";;")
		}
		buffer.WriteString(string(hisData.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *ContractCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	if f == "set" {
		return p.setValue(stub, args)
	}
	if f == "query" {
		return p.query(stub, args)
	}
	if f == "delete" {
		return p.delete(stub, args)
	}
	if f == "update" {
		return p.update(stub, args)
	}
	if f == "queryByProductName" {
		return p.queryByProductName(stub, args)
	}
	if f == "queryContractHistoryByID" {
		return p.queryContractHistoryByID(stub, args)
	}
	return shim.Error("函数名称错误，只能是set、query、delete、queryByProductName或者queryContractDetailByID")
}

func main() {
	err := shim.Start(new(ContractCc))
	if err != nil {
		fmt.Println("启动fabric")
	}
}
