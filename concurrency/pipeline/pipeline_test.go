package pipeline

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	type createUserRequest struct {
		Username string
		Email    string
		Password string
	}

	type args struct {
		args  createUserRequest
		funcs []Func[createUserRequest]
	}
	tests := []struct {
		name string
		args args
		want ResponsesImplementor
	}{
		{
			name: "create user",
			args: args{
				args: createUserRequest{
					Username: "username",
					Email:    "email",
					Password: "hash-password",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						username := Get[string](responses)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						password := Get[string](responses)
						fmt.Println(password)
						return args.Password, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						return nil, nil
					},
				},
			},
			want: PipelineResponses{
				resps: []any{
					"username",
					"email",
					"hash-password",
					nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := Pipeline(tt.args.funcs...)
			got, _ := exec(tt.args.args)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPipeAndPipeGo(t *testing.T) {
	type createUserRequest struct {
		Username string
		Email    string
		Password string
	}

	type args struct {
		args  createUserRequest
		funcs []Func[createUserRequest]
	}
	tests := []struct {
		name string
		args args
		want ResponsesImplementor
	}{
		{
			name: "create user",
			args: args{
				args: createUserRequest{
					Username: "username",
					Email:    "email",
					Password: "fariz-jnck",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						username := Get[string](responses)
						fmt.Println(username)
						return args.Email, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						email := Get[string](responses)
						fmt.Println(email, responses)
						return args.Password, nil
					},
					PipelineGo(
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							time.Sleep(1 * time.Second)
							// fmt.Println("UUYEAA", responses)
							username, _ := Index[string](responses, 0)

							return "pipego-" + username, nil
						},
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							time.Sleep(1 * time.Second)
							fmt.Println(responses)
							email, _ := Index[string](responses, 1)

							return "pipego-" + email, nil
						},
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							time.Sleep(1 * time.Second)
							fmt.Println(responses)
							return "pipego-" + args.Password, nil
						},
					),
				},
			},
			want: PipelineResponses{
				[]any{
					"username",
					"email",
					"fariz-jnck",
					"pipego-username",
					"pipego-email",
					"pipego-fariz-jnck",
				},
			},
		},
		{
			name: "create with several pipe",
			args: args{
				args: createUserRequest{
					Username: "username",
					Email:    "email",
					Password: "bisacepatluarbinasa",
				},
				funcs: []Func[createUserRequest]{
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						return args.Username, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						username := Get[string](responses)
						fmt.Println(username)

						return args.Email, nil
					},
					func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
						email := Get[string](responses)
						fmt.Println(email)
						return args.Password, nil
					},
					PipelineGo(
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							time.Sleep(1 * time.Second)
							username, _ := Index[string](responses, 0)

							return "pipego-" + username, nil
						},
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							time.Sleep(1 * time.Second)
							email, _ := Index[string](responses, 1)

							return "pipego-" + email, nil
						},
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							time.Sleep(1 * time.Second)
							return "pipego-" + args.Password, nil
						},
					),
					Pipe(
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							username, _ := Index[string](responses, 0)
							return "pipe2-" + username, nil
						},
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							email, _ := Index[string](responses, 1)
							return "pipe2-" + email, nil
						},
						func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
							return "pipe2-" + args.Password, nil
						},
						Pipe(
							func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
								username, _ := Index[string](responses, 0)
								return "nestedpipe2-" + username, nil
							},
							func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
								email, _ := Index[string](responses, 1)
								return "nestedpipe2-" + email, nil
							},
							func(args createUserRequest, responses ResponsesImplementor) (response any, err error) {
								return "nestedpipe2-" + args.Password, nil
							},
						),
					),
				},
			},
			want: PipelineResponses{
				[]any{
					"username",
					"email",
					"bisacepatluarbinasa",
					"pipego-username",
					"pipego-email",
					"pipego-bisacepatluarbinasa",
					"pipe2-username",
					"pipe2-email",
					"pipe2-bisacepatluarbinasa",
					"nestedpipe2-username",
					"nestedpipe2-email",
					"nestedpipe2-bisacepatluarbinasa",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := Pipeline(tt.args.funcs...)
			got, _ := exec(tt.args.args)

			// fmt.Println("DAPETNYA", got)
			assert.Equal(t, tt.want, got)
		})
	}
}
