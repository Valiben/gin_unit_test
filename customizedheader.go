package gin_unit_test

// customized request headers
type CustomizedHeader struct {
	Key     string // header's name
	Value   string // header's value
	IsValid bool   // whether add the header to the request
}

// set the header
func SetHeader(key, value string, isValid bool) {
	if myHeaders == nil {
		return
	}

	for i, data := range myHeaders {
		if data.Key == key {
			myHeaders[i].IsValid = isValid
			myHeaders[i].Value = value
			return
		}
	}
}

// get header
func GetHeader(key string) *CustomizedHeader {
	if myHeaders == nil {
		return nil
	}

	for _, data := range myHeaders {
		if data.Key == key {
			return &data
		}
	}

	return nil
}

// add header
func AddHeader(key, value string, isValid bool) {
	if myHeaders == nil {
		myHeaders = make([]CustomizedHeader, 3)
	}

	header := CustomizedHeader{
		Key:     key,
		Value:   value,
		IsValid: isValid,
	}

	myHeaders = append(myHeaders, header)
}

// delete header
func DeleteHeader(key string) {
	if myHeaders == nil {
		return
	}

	for i, data := range myHeaders {
		if data.Key == key {
			myHeaders = append(myHeaders[:i], myHeaders[i+1:]...)
			return
		}
	}
}
