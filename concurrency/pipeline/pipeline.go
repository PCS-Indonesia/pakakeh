package pipeline

import (
	"context"
	"fmt"
)

// Func is a series of pipe functions arguments contract
// when using pipe the function callback must comply with this contract
type Func[T any] func(args T, responses ResponsesImplementor) (response any, err error)

/*
------------------------------ Without Context and Concurrency --------------------------------
*/

// Pipeline to start compose functions with Pipe. Use this when initiate pipe at the beginning
func Pipeline[T any](seqFuncs ...Func[T]) func(args T) (responses ResponsesImplementor, err error) {
	return func(args T) (responses ResponsesImplementor, err error) {
		execFunc := Pipe(seqFuncs...)
		fmt.Println("AWALNYA", responses)
		resp, errFunc := execFunc(args, nil)

		r, _ := resp.(ResponsesImplementor)

		fmt.Println("HASILNYA DONG", r)
		return r, errFunc
	}
}

// Pipe composes a series of functions into a single processing unit.
// It takes a variadic number of Func[T] functions and returns a new Func[T] function
// that executes them sequentially with response handling.
//
// Parameters:
//   - seqFuncs: Variadic Func[T] functions to be executed in sequence
//
// Returns:
//   - Func[T]: A composed function that:
//     1. Converts regular functions to context-aware functions
//     2. Uses PipeCtx to compose the context-aware functions
//     3. Executes the composed function with a background context
//
// The function acts as a bridge between non-context and context-aware pipelines by:
// 1. Converting each Func[T] to FuncCtx[T] while preserving the original function
// 2. Using PipeCtx for the actual composition and execution
// 3. Providing a background context automatically for the execution
func Pipe[T any](seqFuncs ...Func[T]) Func[T] {
	var funcCtx []FuncCtx[T]

	for _, fn := range seqFuncs {
		_fi := fn

		funcCtx = append(funcCtx, func(ctx context.Context, args T, responses ResponsesImplementor) (response any, err error) {
			return _fi(args, responses)
		})
	}

	fCtx := PipeCtx(funcCtx...)

	return func(args T, responses ResponsesImplementor) (response any, err error) {
		return fCtx(context.Background(), args, responses)
	}
}

/*
------------------------------ With Context --------------------------------
*/

// FuncCtx is a type of function series of pipe functions arguments with context
// when using pipe the function callback must comply with this contract
type FuncCtx[T any] func(ctx context.Context, args T, responses ResponsesImplementor) (response any, err error)

func PipelineWithCtx[T any](seqFuncs ...FuncCtx[T]) func(ctx context.Context, args T) (responses ResponsesImplementor, err error) {
	return func(ctx context.Context, args T) (responses ResponsesImplementor, err error) {
		execFunc := PipeCtx(seqFuncs...)
		resp, err := execFunc(ctx, args, nil)

		r, _ := resp.(ResponsesImplementor)

		//fmt.Println("EH PIPELINE WITH CTX", r)
		return r, err
	}
}

// PipeCtx composes a series of context-aware functions into a single processing unit.
// It takes a variadic number of FuncCtx[T] functions and returns a new FuncCtx[T] function
// that executes them sequentially with context and response handling. This functions handle two cases of responses
//
// Parameters:
//   - seqFuncs: Variadic FuncCtx[T] functions to be executed in sequence
//
// Returns:
//   - FuncCtx[T]: A composed function that:
//     1. Initializes response handling
//     2. Executes each function in sequence with context
//     3. Aggregates responses appropriately
//     4. Returns either new or accumulated responses based on pipe state
//
// The function handles two cases:
// 1. When called directly - accumulates all responses in the original responses object
// 2. When called from another pipe - creates new responses object to avoid modifying original
func PipeCtx[T any](seqFuncs ...FuncCtx[T]) FuncCtx[T] {
	return func(ctx context.Context, args T, responses ResponsesImplementor) (response any, err error) {
		_, fromPipe := responses.(PipelineResponses)

		if responses == nil {
			responses = PipelineResponses{} // init Response struct
		}

		var newResponses ResponsesImplementor = PipelineResponses{}
		for _, fn := range seqFuncs {
			if response, err = fn(ctx, args, responses); err != nil {
				return nil, err
			}

			// Add responses
			responses = responses.Add(response)

			// if response is get from Pipe
			if fromPipe {
				newResponses = newResponses.Add(response)
			}
		}

		if fromPipe {
			return newResponses, nil
		}

		return responses, nil
	}
}

