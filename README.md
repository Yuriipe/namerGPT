The following app is used to  translate and create product descriptions in different languages based on available GPT models using API connection.
OpenAI token is required.

"cfg.json" options explained:

	"model":"gpt-3.5-turbo-1106" - default gpt processor model
 
	"language":"english" - default language
 
	"sourceFile":"db/example.csv" - source file path and name (if the file is changed, change the filename only)
 
	"tokenLimName":30 - response product name token limit
 
	"tokenLimDescript":200 - response product description token limit
 
	"temp":0.2 - creativity level (0.1 - the least creative, 1 - the most creative)
 
	"APIKey":"YOUR_OPENAI_API_TOKEN" - change the desription to your token
 
	"APIEndpoint":"https://api.openai.com/v1/chat/completions" - default chat endpoint

"example.csv" shows the requred input format (only formated csv files are supported) - fill according to your preferences and needs

some GPT models estimated processing time (v - values, t - token limit for single value, temp - response temperature)

	// "gpt-3.5-turbo-0301" - 3m 25s : 90 v : 30 t : 0.2 temp
	// "gpt-3.5-turbo-16k-0613" - 5m 14s : 90 v : 30 t : 0.2 temp
	// "gpt-3.5-turbo-1106" - 3m 25s : 90 v : 30 t : 0.1 temp
