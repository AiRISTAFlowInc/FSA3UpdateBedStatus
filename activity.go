package UpdateBedStatus

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/project-flogo/core/activity"
)

func init() {
	_ = activity.Register(&Activity{}) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Input{}, &Output{})

// Activity is an sample Activity that can be used as a base to create a custom activity
type Activity struct {
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	updateStatus, bedStatus := UpdateBedStatus(input.IP, input.CustomerId, input.Username, input.Password, input.MAC)

	output := &Output{Status: updateStatus, BedStatus: bedStatus}

	// fmt.Println("Output: ", output.Status)
	 ctx.Logger().Info("Output: ", output)

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func UpdateBedStatus(IP string, customerId string, uname string, pword string, MAC string) (bool, string) {
	itemID := GetByMACAddress(IP, customerId, uname, pword, MAC)
	staff := GetByStaffId(IP, customerId, uname, pword, itemID)
	status, bedStatus := nextBedStatus(IP, customerId, uname, pword, itemID, staff)
	return status, bedStatus
}

func GetByMACAddress(IP string, customerId string, uname string, pword string, MAC string) string {
	// Create an HTTP client
	client := &http.Client{}
	cleanMAC := url.QueryEscape(MAC)

	// Create the request
	url := "http://" + IP + "/XpertRestApi/api/Device/GetByMacAddress?MacAddress="+cleanMAC+"&CustomerId="+customerId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	// Add basic authentication to the request header
	auth := uname + ":" + pword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	// Unmarshal the config JSON into an object
	var device Device
	errUnmarshal := json.Unmarshal(body, &device)
	if errUnmarshal != nil {
	 	fmt.Println(errUnmarshal)
		return ""
	}

	return strconv.Itoa(device.ItemID)
}

func GetByStaffId(IP string, customerId string, uname string, pword string, staffId string) Staff {
	// Staff to return 
	var staff Staff

	// Create an HTTP client
	client := &http.Client{}

	// Create the request
	url := "http://" + IP + "/XpertRestApi/api/Staff/GetByStaffId?StaffId="+staffId+"&CustomerId="+customerId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return staff
	}

	// Add basic authentication to the request header
	auth := uname + ":" + pword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return staff
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return staff
	}

	// Unmarshal the config JSON into an object
	errUnmarshal := json.Unmarshal(body, &staff)
	if errUnmarshal != nil {
	 	fmt.Println(errUnmarshal)
		return staff
	}

	return staff
}

func nextBedStatus(IP string, customerId string, uname string, pword string, staffId string, staff Staff) (bool, string) {
	assocItemId := strconv.Itoa(staff.AssocItemID)
	status:= false
	bedStatusOutput := "ERROR"

	switch bedStatus := staff.BedStatus; bedStatus{
	case "ASSIGNED":
		println("assigned")
		status = changeItemAssociation(IP, customerId, uname, pword, staffId, assocItemId, "DISCHARGING")
		bedStatusOutput = "DISCHARGING"
	case "DISCHARGING":
		println("discharging")
		endItemAssociation(IP, customerId, uname, pword, staffId, assocItemId)
		status = createItemAssociation(IP, customerId, uname, pword, staffId, assocItemId, "CLEANING")
		bedStatusOutput = "CLEANING"
	case "CLEANING":
		println("cleaning")
		status = endItemAssociation(IP, customerId, uname, pword, staffId, assocItemId)
		bedStatusOutput = "AVAILABLE"
	}
	println("stauts: ", bedStatusOutput)
	
	return status, bedStatusOutput
}

func changeItemAssociation(IP string, customerId string, uname string, pword string, itemId string, assocItemId string, associationType string) bool{
	// Create an HTTP client
	client := &http.Client{}
	// Create the request
	url := "http://" + IP + "/XpertRestApi/api/Staff/ChangeItemAssociation?CustomerId="+customerId+"&AssociationType="+associationType+"&ItemID="+itemId+"&AssociatedItemID="+assocItemId
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Add basic authentication to the request header
	auth := uname + ":" + pword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	// Unmarshal the config JSON into an object
	var response Response
	errUnmarshal := json.Unmarshal(body, &response)
	if errUnmarshal != nil {
	 	fmt.Println(errUnmarshal)
		return false
	}
	if(response.ErrorMessage == ""){// If successful return true 
		return true
	}

	return false
}

func endItemAssociation(IP string, customerId string, uname string, pword string, itemId string, assocItemId string) bool{
	// Create an HTTP client
	client := &http.Client{}
	// Create the request
	url := "http://" + IP + "/XpertRestApi/api/Staff/EndItemAssociation?CustomerId="+customerId+"&ItemID="+itemId+"&AssociatedItemID="+assocItemId
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Add basic authentication to the request header
	auth := uname + ":" + pword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	// Unmarshal the config JSON into an object
	var response Response
	errUnmarshal := json.Unmarshal(body, &response)
	if errUnmarshal != nil {
	 	fmt.Println(errUnmarshal)
		return false
	}
	if(response.ErrorMessage == ""){// If successful return true 
		return true
	}else{println("EndItemAssociationError: ", response.ErrorMessage)}

	return false
}

func createItemAssociation(IP string, customerId string, uname string, pword string, itemId string, assocItemId string, associationType string) bool{
	// Create an HTTP client
	client := &http.Client{}
	// Create the request
	url := "http://" + IP + "/XpertRestApi/api/Staff/CreateItemAssociation?CustomerId="+customerId+"&AssociationType="+associationType+"&ItemID="+itemId+"&AssociatedItemID="+assocItemId
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	// Add basic authentication to the request header
	auth := uname + ":" + pword
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", basicAuth)
	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return false
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return false
	}

	// Unmarshal the config JSON into an object
	var response Response
	errUnmarshal := json.Unmarshal(body, &response)
	if errUnmarshal != nil {
	 	fmt.Println(errUnmarshal)
		return false
	}
	if(response.ErrorMessage == ""){// If successful return true 
		return true
	}

	return false
}