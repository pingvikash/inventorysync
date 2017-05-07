package BasukiChainCode

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

type inventorymaster struct {

}

type ItemDetails struct{	
	ItemID string `json:"itemid"`
	ItemDesc string `json:"itemdesc"`
	ItemPrice string `json:"itemprice"`
	UOM string `json:"uom"`
	ProductClass string `json:"prodClass"`
	CreatedBy string `json:"createdBy"`
	TotalQty string `json:"totQty"`
}

type Transaction struct{	
	TrxId string `json:"trxId"`
	TimeStamp string `json:"timeStamp"`
	ItemID string `json:"itemid"`
	SourceSystem string `json:"srcSystem"`
	Quantity string `json:"qty"`
	Trxntype string `json:"trxntype"`
	TrxnSubType string `json:"trxnSubType"`
	Remarks string `json:"remarks"`
}

type Inventory struct{
	ItemID string `json:"itemid"`
	Quantity string `json:"totQty"`
}

// Init initializes the smart contracts
func (t *inventorymaster) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Check if table already exists
	_, err := stub.GetTable("ItemDetails")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create application Table
	err = stub.CreateTable("ItemDetails", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "itemid", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "itemdesc", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "itemprice", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "uom", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "prodClass", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "createdBy", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "totQty", Type: shim.ColumnDefinition_STRING, Key: false},

	})
	if err != nil {
		return nil, errors.New("Failed creating ItemDetails.")
	}
	


	// Check if table already exists
	_, err = stub.GetTable("Transaction")
	if err == nil {
		// Table already exists; do not recreate
		return nil, nil
	}

	// Create transaction Table
	err = stub.CreateTable("Transaction", []*shim.ColumnDefinition{
		&shim.ColumnDefinition{Name: "trxId", Type: shim.ColumnDefinition_STRING, Key: true},
		&shim.ColumnDefinition{Name: "timeStamp", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "itemid", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "srcSystem", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "qty", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "trxntype", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "trxnSubType", Type: shim.ColumnDefinition_STRING, Key: false},
		&shim.ColumnDefinition{Name: "remarks", Type: shim.ColumnDefinition_STRING, Key: false},
	})
	if err != nil {
		return nil, errors.New("Failed creating Transaction Table.")
	}
		
	// setting up the users role
	stub.PutState("user_type1_1", []byte("pos"))
	stub.PutState("user_type1_2", []byte("jde"))
	stub.PutState("user_type1_3", []byte("fulfillment"))
	stub.PutState("user_type1_4", []byte("ecomm"))	
	
	return nil, nil
}


func (t *inventorymaster) registerItem(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

if len(args) != 7 {
			return nil, fmt.Errorf("Incorrect number of arguments. Expecting 7. Got: %d.", len(args))
		}
		
		itemid:=args[0]
		itemdesc:=args[1]
		itemprice:=args[2]
		uom:=args[3]
		prodClass:=args[4]
		totQty:=args[5]
		
		assignerOrg1, err := stub.GetState(args[6])
		assignerOrg := string(assignerOrg1)
		
		createdBy:=assignerOrg


		// Insert a row
		ok, err := stub.InsertRow("ItemDetails", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: itemid}},
				&shim.Column{Value: &shim.Column_String_{String_: itemdesc}},
				&shim.Column{Value: &shim.Column_String_{String_: itemprice}},
				&shim.Column{Value: &shim.Column_String_{String_: uom}},
				&shim.Column{Value: &shim.Column_String_{String_: prodClass}},
				&shim.Column{Value: &shim.Column_String_{String_: createdBy}},
				&shim.Column{Value: &shim.Column_String_{String_: totQty}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}
			
		return nil, nil

}