/*
------------------------------ With Concurrency --------------------------------
*/
// PipeGo converts a sequence of Func[T] functions to a concurrent pipeline.
// It takes a variadic number of Func[T] functions and returns a single Func[T] that executes them concurrently.
//
// Parameters:
//   - seqFuncs: Variadic Func[T] functions to be executed concurrently
//
// Returns:
//   - Func[T]: A composed function that:
//     1. Converts each Func to a FuncCtx by wrapping it with context handling
//     2. Creates a concurrent pipeline using PipeCtxGo
//     3. Returns a function that executes the pipeline with a background context
func PipelineGo[T any](seqFuncs ...Func[T]) Func[T] {
	var funcCtx []FuncCtx[T]

	// Convert each Func to FuncCtx by wrapping with context
	for _, fn := range seqFuncs {
		_fi := fn

		funcCtx = append(funcCtx, func(ctx context.Context, args T, responses ResponsesImplementor) (response any, err error) {
			return _fi(args, responses)
		})
	}

	// Create concurrent pipeline with context
	fCtxGo := PipeCtxGo(funcCtx...)

	// Return function that executes pipeline with background context
	return func(args T, responses ResponsesImplementor) (response any, err error) {
		return fCtxGo(context.Background(), args, responses)
	}
}

// PipeCtxGo executes a sequence of functions concurrently with context as first class interface.
// It takes a variadic number of FuncCtx[T] functions and returns a single FuncCtx[T].
// Each function in the sequence is executed in its own goroutine.
func PipeCtxGo[T any](seqFuncs ...FuncCtx[T]) FuncCtx[T] {
	return func(ctx context.Context, args T, responses ResponsesImplementor) (response any, err error) {
		// Check if response is from another pipe to handle nested pipelines
		_, fromPipe := responses.(PipelineResponses)

		// Initialize responses if nil
		if responses == nil {
			responses = PipelineResponses{} // init Response struct
		}

		// Channel to collect results from goroutines
		structChan := make(chan struct {
			index    int
			response any
			err      error
		})

		// Execute each function sequence concurrently
		// Each goroutine executes a function and sends its result through the channel
		for ind, fn := range seqFuncs {
			go func(f FuncCtx[T], i int) {
				response, err = f(ctx, args, responses)

				structChan <- struct {
					index    int
					response any
					err      error
				}{
					index:    i,
					response: response,
					err:      err,
				}
			}(fn, ind)
		}

		// Collect responses from goroutines in order of completion
		mapResponse := make(map[int]any)

		// Wait for all goroutines to complete and collect their responses
		for i := 0; i < len(seqFuncs); i++ {
			res := <-structChan

			// Return early if any function returns an error
			if res.err != nil {
				return response, res.err
			}

			// Store response in map using original function index
			mapResponse[res.index] = res.response
		}

		// Create new responses object for nested pipeline handling
		var newResponses ResponsesImplementor = PipelineResponses{}

		// Process responses in original sequence order
		for j := 0; j < len(mapResponse); j++ {
			// Add response to main responses collection
			responses = responses.Add(mapResponse[j])

			// If from nested pipe, also add to new responses
			if fromPipe {
				newResponses = newResponses.Add(mapResponse[j])
			}
		}

		// fmt.Println("KESINI GA 2", mapResponse, newResponses, responses)

		// Return appropriate responses based on pipeline nesting
		if fromPipe {
			return newResponses, nil
		}

		// in case the response not from Pipeline
		return responses, nil
	}
}
