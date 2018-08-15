/*
 * The sample smart contract for documentation topic:
 * cross border funds transfer
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the hospital structure, with 3 properties.  Structure tags are used by encoding/json library
type Hospital struct {
	HospitalID string  `json:"hospitalID"`
	Name       string  `json:"name"`
	Country    string  `json:"country"`
	Balance    float64 `json:"Balance"`
}

// Define the doctor structure, with 3 properties.  Structure tags are used by encoding/json library
type Doctor struct {
	DoctorID   string  `json:"doctorID"`
	Name       string  `json:"name"`
	HospitalID string  `json:"hospitalID"`
	Balance    float64 `json:"Balance"`
}

// Define the patient structure, with 3 properties.  Structure tags are used by encoding/json library
type Patient struct {
	PatientID  string  `json:"patientID"`
	Name       string  `json:"name"`
	ReportID   string  `json:"reportID"`
	HospitalID string  `json:"hospitalID"`
	Balance    float64 `json:"balance"`
}

// Define the report structure, with 3 properties.  Structure tags are used by encoding/json library
type Report struct {
	ReportID   string  `json:"reportID"`
	PatientID  string  `json:"patientID"`
	HospitalID string  `json:"hospitalID"`
	Fee        float64 `json:"fee"`
}

/*
 * The Init method is called when the Smart Contract "health" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract ""
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryAll" { //return all the assets on the ledger
		return s.queryAll(APIstub, args)
	} else if function == "query" { //single hospital or doctor or patient or report
		return s.query(APIstub, args)
	} else if function == "createHospital" {
		return s.createHospital(APIstub, args)
	} else if function == "createDoctor" {
		return s.createDoctor(APIstub, args)
	} else if function == "createPatient" {
		return s.createPatient(APIstub, args)
	} else if function == "createReport" {
		return s.createReport(APIstub, args)
	} else if function == "queryAllHospitals" { //return all the hospitals on the ledger
		return s.queryAllHospitals(APIstub)
	} else if function == "queryAllDoctors" { //return all the doctors on the ledger
		return s.queryAllDoctors(APIstub)
	} else if function == "queryAllPatients" { //return all the patients on the ledger
		return s.queryAllPatients(APIstub)
	} else if function == "queryAllReports" { //return all the reports on the ledger
		return s.queryAllReports(APIstub)
	} else if function == "transferPatient" { //change HospitalID
		return s.transferPatient(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/** ----------------------------------------------------------------------**/
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	hospitals := []Hospital{
		{HospitalID: "H001", Name: "HOSPTITAL_1", Country: "INDIA", Balance: 1000000.0},
		{HospitalID: "H002", Name: "HOSPTITAL_2", Country: "INDIA", Balance: 1000000.0},
		{HospitalID: "H003", Name: "HOSPTITAL_3", Country: "INDIA", Balance: 1000000.0},
	}

	doctors := []Doctor{
		{DoctorID: "D001", Name: "DOCTOR_1", HospitalID: "H001", Balance: 100000.0},
		{DoctorID: "D002", Name: "DOCTOR_2", HospitalID: "H001", Balance: 100000.0},
		{DoctorID: "D003", Name: "DOCTOR_3", HospitalID: "H002", Balance: 100000.0},
		{DoctorID: "D004", Name: "DOCTOR_4", HospitalID: "H002", Balance: 100000.0},
		{DoctorID: "D005", Name: "DOCTOR_5", HospitalID: "H003", Balance: 100000.0},
		{DoctorID: "D006", Name: "DOCTOR_6", HospitalID: "H003", Balance: 100000.0},
	}

	patients := []Patient{
		{PatientID: "P001", Name: "PATIENT_1", ReportID: "R001", HospitalID: "H001", Balance: 10000.0},
		{PatientID: "P002", Name: "PATIENT_2", ReportID: "R002", HospitalID: "H001", Balance: 10000.0},
		{PatientID: "P003", Name: "PATIENT_3", ReportID: "R003", HospitalID: "H002", Balance: 10000.0},
		{PatientID: "P004", Name: "PATIENT_4", ReportID: "R004", HospitalID: "H002", Balance: 10000.0},
		{PatientID: "P005", Name: "PATIENT_5", ReportID: "R005", HospitalID: "H003", Balance: 10000.0},
		{PatientID: "P006", Name: "PATIENT_6", ReportID: "R006", HospitalID: "H003", Balance: 10000.0},
	}

	reports := []Report{
		{ReportID: "R001", PatientID: "P001", HospitalID: "H001", Fee: 1000.0},
		{ReportID: "R002", PatientID: "P002", HospitalID: "H001", Fee: 500.0},
		{ReportID: "R003", PatientID: "P003", HospitalID: "H002", Fee: 800.0},
		{ReportID: "R004", PatientID: "P004", HospitalID: "H002", Fee: 1000.0},
		{ReportID: "R005", PatientID: "P005", HospitalID: "H003", Fee: 600.0},
		{ReportID: "R006", PatientID: "P006", HospitalID: "H003", Fee: 1000.0},
	}

	writeHospitalToLedger(APIstub, hospitals)
	writeDoctorToLedger(APIstub, doctors)
	writePatientToLedger(APIstub, patients)
	writeReportToLedger(APIstub, reports)

	return shim.Success(nil)
}

