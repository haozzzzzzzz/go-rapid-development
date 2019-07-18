package test_comment

/**
@api_doc_start
{
	"http_method": "GET",
	"relative_paths": ["/hello_world"],
	"query_data": {
		"name": "姓名|string|required"
	},
	"post_data": {
		"age": "年龄|string|required",
		"nicknames": ["string"],
		"__nicknames": "别名|array",
		"contacts": [
			{
				"location": "地址|string|required"
			}
		],
		"__contacts": "联系方式|array|required",
		"info": {
			"has_dog": "是否养了狗|array"
		},
		"__info": "其他信息|object|required"
	},
	"resp_data": {
	    "a": "这个是a|int",
	    "b": "这个是b|int",
	    "c": {
			"d": "这个是d|string"
		},
		"__c": "这个是c|object",
		"f": [
			"string"
		],
		"__f": "这个是f|object|required",
		"g": [
			{
				"h": "这个是h|string|required"
			}
		],
		"__g": "这个是g|array|required"
	}
}
@api_doc_end
*/
func something() {

}
