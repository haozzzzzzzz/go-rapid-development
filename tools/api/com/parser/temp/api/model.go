package api

type Request struct {
	ReqKey1 string `json:"req_key_1" form:"req_key_1"`
}

type Field struct {
	//FieldKey1 string      `json:"field_key_1" form:"field_key_1"`
	FieldKey2 FieldChild1 `json:"field_key_2" form:"field_key_2"`
}

type FieldChild1 struct {
	//FieldChild1Key1 string      `json:"field_child_1_key_1" form:"field_child_1_key_1"`
	FiledChild1Key2 FiledChild2 `json:"filed_child_1_key_2" form:"filed_child_1_key_2"`
}

type FiledChild2 struct {
	FiledChild2Key string `json:"filed_child_2_key" form:"filed_child_2_key"`
}
