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
		"location": "地址|string|required"
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
