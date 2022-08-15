package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	//"strconv"
	"time"
)

type DrawingCc struct {
}

func (p *DrawingCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	drawing_id := args[0]   		  //图纸标识号
	drawing_name := args[1]   		  //图纸名称
	drawing_file := args[2]   		  //图纸文件哈希
	contractId := args[3]   		  //合同
	technology:= args[4]   		  	  //工艺名称
	txid:=stub.GetTxID()
	timestamp := ts
	value := fmt.Sprintf(`{"drawing_id":"%s","drawing_name":"%s","drawing_file":"%s","contractId":"%s","technology":"%s","timestamp":"%s","txid":"%s"}`,drawing_id,drawing_name,drawing_file,contractId,technology,timestamp,txid)
	err = stub.PutState(drawing_id, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *DrawingCc) setValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	drawing_id := args[0]   		  //图纸标识号
	drawing_name := args[1]   		  //图纸名称
	drawing_file := args[2]   		  //图纸文件哈希
	contractId := args[3]   		  //合同
	technology:= args[4]   		  	  //工艺名称
	txid:=stub.GetTxID()
	timestamp := ts
	rsp, err := stub.GetState(drawing_id)
	if string(rsp) != "" {
		return shim.Error("该图纸标识号已存在")
	}
	value := fmt.Sprintf(`{"drawing_id":"%s","drawing_name":"%s","drawing_file":"%s","contractId":"%s","technology":"%s","timestamp":"%s","txid":"%s"}`,drawing_id,drawing_name,drawing_file,contractId,technology,timestamp,txid)
	err = stub.PutState(drawing_id, []byte(value))
	if err != nil {
		return shim.Error("新增图纸失败")
	}
	return shim.Success([]byte("新增图纸成功"))
}

func (p *DrawingCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为图纸标识码，参数个数必须为1")
	}
	drawing_id := args[0]
	drawing_id_byte, err := stub.GetState(drawing_id)
	if err != nil {
		return shim.Error("没有图纸信息")
	}
	return shim.Success(drawing_id_byte)
}

func (p *DrawingCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为图纸标识码，参数个数必须为1")
	}
	drawing_id := args[0]
	rsp,err:=stub.GetState(drawing_id)
	if string(rsp) == "" || err!=nil {
		return shim.Error("该图纸编号不存在")
	}
	err = stub.DelState(drawing_id)
	if err != nil {
		return shim.Error("没有图纸信息")
	}
	return shim.Success([]byte("删除图纸成功"))
}

func (p *DrawingCc) queryByContractId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为合同名称编号，个数必须为1")
	}
	//contractName := args[0]
	contractId := fmt.Sprintf("{\"selector\":{\"contractId\":\"%s\"}}", args[0])
	resultsIterator, err :=stub.GetQueryResult(contractId)
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
		//key :=queryResponse.Key
		//value :=queryResponse.Value
		//kv :=key+value
		//buffer.WriteString(queryResponse.Key+string(queryResponse.Value))
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *DrawingCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp, _ :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if len(args) != 5 {
		return shim.Error("参数必须为5个")
	}
	drawing_id := args[0]   		  //图纸标识号
	drawing_name := args[1]   		  //图纸名称
	drawing_file := args[2]   		  //图纸文件哈希
	contractId := args[3]   		  //合同
	technology:= args[4]   		  	  //工艺名称
	txid:=stub.GetTxID()
	timestamp := ts
	result, err :=stub.GetState(drawing_id)

	if err!=nil ||result ==nil {
		return shim.Error("未找到需要更新的图纸")
	}
	value := fmt.Sprintf(`{"drawing_id":"%s","drawing_name":"%s","drawing_file":"%s","contractId":"%s","technology":"%s","timestamp":"%s","txid":"%s"}`,drawing_id,drawing_name,drawing_file,contractId,technology,timestamp,txid)
	err = stub.PutState(drawing_id, []byte(value))
	if err != nil {
		return shim.Error("更新图纸失败")
	}
	return shim.Success([]byte("更新图纸成功"))
}

func (p *DrawingCc) queryDrawingHistoryByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	drawing_id :=args[0]
	iterator, err := stub.GetHistoryForKey(drawing_id)
	if err != nil {
		return shim.Error("根据指定的标识号查询对应的历史变更数据失败")
	}
	defer iterator.Close()
	var buffer bytes.Buffer
	bArrayMemberAlreadyWritten := false
	for iterator.HasNext() {
		hisData, err := iterator.Next()//获取当前i的的迭代器数据
		if err != nil {
			return shim.Error("获取sco的历史变更数据失败")
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


func (p *DrawingCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
	if f == "queryDrawingHistoryByID" {
		return p.queryDrawingHistoryByID(stub, args)
	}
	return shim.Error("函数名称错误")
}

func main() {
	err := shim.Start(new(DrawingCc))
	if err != nil {
		fmt.Println("启动fabric")
	}
}

