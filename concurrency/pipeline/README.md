# Pakakeh Pipeline

Pipeline is a `Go` function to organize our code to be more readable and clean in form of functions series or chains. We can use pipeline to serialize our logic into several functions with Single Responsibilty from SOLID principle.


## How To Use

This package can be used by adding the following import `pipeline` statement to your `.go` files.

```go
import "github.com/PCS-Indonesia/pakakeh/concurrency/pipeline" 
```

## Example Usage
#### Example 1 - Simple pipeline
```go
package main

import (
	"errors"
	"strings"

	"github.com/PCS-Indonesia/pakakeh/concurrency/pipeline"
)

func main() {
	e := pipeline.Pipeline(
		isUserEmailExists,
		validateUserEmail,
		insertNewUser,
	)

	_, err := e(UserInput{
		Email: "elvizar.kh@pcsindonesia.co.id",
	})
	if err != nil {
		panic(err)
	}
}

// mock Database as map for quick example
var DB = make(map[string]any)

type UserInput struct {
	Email    string
	Password string
}

func isUserEmailExists(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	_, exists := DB[args.Email]
	if exists {
		return nil, errors.New("email already exists")
	}
	return nil, nil
}

func validateUserEmail(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	if !strings.Contains(args.Email, "@") {
		return nil, errors.New("incorrect email address")
	}
	return nil, nil
}

func insertNewUser(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	DB[args.Email] = args
	return nil, nil
}
```
The example above demonstrates a simple user registration flow using Pipeline. It shows how to:

1. Chain multiple functions together using `pipeline.Pipeline()`
2. Pass data between functions using a shared `UserInput` struct
3. Handle errors at each step
4. Access a mock database

The pipeline consists of 3 functions:

- `isUserEmailExists`: Checks if email already exists in DB
- `validateUserEmail`: Validates email format 
- `insertNewUser`: Inserts the new user into DB

Each function follows Single Responsibility Principle and can be tested independently. The pipeline executes these functions sequentially, jumping to each function to get deeper context, stopping if any function returns an error (See [Example 2](#example-2---using)).

#### Example 2 - using `Responses` 
As a `Pipe` function defined, we can use a previous response as an input for the next function (This is optional if you want to use it)
Here's how the example works:

```go
package main

import (
    "strings"
	"errors"

	"github.com/PCS-Indonesia/pakakeh/concurrency/pipeline"
)

func main() {
	e := pipeline.Pipeline(
		getBlacklistUsers,
		isBlacklistUser,
		... // and so on
	)

	_, err := e(UserInput{
		Email: "elvizar.blacklist@ismail.com",
	})
	if err != nil {
		panic(err)
	}
}

var (
	DB = make(map[string]any) // Mock DB as a mapß
)

type UserInput struct {
	Email    string
	Password string
}

func getBlacklistUsers(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	// Mock get all blacklist users from DB
    return map[string]any{
		"gammarizi@pcsindonesia.co.id": struct{}{},
		"elvizar.blacklist@ismail.com":  struct{}{},
	}, nil
}

func isBlacklistUser(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	blacklistUsers := pipeline.Get[map[string]any](responses)
    
	_, isBlacklist := blacklistUsers[args.Email]
	if isBlacklist {
		return nil, errors.New("this email is blacklisted")
	}

	return nil, nil
}
```

From code above, we can see this scenario is to validate the incoming users whether is blackisted or not, so at the beginning we can get the blacklist users. Then, we utilize the `responses` on the next function to validate the users.
1. `getBlacklistUsers` returns a map of blacklisted email addresses
2. `isBlacklistUser` uses `pipeline.Get[T]` to retrieve the blacklist map from previous response
3. Checks if current user's email exists in blacklist
4. Returns error if email is blacklisted

This demonstrates how to:

- Pass data between pipeline functions using `Responses`
- Type-safely access previous responses with `pipeline.Get[T]`
- Build more complex validation flows
- Reuse response data across multiple functions


#### Example 3 - Concurrency Pipeline with `PipelineGo`
The pipeline also support concurrency by simply use `PipelineGo` instead. 

_PS : This abstract concurrency under the hood, so we don't need to write Go routine manually_

```go
package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PCS-Indonesia/pakakeh/concurrency/pipeline"
)

func main() {
	e := pipeline.Pipeline(
		getBlacklistUsers,
		isBlacklistUser,
		validateUserEmail,
		pipeline.PipelineGo( // this will be executed concurrently
			insertUser,
			sendNotification,
		),
	)

	_, err := e(UserInput{
		Email: "ndasmu_wowo@hotbabe.mail",
	})
	if err != nil {
		panic(err)
	}
}

func insertUser(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	DB[args.Email] = args
	return nil, nil
}

func sendNotification(args UserInput, responses pipeline.ResponsesImplementor) (response any, err error) {
	fmt.Println("send notification")
	return nil, nil
}
```



## Test Coverage 
`ok  	github.com/PCS-Indonesia/pakakeh/concurrency/pipeline	2.418s	coverage: 80.2% of statements` see **_pipeline_test.go_** file and execute it using `go test -v` command



## PROS & CONS
Pros :
- You can compose all your functions into 1 single logic. This makes code more readable and clean in forms of function series
- **Single Responsibility**: Each function focuses on one specific task, following SOLID principles
- **Concurrency Support**: Built-in support for concurrent execution using PipelineGo without manual goroutine management

Cons :
- Need to update your function structure to follow pipeline format (if it applies to your existing function)
- **State Management**: Sharing state between pipeline stages requires passing through responses
- currently , it is harder to debug issues in concurrent pipeline stages
- **Error Recovery**: No built-in way to recover from errors and continue pipeline execution