# httpAsyncStress
Http Async Stress Test written in Go Lang.


# The json to use in configure.
```
{
	"body": "<BODY WITH MARKS THAT CAN BE REPLACED>",
	"replaces": [
		{ 
				"values": <[ array with numbers ]>,
				"sort": "rand",
				"type": "array",
				"mark": "##CITY##"
		},
		{
				"value": 15780000,
				"format": "2006-01-02",
				"type": "date",
				"mark": "##INDATE##"
		},
		{
				"value": 16384800,
				"format": "2006-01-02",
				"type": "date",
				"mark": "##OUTDATE##"
		}
	],
	"headers": [
		{ 
			"name" : "<HEADER NAME>",
			"value" : "<HEADER VALUE>"
		}
	],
	"query": [
		{ 
			"name" : "<QUERY NAME>",
			"values" : ["ARRAY WITH QUERY VALUES"]
		}
	],
	"url": "<URL>",
	"type": "<METHOD>",
	"grep": "<Text to FIND>"
} 
```
