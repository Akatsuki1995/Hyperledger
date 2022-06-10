package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Identification : 
type Identification struct {
	Tattoo_number int `json:"Tattoo_number"`
	Microship_number int `json:"Microship_number"`
	Date string `json:"Date"`
	Location string `json:"Location"`
}
// Define the Pet structure, with 4 properties.  Structure tags are used by encoding/json library
type Pet struct {
	Name string `json:"Name"`
	Species string `json:"Species"`
	Breed string `json:"Breed"`
	Owner string `json:"Owner"`
}
type petPrivateDetails struct {
	Owner string `json:"Owner"`
	Date_of_birth string `json:"Date_of_birth"`
	Sex string `json:"Sex"`
	Colour string `json:"Colour"`
	Distinguishing_marks string `json:"Distinguishing_marks"`
	Known_allergies string `json:"Known_allergies"`
	Price string `json:"Price"` 
}


// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("diggipet_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryPet":
		return s.queryPet(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createPet":
		return s.createPet(APIstub, args)
	case "queryAllPets":
		return s.queryAllPets(APIstub)
	case "changePetOwner":
		return s.changePetOwner(APIstub,args)
	case "getHistoryForAsset":
		return s.getHistoryForAsset(APIstub, args)
	case "queryPetsByOwner":
		return s.queryPetsByOwner(APIstub, args)
	case "restictedMethod":
		return s.restictedMethod(APIstub, args)
	case "test":
		return s.test(APIstub, args)
	case "createPrivatePet":
		return s.createPrivatePet(APIstub, args)
	case "readPrivatePet":
		return s.readPrivatePet(APIstub,args)
	case "updatePrivateData":
		return s.updatePrivateData(APIstub, args)
	case "readPetPrivateDetails":
		return s.readPetPrivateDetails(APIstub, args)
	case "createPrivatePetImplicitForDiggipet":
		return s.createPrivatePetImplicitForDiggipet(APIstub, args)
	case "createPrivatePetImplicitForVetOrg":
		return s.createPrivatePetImplicitForVetOrg(APIstub, args)	
	case "createPrivatePetImplicitForBreederOrg":
		return s.createPrivatePetImplicitForBreederOrg(APIstub, args)

	case "queryPrivateDataHash":
		return s.queryPrivateDataHash(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

	// return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryPet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	petAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(petAsBytes)
}
func (s *SmartContract) queryPerson(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	personAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(personAsBytes)
}

func (s *SmartContract) readPrivatePet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// collectionPets, collectionPetPrivateDetails, _implicit_org_DiggipetMSP, _implicit_org_VetorgMSP, _implicit_org_BreedersorgMSP
	petAsBytes, err := APIstub.GetPrivateData(args[0], args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[1] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if petAsBytes == nil {
		jsonResp := "{\"Error\":\"Pet private details does not exist: " + args[1] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(petAsBytes)
}

func (s *SmartContract) readPetPrivateDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	petAsBytes, err := APIstub.GetPrivateData("collectionPetPrivateDetails", args[0])

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[0] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if petAsBytes == nil {
		jsonResp := "{\"Error\":\"Marble private details does not exist: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(petAsBytes)
}

func (s *SmartContract) test(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	petAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(petAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	pets := []Pet{
		Pet{Name: "Lara", Species: "Cat", Breed: "Persan", Owner: "Hakam"},
		Pet{Name: "Jack", Species: "Dog", Breed: "Labrador", Owner: "Peter"},
		Pet{Name: "Kika", Species: "Cat", Breed: "Siamois", Owner: "Emil"},
		Pet{Name: "Angry", Species: "Cat", Breed: "Ragdoll", Owner: "Robinson"},
		Pet{Name: "Black", Species: "Dog", Breed: "Caniche", Owner: "Joel"},
		Pet{Name: "Snow", Species: "Dog", Breed: "Boxer", Owner: "Brad"},
		Pet{Name: "Blacky", Species: "Cat", Breed: "Maine_Coon", Owner: "Hakam"},
		Pet{Name: "Lucy", Species: "Dog", Breed: "Golden_Retriever", Owner: "Helen"},
		Pet{Name: "Cesar", Species: "Cat", Breed: "Bengal", Owner: "Adriana"},
	}

	i := 0
	for i < len(pets) {
		petAsBytes, _ := json.Marshal(pets[i])
		APIstub.PutState("PET"+strconv.Itoa(i), petAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createPrivatePet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type petTransientInput struct {
		Name string `json:"Name"` //the fieldtags are needed to keep case from bouncing around
		Species string `json:"Species"`
		Breed string `json:"Breed"`
		Owner string `json:"Owner"`
		Date_of_birth string `json:"Date_of_birth"`
		Sex string `json:"Sex"`
		Colour string `json:"Colour"`
		Distinguishing_marks string `json:"Distinguishing_marks"`
		Known_allergies string `json:"Known_allergies"`
		Price string `json: "Price"`
		Key   string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	petDataAsBytes, ok := transMap["pet"]
	if !ok {
		return shim.Error("pet must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(petDataAsBytes))

	if len(petDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var petInput petTransientInput
	err = json.Unmarshal(petDataAsBytes, &petInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(petDataAsBytes) + "Error is : " + err.Error())
	}

	logger.Infof("3333")

	if len(petInput.Key) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Name) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Species) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Breed) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Owner) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Date_of_birth) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Sex) == 0 {
		return shim.Error("Key field must be a non-empty string")
	}
	if len(petInput.Colour) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if len(petInput.Distinguishing_marks) == 0 {
		return shim.Error("model field must be a non-empty string")
	}
	if len(petInput.Known_allergies) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if len(petInput.Price) == 0 {
		return shim.Error("price field must be a non-empty string")
	}

	logger.Infof("444444")

	// ==== Check if pet already exists ====
	petAsBytes, err := APIstub.GetPrivateData("collectionPets", petInput.Key)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if petAsBytes != nil {
		fmt.Println("This pet already exists: " + petInput.Key)
		return shim.Error("This pet already exists: " + petInput.Key)
	}

	logger.Infof("55555")
	if err != nil {
        return shim.Error("Cannot get ID" + err.Error())
    }
	mspid, err := cid.GetMSPID(APIstub)
	err = verifyClientOrgMatchesPeerOrg(mspid)
    if err != nil {
        return shim.Error(err.Error())
    }
	var pet = Pet{Name: petInput.Name, Species: petInput.Species, Breed: petInput.Breed, Owner: petInput.Owner}

	petAsBytes, err = json.Marshal(pet)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutPrivateData("collectionPets", petInput.Key, petAsBytes)
	if err != nil {
		logger.Infof("6666666")
		return shim.Error(err.Error())
	}

	petPrivateDetails := &petPrivateDetails{Owner: petInput.Owner, Date_of_birth: petInput.Date_of_birth, Sex: petInput.Sex, Colour: petInput.Colour, Distinguishing_marks: petInput.Distinguishing_marks, Known_allergies: petInput.Known_allergies, Price: petInput.Price}

	petPrivateDetailsAsBytes, err := json.Marshal(petPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionPetPrivateDetails", petInput.Key, petPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(petAsBytes)
}

func (s *SmartContract) updatePrivateData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	type petTransientInput struct {
		Owner string `json:"Owner"`
		Date_of_birth string `json:"Date_of_birth"`
		Sex string `json:"Sex"`
		Colour string `json:"Colour"`
		Distinguishing_marks string `json:"Distinguishing_marks"`
		Known_allergies string `json:"Known_allergies"`
		Price string `json:"Price"` 
		Key   string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	petDataAsBytes, ok := transMap["pet"]
	if !ok {
		return shim.Error("pet must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(petDataAsBytes))

	if len(petDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")
	if err != nil {
        return shim.Error(err.Error())
    }
	mspid, err := cid.GetMSPID(APIstub)
	err = verifyClientOrgMatchesPeerOrg(mspid)
    if err != nil {
        return shim.Error(err.Error())
    }
	var petInput petTransientInput
	err = json.Unmarshal(petDataAsBytes, &petInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(petDataAsBytes) + "Error is : " + err.Error())
	}

	petPrivateDetails := &petPrivateDetails{Owner: petInput.Owner, Date_of_birth: petInput.Date_of_birth, Sex: petInput.Sex, Colour: petInput.Colour, Distinguishing_marks: petInput.Distinguishing_marks, Known_allergies: petInput.Known_allergies, Price: petInput.Price}

	petPrivateDetailsAsBytes, err := json.Marshal(petPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionPetPrivateDetails", petInput.Key, petPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(petPrivateDetailsAsBytes)

}
func verifyClientOrgMatchesPeerOrg(clientOrgID string) error {
	peerOrgID, err := shim.GetMSPID()
	if err != nil {
		return fmt.Errorf("failed getting peer's orgID: %v", err)
	}

	if clientOrgID != peerOrgID {
		return fmt.Errorf("client from org %s is not authorized to read or write private data from an org %s peer",
			clientOrgID,
			peerOrgID,
		)
	}

	return nil
}
func (s *SmartContract) createPet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var pet = Pet{Name: args[1], Species: args[2], Breed: args[3], Owner: args[4]}

	petAsBytes, _ := json.Marshal(pet)
	APIstub.PutState(args[0], petAsBytes)
	indexName := "owner~key"
	colorNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{pet.Owner, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(colorNameIndexKey, value)

	return shim.Success(petAsBytes)
}

func (S *SmartContract) queryPetsByOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	owner := args[0]

	ownerAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("owner~key", []string{owner})
	if err != nil {
		return shim.Error(err.Error())
	}

	defer ownerAndIdResultIterator.Close()

	var i int
	var id string

	var pets []byte
	bArrayMemberAlreadyWritten := false

	pets = append([]byte("["))

	for i = 0; ownerAndIdResultIterator.HasNext(); i++ {
		responseRange, err := ownerAndIdResultIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}

		id = compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)

		if bArrayMemberAlreadyWritten == true {
			newBytes := append([]byte(","), assetAsBytes...)
			pets = append(pets, newBytes...)

		} else {
			// newBytes := append([]byte(","), petsAsBytes...)
			pets = append(pets, assetAsBytes...)
		}

		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
		bArrayMemberAlreadyWritten = true

	}

	pets = append(pets, []byte("]")...)

	return shim.Success(pets)
}

func (s *SmartContract) queryAllPets(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "PET0"
	endKey := "PET999"

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

	fmt.Printf("- queryAllPets:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) restictedMethod(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// get an ID for the client which is guaranteed to be unique within the MSP
	//id, err := cid.GetID(APIstub) -

	// get the MSP ID of the client's identity
	//mspid, err := cid.GetMSPID(APIstub) -

	// get the value of the attribute
	//val, ok, err := cid.GetAttributeValue(APIstub, "attr1") -

	// get the X509 certificate of the client, or nil if the client's identity was not based on an X509 certificate
	//cert, err := cid.GetX509Certificate(APIstub) -

	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err != nil {
		// There was an error trying to retrieve the attribute
		shim.Error("Error while retriving attributes")
	}
	if !ok {
		// The client identity does not possess the attribute
		shim.Error("Client identity does not posses the attribute")
	}
	// Do something with the value of 'val'
	if val != "admin" {
		fmt.Println("Attribute role: " + val)
		return shim.Error("Only user with role as ADMIN have access this method!")
	} else {
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		petAsBytes, _ := APIstub.GetState(args[0])
		return shim.Success(petAsBytes)
	}

}

func (s *SmartContract) changePetOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	petAsBytes, _ := APIstub.GetState(args[0])
	pet := Pet{}

	json.Unmarshal(petAsBytes, &pet)
	pet.Owner = args[1]

	petAsBytes, _ = json.Marshal(pet)
	APIstub.PutState(args[0], petAsBytes)

	return shim.Success(petAsBytes)
}

func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	petName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(petName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) createPrivatePetImplicitForDiggipet(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}
	var pet = Pet{Name: args[1], Species: args[2], Breed: args[3], Owner: args[4]}

	petAsBytes, _ := json.Marshal(pet)
	// APIstub.PutState(args[0], petAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_DiggipetMSP", args[0], petAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(petAsBytes)
}

func (s *SmartContract) createPrivatePetImplicitForVetOrg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}

	var pet = Pet{Name: args[1], Species: args[2], Breed: args[3], Owner: args[4]}

	petAsBytes, _ := json.Marshal(pet)
	APIstub.PutState(args[0], petAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_VetorgMSP", args[0], petAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(petAsBytes)
}
func (s *SmartContract) createPrivatePetImplicitForBreederOrg(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect arguments. Expecting 5 arguments")
	}

	var pet = Pet{Name: args[1], Species: args[2], Breed: args[3], Owner: args[4]}

	petAsBytes, _ := json.Marshal(pet)
	APIstub.PutState(args[0], petAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_BreedersorgMSP", args[0], petAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(petAsBytes)
}
func (s *SmartContract) queryPrivateDataHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	petAsBytes, _ := APIstub.GetPrivateDataHash(args[0], args[1])
	return shim.Success(petAsBytes)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
