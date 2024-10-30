package OrmTS

import (
	"errors"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	stsapi "github.com/alibabacloud-go/sts-20150401/v2/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/bytedance/sonic"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func (receiver *Token) GetSTSToken(roleArn *string, ttl *int64, roleSessionName *string, Policy ...*string) error {
	var endPoint string
	switch os.Getenv("USE_VPC") {
	case "false":
		endPoint = "sts." + os.Getenv("REGION") + ".aliyuncs.com"
	case "true":
		endPoint = "sts-vpc." + os.Getenv("REGION") + ".aliyuncs.com"
	default:
		endPoint = "sts." + os.Getenv("REGION") + ".aliyuncs.com"
	}
	if ttl == nil || *ttl <= 900 {
		ttl = new(int64)
		*ttl = 900
	}
	if roleArn == nil || *roleArn == "" {
		return errors.New("must have roleArn")
	}
	openapiConfig := &openapi.Config{
		AccessKeyId:     &receiver.AccessKeyId,
		AccessKeySecret: &receiver.AccessKeySecret,
		Endpoint:        &endPoint,
	}
	stsApiClient, err := stsapi.NewClient(openapiConfig)
	if err != nil {
		return err
	}
	assumeRoleRequest := new(stsapi.AssumeRoleRequest)
	switch len(Policy) {
	case 0:
		assumeRoleRequest = &stsapi.AssumeRoleRequest{
			RoleArn:         roleArn,
			DurationSeconds: ttl,
			RoleSessionName: roleSessionName,
		}
	case 1:
		assumeRoleRequest = &stsapi.AssumeRoleRequest{
			RoleArn:         roleArn,
			DurationSeconds: ttl,
			Policy:          Policy[0],
			RoleSessionName: roleSessionName,
		}
	default:
		return errors.New("too many arguments")
	}
	runtime := &teautil.RuntimeOptions{}
	resp, err := stsApiClient.AssumeRoleWithOptions(assumeRoleRequest, runtime)
	if err != nil {
		return err
	}
	if *resp.StatusCode != 200 {
		return errors.New("assume role failed")
	}
	receiver.AccessKeyId = *resp.Body.Credentials.AccessKeyId
	receiver.AccessKeySecret = *resp.Body.Credentials.AccessKeySecret
	receiver.SecurityToken = *resp.Body.Credentials.SecurityToken
	expTime, err := time.Parse("2006-01-02T15:04:05Z", *resp.Body.Credentials.Expiration)
	if err != nil {
		return errors.Join(err, errors.New(fmt.Sprintf("Request ID: %s", *resp.Body.RequestId)))
	}
	receiver.Expiration = int(expTime.Unix())
	log.Info().Msg("get sts token success")
	log.Info().Msgf("request ID: %s", *resp.Body.RequestId)
	return nil
}
func (receiver *Policy) Decode(policy string) error {
	return sonic.UnmarshalString(policy, receiver)
}
func (receiver *Policy) Encode() (string, error) {
	return sonic.MarshalString(receiver)
}
