package main

import (
"bytes"
"fmt"
"github.com/hyperledger/fabric/core/chaincode/shim"
pb "github.com/hyperledger/fabric/protos/peer"
//"strconv"
"time"
)

type GongdanCc struct {
}

func (p *GongdanCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	work_order_id := args[0]   		  //图纸标识号
	task := args[1]   		  //工单任务文件
	duration := args[2]   		  //持续时间
	drawing_id := args[3]   		  //图纸id
	contract_id:= args[4]   		  	  //合同id
	txid:=stub.GetTxID()
	timestamp := ts
	value := fmt.Sprintf(`{"work_order_id":"%s","task":"%s","duration":"%s","drawing_id":"%s","contract_id":"%s","timestamp":"%s","txid":"%s"}`,work_order_id,task,duration,drawing_id,contract_id,timestamp,txid)
	err = stub.PutState(work_order_id, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *GongdanCc) setValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	work_order_id := args[0]   		  //图纸标识号
	task := args[1]   		  //工单任务
	duration := args[2]   		  //持续时间
	drawing_id := args[3]   		  //图纸id
	contract_id:= args[4]   		  	  //合同id
	txid:=stub.GetTxID()
	timestamp := ts
	rsp, err := stub.GetState(drawing_id)
	if string(rsp) != "" {
		return shim.Error("该图纸标识号已存在")
	}

	value := fmt.Sprintf(`{"work_order_id":"%s","task":"%s","duration":"%s","drawing_id":"%s","contract_id":"%s","timestamp":"%s","txid":"%s"}`,work_order_id,task,duration,drawing_id,contract_id,timestamp,txid)

	err = stub.PutState(work_order_id, []byte(value))

	if err != nil {
		return shim.Error("新增图纸失败")
	}
	return shim.Success([]byte("新增图纸成功"))
}

func (p *GongdanCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("参数为图纸标识码，参数个数必须为1")
	}
	work_order_id := args[0]
	work_order_id_byte, err := stub.GetState(work_order_id)
	if err != nil {
		return shim.Error("没有图纸信息")
	}
	return shim.Success(work_order_id_byte)
}

func (p *GongdanCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为图纸标识码，参数个数必须为1")
	}
	work_order_id := args[0]
	rsp,err:=stub.GetState(work_order_id)
	if string(rsp) == "" || err!=nil {
		return shim.Error("该图纸编号不存在")
	}
	err = stub.DelState(work_order_id)
	if err != nil {
		return shim.Error("没有图纸信息")
	}
	return shim.Success([]byte("删除图纸成功"))
}

func (p *GongdanCc) queryByContractId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为合同名称编号，个数必须为1")
	}
	//contractName := args[0]
	contractId := fmt.Sprintf("{\"selector\":{\"contract_id\":\"%s\"}}", args[0])
	resultsIterator, err :=stub.GetQueryResult(contractId)
	if err != nil {
		return shim.Error("查询有误")
	}
	defer  resultsIterator.Close() //记得结束富查询调用
	// buffer is a JSON array containing QueryRecords
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
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *GongdanCc) queryByDrawingId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为图纸编号，个数必须为1")
	}
	drawing_id := fmt.Sprintf("{\"selector\":{\"drawing_id\":\"%s\"}}", args[0])
	resultsIterator, err :=stub.GetQueryResult(drawing_id)
	if err != nil {
		return shim.Error("查询有误")
	}
	defer  resultsIterator.Close() //记得结束富查询调用
	// buffer is a JSON array containing QueryRecords
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
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *GongdanCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp, _ :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	work_order_id := args[0]   		  //图纸标识号
	task := args[1]   		  //工单任务文件
	duration := args[2]   		  //持续时间
	drawing_id := args[3]   		  //图纸id
	contract_id:= args[4]   		  	  //合同id
	txid:=stub.GetTxID()
	timestamp := ts
	result, err :=stub.GetState(drawing_id)
	if err!=nil ||result ==nil {
		return shim.Error("未找到需要更新的图纸")
	}
	value := fmt.Sprintf(`{"work_order_id":"%s","task":"%s","duration":"%s","drawing_id":"%s","contract_id":"%s","timestamp":"%s","txid":"%s"}`,work_order_id,task,duration,drawing_id,contract_id,timestamp,txid)
	err = stub.PutState(drawing_id, []byte(value))
	if err != nil {
		return shim.Error("更新图纸失败")
	}
	return shim.Success([]byte("更新图纸成功"))

}

func (p *GongdanCc) queryGongdanHistoryByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	gongdan_id :=args[0]
	iterator, err := stub.GetHistoryForKey(gongdan_id)
	if err != nil {
		return shim.Error("根据指定的标识号查询对应的历史变更数据失败")
	}
	defer iterator.Close()

	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for iterator.HasNext() {
		hisData, err := iterator.Next()//获取当前i的的迭代器数据
		if err != nil {
			return shim.Error("获取工单的历史变更数据失败")
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

func (p *GongdanCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()
	if f == "set" {
		return p.setValue(stub, args)
	}
	if f == "query" {
		return p.query(stub, args)
	}
	if f == "update" {
		return p.update(stub, args)
	}
	if f == "delete" {
		return p.delete(stub, args)
	}
	if f == "queryByContractId" {
		return p.queryByContractId(stub, args)
	}
	if f == "queryByDrawingId" {
		return p.queryByDrawingId(stub, args)
	}
	if f == "queryGongdanHistoryByID" {
		return p.queryGongdanHistoryByID(stub, args)
	}
	return shim.Error("函数名称错误，只能是set或query或者delete")
}

func main() {
	err := shim.Start(new(GongdanCc))
	if err != nil {
		fmt.Println("启动fabric")
	}
}