/** -------------------------------------writeHospitalToLedger---------------------------------------------*/

func writeHospitalToLedger(APIStub shim.ChaincodeStubInterface, hospitals []Hospital) sc.Response {
	for i := 0; i < len(hospitals); i++ {
		key := hospitals[i].HospitalID
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(hospitals[i])
			err := APIStub.PutState(key, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Hospital with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}
	}
	return shim.Success(nil)
}

/** --------------------------------------writeDoctorToLedger--------------------------------------------*/

func writeDoctorToLedger(APIStub shim.ChaincodeStubInterface, doctors []Doctor) sc.Response {

	for i := 0; i < len(doctors); i++ {
		key := doctors[i].DoctorID
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(doctors[i])
			err := APIStub.PutState(key, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Doctor with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}

	}
	return shim.Success(nil)
}

/** ---------------------------------------writePatientToLedger-------------------------------------------*/

func writePatientToLedger(APIStub shim.ChaincodeStubInterface, patients []Patient) sc.Response {
	for i := 0; i < len(patients); i++ {
		key := patients[i].PatientID
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(patients[i])
			err := APIStub.PutState(key, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Patient with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}
	}
	return shim.Success(nil)
}

/** -----------------------------------------writeReportToLedger-----------------------------------------*/

func writeReportToLedger(APIStub shim.ChaincodeStubInterface, reports []Report) sc.Response {
	for i := 0; i < len(reports); i++ {
		key := reports[i].ReportID
		chkBytes, _ := APIStub.GetState(key)
		if chkBytes == nil { //only add if it is not already present
			asBytes, _ := json.Marshal(reports[i])
			err := APIStub.PutState(key, asBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		} else {
			msg := " Report with key:" + key + " already exists.. skipping ......."
			return shim.Error(msg)
		}
	}
	return shim.Success(nil)
}

/** ----------------------------------------queryAll-------------------------------------------------------*/

func (s *SmartContract) queryAll(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments for querying all assets. Expecting 1")
	}
	//collection := args[0]
	//startKey := collection + "0"
	//endKey := collection + "99"

	resultsIterator, err := APIstub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString("\n,")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}\n")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("\n]")

	fmt.Printf("- queryAll:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/** ---------------------------------query-------------------------------------**/

func (s *SmartContract) query(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	asBytes, _ := APIstub.GetState(args[0])
	return shim.Success(asBytes)
}

/** ----------------------------------------------------------------------------------------------
create hospital needs 4 args
HospitalID     string  `json:"hospitalID"`
Name           string  `json:"name"`
Country        string  `json:"country"`
Balance        float64 `json:"balance"`
args: ['H004', 'HOSPITAL_4', 'INDIA', '1000000.0'],
*/
func (s *SmartContract) createHospital(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments for creating a hospital. Expecting 4")
	}
	balance, _ := strconv.ParseFloat(args[3], 64)
	hospitals := []Hospital{Hospital{HospitalID: args[0], Name: args[1], Country: args[2], Balance: balance}}

	writeHospitalToLedger(APIstub, hospitals)
	return shim.Success(nil)
}

/** ----------------------------------------------------------------------------------------------
ceate doctor needs 4 args
DoctorID      string  `json:"doctorID"`
Name          string  `json:"name"`
HospitalID    string  `json:"hospitalID"`
Balance       float64 `json:"balance"`
args: ["D007", "DOCTOR_7", "H004", "100000.0"],
*/
func (s *SmartContract) createDoctor(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments for creating a Doctor. Expecting 4")
	}
	balance, _ := strconv.ParseFloat(args[3], 64)
	doctors := []Doctor{Doctor{DoctorID: args[0], Name: args[1], HospitalID: args[2], Balance: balance}}

	writeDoctorToLedger(APIstub, doctors)
	return shim.Success(nil)
}

