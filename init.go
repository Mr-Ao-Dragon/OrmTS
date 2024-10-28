package OrmTS

import (
	"errors"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	otsapi "github.com/alibabacloud-go/tablestore-20201209/client"
	teautil "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if regionID, RegionEnvSet := os.LookupEnv("REGION"); !RegionEnvSet || regionID == "" {
		log.Fatal().
			AnErr("err: ", errors.New("need region env, but is not found")).
			Msgf("plase set %#v", regionID)
	}
	if useVpc, vpcEnvSet := os.LookupEnv("USE_VPC"); !vpcEnvSet || useVpc == "" {
		log.Info().Msgf("default use to no vpc mode.")
		_ = os.Setenv("USE_VPC", "false")
		UseVpc = false
		EndPoint = "https://tablestore." + os.Getenv("REGION") + ".aliyuncs.com"
	} else if useVpc := os.Getenv("USE_VPC"); useVpc != "false" && useVpc != "true" {
		log.Fatal().Msgf("env %#v must be true or false", useVpc)
	} else {
		UseVpc = true
		EndPoint = "https://tablestore-vpc." + os.Getenv("REGION") + ".aliyuncs.com"
	}
	log.Info().Msgf("hello, enduser!")
	log.Info().Msgf("is me! OrmTS!")
	log.Debug().Msgf("hey developer, let's rock!")
}
func New(region string, token Token) (NewQueryObj Query, err error) {
	NewQueryObj = *new(Query)
	NewQueryObj.region = new(RegionData)
	NewQueryObj.region.Region = region
	cfg := new(openapi.Config)
	cfg.SetAccessKeyId(token.AccessKeyId)
	cfg.SetAccessKeySecret(token.AccessKeySecret)
	cfg.SetSecurityToken(token.SecurityToken)
	cfg.SetRegionId(region)
	cfg.SetEndpoint("tablestore." + region + ".aliyuncs.com")
	listQueryObj, err := otsapi.NewClient(cfg)
	runtime := new(teautil.RuntimeOptions)
	headers := make(map[string]*string)
	lsDbsReq := new(otsapi.ListInstancesRequest)
	lsDbsResult, err := listQueryObj.ListInstancesWithOptions(lsDbsReq, headers, runtime)
	if err != nil {
		return Query{}, err
	}
	instanceNums := len(lsDbsResult.Body.Instances)
	if instanceNums == 0 {
		return Query{
			token: token,
			region: &RegionData{
				Region: region,
				Dbs:    []*otsapi.ListInstancesResponseBodyInstances{},
			},
			Setting: new(tablestore.TableStoreClient),
		}, nil
	}

	NewQueryObj.token = token
	NewQueryObj.Setting = new(tablestore.TableStoreClient)
	NewQueryObj.region = new(RegionData)
	NewQueryObj.region.Region = region
	NewQueryObj.region.Dbs = lsDbsResult.Body.Instances
	err = nil
	return
}
