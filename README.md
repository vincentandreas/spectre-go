# Spectre-Go
Unofficial repository, re-create  [Spectre app](https://spectre.app/), based on Go language.

### Important Notes:
This repo just for my own learning purpose, **don't use this repo for generating your real password**. 
Please use the official one in https://spectre.app/ . I don't responsible for anything bad could be happened by using this repo.

## Introduction
### What is Spectre?
Spectre brings a stateless cryptographic identity algorithm to you across multiple platforms.

### How does it work?
With just username and spectre secret (master password), you don't have to remember all site / app passwords.

Better try to test the app in https://spectre.app

![Spectre web](docs/spectre-demo.gif)


### Original Features
There's some features in spectre:

- Choose what you want to generate. It could be:
  - password
  - site name (username)
  - security answer
- Choose the key type: 
  - Short, Medium, Long, Max length password
  - Pin
  - Phrase (for security answer)
  - Name (for username)
- Set the keyCounter, to get different result in the same site. It's really useful for updating your password regularly.
- All passwords are 100% generated locally.


## Spectre-go 

### What's inside this repo?

The repo contains main spectre algorithm to generate password. This repo based on spectre newest algorithm version V3(2015:01)
You'll get the same generated password from my repo and the official one.

I made this using go language, to generate the password, you need to send the request via REST API.

### API List

#### Get health status [GET /api/health ]

Response example:
```
{
    "result": "OK"
}
```


#### Generate password [POST /api/generatePassword ]

Header:
```
Utilize-Cache : [true | false]
```

Body:

```
username    = your username
password    = your master password (spectre secret)
site        = site that you want to generate
keyCounter  = counter for generating different password from the same site, range from 1 - 4294967295 ((2^32) - 1)
keyPurpose  = [password | loginName | answer]
keyType     = [med | long | max | short | basic | pin | name | phrase]
```



Request example:
```json
{
    "username":"a", 
    "password":"a", 
    "site":"twitter.com", 
    "keyCounter":1,
    "keyPurpose":"password",
    "keyType":"med"
}
```

Response success example :

```json
{
    "responseCode": "00", 
    "responseMessage": "Success",
    "result": "RevXep5+"
}
```

Response failed example:
```json
{
    "responseCode": "01",
    "responseMessage": "Request not valid",
    "result": ""
}
```

Result example:

![Go version - long](docs/long-go.png)

![Go version - med](docs/med-go.png)

The official one:

![Ofc version - long](docs/long-web.png)

![Ofc version - med](docs/med-web.png)

### How I recreate this app?
- There's many option to see the algorithm implmentation in android, iOS, web, CLI, etc. I choose the implementation by using web. (Since it's easier to debug)
- I try to understand the flow of how the spectre generate the password, and here's the main flow of spectre web:
![Flow](docs/main-spectre-flow.png)
- Since I understand python, and it's really simple for experiment with that, I try to do the Proof Of Concept, by implement that. you can check the POC in `docs/spectre_poc.py`
- When the python code already work properly, I try to recreate that using Go language. 

### How the algorithms work
![Flow](docs/simple-flow.png)

It's consist of 3 main function:
- newUserKey : generate userKeyData using [Scrypt](https://en.wikipedia.org/wiki/Scrypt)
- newSiteKey : generate siteKey, based on userKeyData, and other params. Use [HMAC-SHA256](https://en.wikipedia.org/wiki/HMAC)
- newSiteResult : generate password based on siteKey and choosen templates.

The value of siteKey will be modulated to get character based on the position 

```go
func NewSiteResult(params models.GenSiteParam) string {
	userKey := newUserKey(params.Username, params.Password, params.KeyPurpose)
	siteKey := newSiteKey(userKey, params.Site, params.KeyCounter, params.KeyPurpose, "")
	
	resTemplates := templates[params.KeyType] // e.g. "CvcnoCvc"
	resTemplate := resTemplates[int(siteKey[0])%len(resTemplates)]
	
	var passRes strings.Builder
	for i := 0; i < len(resTemplate); i++ { // e.g. "C" -> "BCDFGH..." (Consonant template)
		currChar := characters[string(resTemplate[i])] // 
		idx := int(siteKey[i+1]) % len(currChar)
		passRes.WriteRune([]rune(currChar)[idx])
	}

	return passRes.String()
}
```
### Different with the official one:
#### Using REST API
The original one have the web view, but I haven’t implemented that, so I choose to make that as a REST API. It’s not as safe as generating locally, for sure.

#### Caching using Redis
I'm doing experiment, by using redis for saving generated password temporarly, 
it's boost the speed when doing the generating for the same password, again, and again. Really suitable for people with Short-term memory loss ✌️✌️. 

Benchmark result (with same payload, without cache)
```
Benchmark_intTest_genPassword_withOutCache-8           5         254846040 ns/op
```
Benchmark result (with same payload, using redis)
```
Benchmark_intTest_genPassword_withCache-8              8         125318125 ns/op
```

After doing this, I was wondering, does the save cache feature really affect the performance, when all of the payload is different? So I'm experimenting, by using random UUID as the payload, to create different payload and here's the benchmark result:
```
Benchmark_intTest_genPassword_withCache-8              4         264977950 ns/op
Benchmark_intTest_genPassword_withOutCache-8           5         254846040 ns/op

```


### Future work:
- Implement HTTPS, you don't want your plain password intercepted by someone (at least your own Burpsuite) right? 
- Create docker image, so you can run the container, in CLI / REST mode
- Adding more feature, and enhance the code
- Creating a nice web view. Nice view definition for me is: you can use the web easily. You don't have to be worry with the fancy animation effect, color. I won't create that. 

### Credit:
- https://gitlab.com/spectre.app
- https://tutorialedge.net/golang/go-encrypt-decrypt-aes-tutorial/
- https://medium.com/insiderengineering/integration-test-in-golang-899412b7e1bf