/**----------------------------------------------------------------------------------------------
createPatient needs 5 args
	PatientID     string  `json:"patientID"`
	Name          string  `json:"name"`
	ReportID      string  `json:"reportID"`
	HospitalID    string  `json:"hospitalID"`
	Balance       float64 `json:"balance"`
	["P007", "PATIENT_7",  "R007", "H004", "100000.0"],
*/
func (s *SmartContract) createPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments for creating a patient. Expecting 5")
	}
	balance, _ := strconv.ParseFloat(args[4], 64)
	patients := []Patient{Patient{PatientID: args[0], Name: args[1], ReportID: args[2], HospitalID: args[3], Balance: balance}}

	writePatientToLedger(APIstub, patients)

	return shim.Success(nil)
}

/** ----------------------------------------------------------------------------------------------
ceate report needs 4 args
ReportID      string  `json:"reportID"`
PatientID     string  `json:"patientID"`
HospitalID    string  `json:"hospitalID"`
Fee           float64 `json:"fee"`
args: ['R007', 'P007', 'H007', '100000.0'],
*/
func (s *SmartContract) createReport(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments for creating a Report. Expecting 4")
	}
	fee, _ := strconv.ParseFloat(args[3], 64)
	reports := []Report{Report{ReportID: args[0], PatientID: args[1], HospitalID: args[2], Fee: fee}}

	writeReportToLedger(APIstub, reports)

	//get FROM patient from ledger
	fromPatientAsBytes, _ := APIstub.GetState(args[1])
	fromPatient := Patient{}

	json.Unmarshal(fromPatientAsBytes, &fromPatient)
	fromPatientPatientID := fromPatient.PatientID
	fromBalance := float64(fromPatient.Balance)

	//check if patient has enough balance to cover the payment
	if fromBalance < fee {
		errMsg := "Insufficent funds in patient: Patient ID: " + fromPatientPatientID
		return shim.Error(errMsg)
	}

	//get TO hospital from ledger
	toHospitalAsBytes, _ := APIstub.GetState(args[2])
	toHospital := Hospital{}

	json.Unmarshal(toHospitalAsBytes, &toHospital)
	toBalance := toHospital.Balance

	//reduce FROM patient balance by fee amount
	fromPatient.Balance = fromBalance - fee

	//increase TO hospital balance by fee amount
	toHospital.Balance = toBalance + fee

	//write all changed assets to the ledger
	fromPatientAsBytes, _ = json.Marshal(fromPatient)
	err := APIstub.PutState(args[1], fromPatientAsBytes)
	if err != nil {
		return shim.Error("Error writing updates to FROM patient account " + fromPatient.PatientID)
	}

	toHospitalAsBytes, _ = json.Marshal(toHospital)
	err = APIstub.PutState(args[2], toHospitalAsBytes)
	if err != nil {
		return shim.Error("Error writing updates to TO hospital account " + toHospital.HospitalID)
	}

	if err == nil {
		fmt.Println("~~~~~~~~~~~~~~~~~ Successfully created report~~~~~~~~~~~~~~~~~")
	}

	return shim.Success(nil)
}

/**------------------------------------queryAllHospitals-------------------------------------------*/

func (s *SmartContract) queryAllHospitals(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "H001"
	endKey := "H999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllHospitals:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/**------------------------------------queryAllDoctors-------------------------------------------*/

func (s *SmartContract) queryAllDoctors(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "D001"
	endKey := "D999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllDoctors:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/**------------------------------------queryAllPatients-------------------------------------------*/

func (s *SmartContract) queryAllPatients(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "P001"
	endKey := "P999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllPatients:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/**------------------------------------queryAllReports-------------------------------------------*/

func (s *SmartContract) queryAllReports(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "R001"
	endKey := "R999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllReports:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/**------------------------------------transferPatient-------------------------------------------*/

func (s *SmartContract) transferPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	patient := Patient{}

	json.Unmarshal(patientAsBytes, &patient)
	patient.HospitalID = args[1]

	patientAsBytes, _ = json.Marshal(patient)
	APIstub.PutState(args[0], patientAsBytes)

	return shim.Success(nil)
}

/**----------------------------------------------------------------------------------------------
The main function is only relevant in unit test mode. Only included here for completeness.
*/

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}

	fmt.Println("successfully initialized smart contract")

}
