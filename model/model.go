package model

//RequestModel : Request Model for Post API
type RequestModel struct {
	EmpID string `json:"emp_id" sql:"emp_id"`
	Name  string `json:"name" sql:"name"`
	Age   int    `json:"age" sql:"age"`
	Hobby string `json:"hobby" sql:"hobby"`
}

//ResponseModel : Response Model For Get API
type ResponseModel struct {
	Data []RequestModel `json:"data"`
}
