The following application is utilized to translate and generate product descriptions in various languages, leveraging available GPT models through API connection. 
An OpenAI token is required.

Explanation of options in "cfg.json":

	"model":"gpt-3.5-turbo-1106" - Default gpt processor model
 
	"language":"english" - Default language
 
	"sourceFile":"db/example.csv" - Path and name of the source file (if the file is modified, only change the filename).
 
	"tokenLimName":30 - Token limit for the response product name
 
	"tokenLimDescript":200 - Token limit for the response product description
 
	"temp":0.2 - creativity level (0.1 - the least creative, 1 - the most creative)
 
	"APIKey":"YOUR_OPENAI_API_TOKEN" - Update with your token
 
	"APIEndpoint":"https://api.openai.com/v1/chat/completions" - Default chat endpoint

"example.csv" illustrates the required input format (only formated csv files are supported) - fill it according to your preferences and needs

Estimations for processing time of some GPT models (v - values, t - token limit for single value, temp - response temperature)

	// "gpt-3.5-turbo-0301" - 3m 25s : 90 v : 30 t : 0.2 temp
	// "gpt-3.5-turbo-16k-0613" - 5m 14s : 90 v : 30 t : 0.2 temp
	// "gpt-3.5-turbo-1106" - 3m 25s : 90 v : 30 t : 0.1 temp
