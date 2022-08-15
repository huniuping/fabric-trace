package main

import (
	"bytes"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"time"
)

type AssembleCc struct {
}

func (p *AssembleCc) Init(stub shim.ChaincodeStubInterface) pb.Response {
	init_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(init_timestamp.Seconds,int64(init_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 7 {
		return shim.Error("参数必须为7个")
	}
	assemble_product_id := args[0]      //组装产品标识号
	assemble_product_name := args[1]    //组装产品名称
	work_order_id := args[2]               //工单号
	assemble_line_id := args[3]   		//组装线编号
	date := args[4]    			        //组装完成日期
	process_list := args[5]  			//组件编号列表
	technology := args[6] 				//产品组装工艺文件
	txid :=stub.GetTxID()
	timestamp := ts
	value := fmt.Sprintf(`{"assemble_product_id":"%s","assemble_product_name":"%s","work_order_id":"%s","assemble_line_id":"%s","date":"%s","process_list":"%s","technology":"%s","timestamp":"%s","txid":"%s"}`,assemble_product_id,assemble_product_name,work_order_id,assemble_line_id,date,process_list,technology,timestamp,txid)
	err = stub.PutState(assemble_product_id, []byte(value))
	if err != nil {
		return shim.Error("数据初始化失败")
	}
	return shim.Success([]byte("初始化成功"))
}

func (p *AssembleCc) setValue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	if len(args) != 7 {
		return shim.Error("参数必须为7个")
	}
	assemble_product_id := args[0]      //组装产品标识号
	assemble_product_name := args[1]    //组装产品名称
	work_order_id := args[2]               //工单号
	assemble_line_id := args[3]   		//组装线编号
	date := args[4]    			        //组装完成日期
	process_list := args[5]  			//组件编号列表
	technology := args[6] 				//产品组装工艺文件
	txid :=stub.GetTxID()
	timestamp := ts
	rsp, err := stub.GetState(assemble_product_id)
	if string(rsp) != "" {
		return shim.Error("该产品标识号已存在")
	}
	value := fmt.Sprintf(`{"assemble_product_id":"%s","assemble_product_name":"%s","work_order_id":"%s","assemble_line_id":"%s","date":"%s","process_list":"%s","technology":"%s","timestamp":"%s","txid":"%s"}`,assemble_product_id,assemble_product_name,work_order_id,assemble_line_id,date,process_list,technology,timestamp,txid)
	err = stub.PutState(assemble_product_id, []byte(value))
	if err != nil {
		return shim.Error("新增产品失败")
	}
	return shim.Success([]byte("新增产品成功"))
}

func (p *AssembleCc) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品标识码，参数个数必须为1")
	}
	assemble_product_id := args[0]
	assemble_product_id_byte, err := stub.GetState(assemble_product_id)
	if err != nil {
		return shim.Error("没有产品信息")
	}
	return shim.Success(assemble_product_id_byte)
}

func (p *AssembleCc) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数为产品标识码，参数个数必须为1")
	}
	assemble_product_id := args[0]
	rsp,err:=stub.GetState(assemble_product_id)
	if string(rsp) == "" || err!=nil {
		return shim.Error("该产品编号不存在")
	}
	err = stub.DelState(assemble_product_id)
	if err != nil {
		return shim.Error("没有产品信息")
	}
	return shim.Success([]byte("删除产品成功"))
}

func (p *AssembleCc) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	set_timestamp,err :=stub.GetTxTimestamp()
	ts :=time.Unix(set_timestamp.Seconds,int64(set_timestamp.Nanos)).String()
	if err!= nil{
		return shim.Error("时间戳获取错误")
	}
	if len(args) != 7 {
		return shim.Error("参数必须为7个")
	}
	assemble_product_id := args[0]      //组装产品标识号
	assemble_product_name := args[1]    //组装产品名称
	work_order_id := args[2]               //工单号
	assemble_line_id := args[3]   		//组装线编号
	date := args[4]    			        //组装完成日期
	process_list := args[5]  			//组件编号列表
	technology := args[6] 				//产品组装工艺文件
	txid :=stub.GetTxID()
	timestamp := ts
	result, err :=stub.GetState(assemble_product_id)
	if err!=nil ||result ==nil {
		return shim.Error("未找到需要更新的记录")
	}
	value := fmt.Sprintf(`{"assemble_product_id":"%s","assemble_product_name":"%s","work_order_id":"%s","assemble_line_id":"%s","date":"%s","process_list":"%s","technology":"%s","timestamp":"%s","txid":"%s"}`,assemble_product_id,assemble_product_name,work_order_id,assemble_line_id,date,process_list,technology,timestamp,txid)
	err = stub.PutState(assemble_product_id, []byte(value))
	if err != nil {
		return shim.Error("更新组装信息记录失败")
	}
	return shim.Success([]byte("更新组装信息记录成功"))
}

func (p *AssembleCc) queryAssembleHistoryByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("给定的参数个数不符合要求")
	}
	assemble_product_id :=args[0]
	// 获取历史变更数据，返回一个迭代数组,
	//the historic value and associated transaction id and timestamp are returned.
	iterator, err := stub.GetHistoryForKey(assemble_product_id)
	if err != nil {
		return shim.Error("根据指定的标识号查询对应的历史变更数据失败")
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

func (p *AssembleCc) queryByGongdanId(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

func (p *AssembleCc) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
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
	if f == "queryAssembleHistoryByID" {
		return p.queryAssembleHistoryByID(stub, args)
	}
	if f == "queryByGongdanId" {
		return p.queryByGongdanId(stub, args)
	}
	return shim.Error("函数名称错误，只能是set、query、delete、update或者queryAssembleDetailByID")
}

func main() {
	err := shim.Start(new(AssembleCc))
	if err != nil {
		fmt.Println("启动fabric")
	}
}
