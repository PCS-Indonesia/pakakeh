package pipeline

// ResponseImplementor define each function response inside the pipeline
type ResponsesImplementor interface {
	Get() []any
	Add(any) ResponsesImplementor
}

// PipelineResponse struct that stored all pipeline response
type PipelineResponses struct {
	resps []any
}

// Get returns all responses stored in the pipeline response
func (pr PipelineResponses) Get() []any {
	return pr.resps // retunr slice of response
}

// Add appends a response to the pipeline response collection.
// If the input is another PipelineResponse, its responses are flattened and added individually.
// Otherwise, the input is added as-is. Returns the updated PipelineResponse.
func (pr PipelineResponses) Add(anyResp any) ResponsesImplementor {
	if res, ok := anyResp.(PipelineResponses); ok {
		pr.resps = append(pr.resps, res.resps...)
	} else {
		pr.resps = append(pr.resps, anyResp)
	}

	return pr
}
