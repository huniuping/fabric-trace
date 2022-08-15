package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

type ProcessCc struct {
}

func (p *ProcessCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 7 {
		return shim.Error("参数必须为7个")
	}
	process_product_id := args[0]     //半成品产品标识
	product_name := args[1]           //半成品产品名字
	work_order_id := args[2]           //工单号
	producted_id := args[3]       	  //所用于组装产品id
	technology := args[4] 			  //工艺参数文件
	technology_sequence := args[5] 			  //工艺顺序
	sequential_aggregate_signature := args[6]  //有序聚合签名
	txid:=stub.GetTxID()
	timestamp := ts
	value := fmt.Sprintf(`{"process_product_id":"%s","product_name":"%s","work_order_id":"%s","producted_id":"%s","technology":"%s","technology_sequence":"%s","sequential_aggregate_signature":"%s","timestamp":"%s","txid":"%s"}`,process_product_id,product_name,work_order_id,producted_id,technology,technology_sequence,sequential_aggregate_signature,timestamp,txid)
	err = stub.PutState(process_product_id, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}


func (p *ProcessCc) setValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	if len(args) != 7 {
		return shim.Error("参数必须为7个")
	}
	process_product_id := args[0]     //半成品产品标识
	product_name := args[1]           //半成品产品名字
	work_order_id := args[2]           //工单号
	producted_id := args[3]       	  //所用于组装产品id
	technology := args[4] 			  //工艺参数文件
	technology_sequence := args[5] 			  //工艺顺序
	sequential_aggregate_signature := args[6]  //有序聚合签名
	txid:=stub.GetTxID()
	timestamp := ts
	rsp, err := stub.GetState(process_product_id)
	if string(rsp) != "" {
		return shim.Error("该产品标识号已存在")
	}
	value := fmt.Sprintf(`{"process_product_id":"%s","product_name":"%s","work_order_id":"%s","producted_id":"%s","technology":"%s","technology_sequence":"%s","sequential_aggregate_signature":"%s","timestamp":"%s","txid":"%s"}`,process_product_id,product_name,work_order_id,producted_id,technology,technology_sequence,sequential_aggregate_signature,timestamp,txid)
	err = stub.PutState(process_product_id, []byte(value))
	if err != nil {
		return shim.Error("新增半成品产品失败")
	}
	return shim.Success([]byte("新增半成品产品成功"))
}


func (p *ProcessCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品标识码，参数个数必须为1")
	}
	process_product_id := args[0]
	process_product_id_byte, err := stub.GetState(process_product_id)
	if err != nil {
		return shim.Error("没有产品信息")
	}
	return shim.Success(process_product_id_byte)
}

func (p *ProcessCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品标识码，参数个数必须为1")
	}
	process_product_id := args[0]
	rsp,err:=stub.GetState(process_product_id)
	if string(rsp) == "" || err!=nil {
		return shim.Error("该产品编号不存在")
	}
	err = stub.DelState(process_product_id)
	if err != nil {
		return shim.Error("没有产品信息")
	}
	return shim.Success([]byte("删除半成品产品成功"))
}

func (p *ProcessCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	if len(args) != 7 {
		return shim.Error("参数必须为7个")
	}
	process_product_id := args[0]     //半成品产品标识
	product_name := args[1]           //半成品产品名字
	work_order_id := args[2]           //工单号
	producted_id := args[3]       	  //所用于组装产品id
	technology := args[4] 			  //工艺参数文件
	technology_sequence := args[5] 			  //工艺顺序
	sequential_aggregate_signature := args[6]  //有序聚合签名
	txid:=stub.GetTxID()
	timestamp := ts
	result, err :=stub.GetState(process_product_id)
	if err!=nil ||result ==nil {
		return shim.Error("未找到需要更新的记录")
	}
	value := fmt.Sprintf(`{"process_product_id":"%s","product_name":"%s","work_order_id":"%s","producted_id":"%s","technology":"%s","technology_sequence":"%s","sequential_aggregate_signature":"%s","timestamp":"%s","txid":"%s"}`,process_product_id,product_name,work_order_id,producted_id,technology,technology_sequence,sequential_aggregate_signature,timestamp,txid)
	err = stub.PutState(process_product_id, []byte(value))
	if err != nil {
		return shim.Error("更新半成品失败")
	}
	return shim.Success([]byte("更新半成品成功"))
}

func (p *ProcessCc) queryByProductId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为半成品组装母体编号，个数必须为1")
	}
	producted_id := fmt.Sprintf("{\"selector\":{\"producted_id\":\"%s\"}}", args[0])
	resultsIterator, err :=stub.GetQueryResult(producted_id)
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
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *ProcessCc) queryByGongdanId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为工单编号，个数必须为1")
	}
	work_order_id := fmt.Sprintf("{\"selector\":{\"work_order_id\":\"%s\"}}", args[0])
	resultsIterator, err :=stub.GetQueryResult(work_order_id)
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

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}
	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (p *ProcessCc) queryProcessHistoryByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	process_product_id :=args[0]
	// 获取历史变更数据，返回一个迭代数组,
	//the historic value and associated transaction id and timestamp are returned.
	iterator, err := stub.GetHistoryForKey(process_product_id)
	if err != nil {
		return shim.Error("根据指定的标识号查询对应的历史变更数据失败")
	}
	if iterator==nil{
		return shim.Error("iterator is nil")
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

func (p *ProcessCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
	if f == "queryByProductId" {
		return p.queryByProductId(stub, args)
	}
	if f == "queryByGongdanId" {
		return p.queryByGongdanId(stub, args)
	}
	if f == "queryProcessHistoryByID" {
		return p.queryProcessHistoryByID(stub, args)
	}
	return shim.Error("函数名称错误，只能是set、query、delete、update、queryByProductId或queryProcessDetailByID")
}

func main() {
	err := shim.Start(new(ProcessCc))
	if err != nil {
		fmt.Println("启动fabric")
	}
}
