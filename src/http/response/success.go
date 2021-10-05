package response

import "net/http"

// SucessMsg represents the structure of errors as they occur in the application.
type SucessMsg struct {
	OK          bool        `json:"ok"`
	Status      int         `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

// Success TODO
func Success(msg SucessMsg, data interface{}) (int, SucessMsg) {
	msg.Data = data
	return msg.Status, msg
}

// SuccessCreate TODO
var SuccessCreate SucessMsg = SucessMsg{
	OK:          true,
	Status:      http.StatusCreated,
	Description: "Your resource was successfully created.",
}

// SuccessDelete TODO
var SuccessDelete SucessMsg = SucessMsg{
	OK:          true,
	Status:      http.StatusAccepted,
	Description: "Your resource was successfully deleted.",
}

// SuccessLogin TODO
var SuccessLogin SucessMsg = SucessMsg{
	OK:          true,
	Status:      http.StatusOK,
	Description: "Authentication successful",
}

// SuccessLogout TODO
var SuccessLogout SucessMsg = SucessMsg{
	OK:          true,
	Status:      http.StatusOK,
	Description: "Log out successful",
}

// SuccessRead TODO
var SuccessRead SucessMsg = SucessMsg{
	OK:          true,
	Status:      http.StatusOK,
	Description: "Your resource was successfully read.",
}

// SuccessUpdate TODO
var SuccessUpdate SucessMsg = SucessMsg{
	OK:          true,
	Status:      http.StatusAccepted,
	Description: "Your resource was successfully updated.",
}