func (t *inventorymaster) updateInventory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8.")
	}

	trxId := args[0]
	timeStamp:=args[1]
	itemid := args[2]
	
	assignerOrg1, err := stub.GetState(args[3])
	assignerOrg := string(assignerOrg1)
	
	srcSystem := assignerOrg
	qty := args[4]
	trxntype := args[5]
	trxnSubType := args[6]
	remarks := args[7]
	
	newQty, _ := strconv.ParseInt(qty, 10, 0)
	
	//whether ADD_PENDING, DELETE_PENDING 
	/**if trxnSubType == "ADD_PENDING" || trxnSubType == "DELETE_PENDING"{
		newPoints = 0
	}**/
	

	// Get the row pertaining to this ffid
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: itemid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ItemDetails", columns)
	if err != nil {
		return nil, fmt.Errorf("Error: Failed retrieving user with itemid %s. Error %s", ffId, err.Error())
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		return nil, nil
	}

	currQty := row.Columns[6].GetString_()
	
	if trxntype=="add"{
		earlierQty:=row.Columns[6].GetString_()
		earlierInventory, _:=strconv.ParseInt(earlierQty, 10, 0)
		newInventory = strconv.Itoa(int(earlierInventory) + int(newQty))
	}else if trxntype=="delete"{
	
		earlierQty:=row.Columns[6].GetString_()
		earlierInventory, _:=strconv.ParseInt(earlierQty, 10, 0)
		newIntermediateInv := int(earlierInventory) - int(newQty)
		
		if newIntermediateInv < 0 {
			return nil, errors.New("can't deduct as the resulting inventory becoming less than zero.")
		}
		newInventory = strconv.Itoa(int(earlierInventory) - int(newQty))
	}else{
		return nil, fmt.Errorf("Error: Failed retrieving user with itemid %s. Error %s", itemid, err.Error())
	}
	
	
	//End- Check that the currentStatus to newStatus transition is accurate
	// Delete the row pertaining to this ffid
	err = stub.DeleteRow(
		"ItemDetails",
		columns,
	)
	if err != nil {
		return nil, errors.New("Failed deleting row.")
	}

	
	//ffId := row.Columns[0].GetString_()
	
	itemid := row.Columns[1].GetString_()
	itemdesc := row.Columns[2].GetString_()
	itemprice := row.Columns[3].GetString_()
	uom := row.Columns[4].GetString_()
	prodClass := row.Columns[5].GetString_()
	createdBy := row.Columns[6].GetString_()
	totQty := newInventory


		// Insert a row
		ok, err := stub.InsertRow("ItemDetails", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: itemid}},
				&shim.Column{Value: &shim.Column_String_{String_: itemdesc}},
				&shim.Column{Value: &shim.Column_String_{String_: itemprice}},
				&shim.Column{Value: &shim.Column_String_{String_: uom}},
				&shim.Column{Value: &shim.Column_String_{String_: prodClass}},
				&shim.Column{Value: &shim.Column_String_{String_: createdBy}},
				&shim.Column{Value: &shim.Column_String_{String_: totQty}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}

		
		//inserting the transaction
		
		// Insert a row
		ok, err = stub.InsertRow("Transaction", shim.Row{
			Columns: []*shim.Column{
				&shim.Column{Value: &shim.Column_String_{String_: trxId}},
				&shim.Column{Value: &shim.Column_String_{String_: timeStamp}},
				&shim.Column{Value: &shim.Column_String_{String_: itemid}},
				&shim.Column{Value: &shim.Column_String_{String_: srcSystem}},
				&shim.Column{Value: &shim.Column_String_{String_: qty}},
				&shim.Column{Value: &shim.Column_String_{String_: trxntype}},
				&shim.Column{Value: &shim.Column_String_{String_: trxnSubType}},
				&shim.Column{Value: &shim.Column_String_{String_: remarks}},
			}})

		if err != nil {
			return nil, err 
		}
		if !ok && err == nil {
			return nil, errors.New("Row already exists.")
		}		
	return nil, nil

}

func (t *inventorymaster) getQty(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting itemid to query")
	}

	ffId := args[0]
	

	// Get the row pertaining to this ffId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: itemid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ItemDetails", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the itemid " + itemid + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the itemid " + itemid + "\"}"
		return nil, errors.New(jsonResp)
	}

	
	
	res2E := Inventory{}
	
	res2E.ItemID = row.Columns[6].GetString_()
	res2E.Quantity = row.Columns[6].GetString_()
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}


