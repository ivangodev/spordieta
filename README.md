## HOWTO 

Start Core and Storage microservices via docker-compose:

```
$ docker-compose up --build
...
$ docker-compose down --volume
```

Example of API usage:

```
     Method: POST //create user
        URL: http://localhost:8080/user/0
     Response status:
        200 OK
     Response body:
        
        
        
     Method: POST //create bet
        URL: http://localhost:8080/user/0/bet
        Body: 
        	{
        		"currWeight":  80,
        		"goalWeight":  75,
        		"money":	   100,
        		"durationSec": 60
        	}
        	
     Response status:
        200 OK
     Response body:
        {"bet_id": "0"}
        
        
     Method: GET //get opened bet
        URL: http://localhost:8080/user/0/bet/opened
     Response status:
        200 OK
     Response body:
        {"bet_id": "0"}
        
        
     Method: GET //get bet info
        URL: http://localhost:8080/user/0/bet/0
     Response status:
        200 OK
     Response body:
	{
	  "Cond": {
		"currWeight": 80,
		"goalWeight": 75,
		"money": 100,
		"durationSec": 60,
		"deadline": "2022-01-31T19:52:17.678538638+03:00"
	  },
	  "Status": {
		"opened": true,
		"state": 0,
		"uploaded": false,
		"adminComment": ""
	  }
	}
        
        
     Method: PUT //upload win proof
        URL: http://localhost:8081/user/0/bet/0/winproof
        Body: I lost weight!
     Response status:
        200 OK
     Response body:
        
        
        
     Method: GET //get links to review
        URL: http://localhost:8080/toreview
     Response status:
        200 OK
     Response body:
        ["/user/0/bet/0/winproof",]


     Method: GET //get link where the data resides
        URL: http://localhost:8081/user/0/bet/0/winproof
     Response status:
        200 OK
     Response body:
        /user/0/bet/0/winproof/data


     Method: GET //get actual data
        URL: http://localhost:8081/user/0/bet/0/winproof/data
     Response status:
        200 OK
     Response body:
        I lost weight!
        
        
     Method: PATCH //admin didn't like the proof
        URL: http://localhost:8080/user/0/bet/0
        Body: 
        	{
        		"opened": true,
        		"state": 1,
        		"uploaded": false,
        		"adminComment": "You failed to lose weight. You must pay."
        	}
        	
     Response status:
        200 OK
     Response body:
        
        
        
     Method: PUT //upload pay proof
        URL: http://localhost:8081/user/0/bet/0/payproof
        Body: I payed!
     Response status:
        200 OK
     Response body:
        
        
        
     Method: PATCH //admin liked the proof and closed the bet
				   //(as with win proof, he had to call /toreview, get data URL
				   //and so on, but here it's omitted since it's literally same)
        URL: http://localhost:8080/user/0/bet/0
        Body: 
        	{
        		"opened": false
        	}
        	
     Response status:
        200 OK
     Response body:
        
        
```
