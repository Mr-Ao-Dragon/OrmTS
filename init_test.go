package OrmTS

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	token := Token{
		AccessKeyId:     os.Getenv("ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ACCESS_KEY_SECRET"),
	}
	roleArn := os.Getenv("roleArnWithTest")
	roleSessionName := "test"
	Policy := os.Getenv("Policy")
	if os.Getenv("SECURITY_TOKEN") == "" {
		_ = token.GetSTSToken(&roleArn, nil, &roleSessionName, &Policy)
	}
	type args struct {
		region string
		token  Token
	}
	tests := []struct {
		name            string
		args            args
		wantNewQueryObj Query
		wantErr         bool
	}{
		{
			name: "test",
			args: args{
				region: "cn-beijing",
				token:  token,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.args.region, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