func (t *inventorymaster) getTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting itemid to query")
	}

	itemid := args[0]
	assignerRole := args[1]

	var columns []shim.Column

	rows, err := stub.GetRows("Transaction", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	
	assignerOrg1, err := stub.GetState(assignerRole)
	assignerOrg := string(assignerOrg1)
	
		
	res2E:= []*Transaction{}	
	
	for row := range rows {		
		newApp:= new(Transaction)
		newApp.TrxId = row.Columns[0].GetString_()
		newApp.TimeStamp = row.Columns[1].GetString_()
		newApp.ItemID = row.Columns[2].GetString_()
		newApp.SourceSystem = row.Columns[3].GetString_()
		newApp.Quantity = row.Columns[4].GetString_()
		newApp.Trxntype = row.Columns[5].GetString_()
		newApp.TrxnSubType = row.Columns[6].GetString_()
		newApp.Remarks = row.Columns[7].GetString_()
		
		if newApp.ItemID == itemid && newApp.SourceSystem == assignerOrg{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}



func (t *inventorymaster) getAllTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting itemid to query")
	}

	ffId := args[0]

	var columns []shim.Column

	rows, err := stub.GetRows("Transaction", columns)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve row")
	}
	
	//assignerOrg1, err := stub.GetState(assignerRole)
	//assignerOrg := string(assignerOrg1)
	
		
	res2E:= []*Transaction{}	
	
	for row := range rows {		
		newApp:= new(Transaction)
		newApp.TrxId = row.Columns[0].GetString_()
		newApp.TimeStamp = row.Columns[1].GetString_()
		newApp.ItemID = row.Columns[2].GetString_()
		newApp.SourceSystem = row.Columns[3].GetString_()
		newApp.Quantity = row.Columns[4].GetString_()
		newApp.Trxntype = row.Columns[5].GetString_()
		newApp.TrxnSubType = row.Columns[6].GetString_()
		newApp.Remarks = row.Columns[7].GetString_()
		
		if newApp.ItemID == itemid{
		res2E=append(res2E,newApp)		
		}				
	}
	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}


func (t *inventorymaster) getItem(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting itemid to query")
	}

	itemid := args[0]
	

	// Get the row pertaining to this ffId
	var columns []shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: itemid}}
	columns = append(columns, col1)

	row, err := stub.GetRow("ItemDetails", columns)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get the data for the application " + itemid + "\"}"
		return nil, errors.New(jsonResp)
	}

	// GetRows returns empty message if key does not exist
	if len(row.Columns) == 0 {
		jsonResp := "{\"Error\":\"Failed to get the data for the application " + itemid + "\"}"
		return nil, errors.New(jsonResp)
	}

	
	res2E := ItemDetails{}
	
	res2E.ItemID = row.Columns[0].GetString_()
	res2E.ItemDesc = row.Columns[1].GetString_()
	res2E.ItemPrice = row.Columns[2].GetString_()
	res2E.UOM = row.Columns[3].GetString_()
	res2E.ProductClass = row.Columns[4].GetString_()
	res2E.CreatedBy = row.Columns[5].GetString_()
	res2E.TotalQty = row.Columns[6].GetString_()

	
    mapB, _ := json.Marshal(res2E)
    fmt.Println(string(mapB))
	
	return mapB, nil

}




// Invoke invokes the chaincode
func (t *inventorymaster) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "registerItem" {
		t := inventorymaster{}
		return t.registerItem(stub, args)	
	} else if function == "updateInventory" { 
		t := inventorymaster{}
		return t.updateInventory(stub, args)
	}

	return nil, errors.New("Invalid invoke function name.")

}

// query queries the chaincode
func (t *inventorymaster) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getQty" {
		t := inventorymaster{}
		return t.getMile(stub, args)		
	} else if function == "getTransaction" { 
		t := inventorymaster{}
		return t.getTransaction(stub, args)
	}else if function == "getAllTransaction" { 
		t := inventorymaster{}
		return t.getAllTransaction(stub, args)
	} else if function == "getItem" { 
		t := inventorymaster{}
		return t.getUser(stub, args)
	}
	
	return nil, nil
}

func main() {
	primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(inventorymaster))
	if err != nil {
		fmt.Printf("Error starting FFP: %s", err)
	}
} 
