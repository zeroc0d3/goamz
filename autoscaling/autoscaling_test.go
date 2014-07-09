package autoscaling_test

import (
	"github.com/mitchellh/goamz/autoscaling"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/testutil"
	. "github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	autoscaling *autoscaling.AutoScaling
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	auth := aws.Auth{"abc", "123", ""}
	s.autoscaling = autoscaling.NewWithClient(auth, aws.Region{AutoScalingEndpoint: testServer.URL}, testutil.DefaultClient)
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_CreateAutoScalingGroup(c *C) {
	testServer.Response(200, nil, CreateAutoScalingGroupExample)

	options := autoscaling.CreateAutoScalingGroup{
		AvailZone:               []string{"us-east-1a"},
		DefaultCooldown:         30,
		DesiredCapacity:         2,
		HealthCheckGracePeriod:  30,
		HealthCheckType:         "elb",
		InstanceId:              "i-foo",
		LaunchConfigurationName: "foobar",
		MinSize:                 2,
		MaxSize:                 2,
		PlacementGroup:          "foobar",
		Name:                    "foobar",
		Tags: []autoscaling.Tag{
			autoscaling.Tag{
				Key:   "foo",
				Value: "bar",
			},
		},
		VPCZoneIdentifier: []string{"foo", "bar"},
	}

	resp, err := s.autoscaling.CreateAutoScalingGroup(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateAutoScalingGroup"})
	c.Assert(req.Form["InstanceId"], DeepEquals, []string{"i-foo"})
	c.Assert(req.Form["VPCZoneIdentifier"], DeepEquals, []string{"foo,bar"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) Test_CreateLaunchConfiguration(c *C) {
	testServer.Response(200, nil, CreateLaunchConfigurationExample)

	options := autoscaling.CreateLaunchConfiguration{
		SecurityGroups: []string{"sg-1111"},
		ImageId:      gs  "i-141421",
		InstanceId:     "i-141421",
		InstanceType:   "m1.small",
		KeyName:        "foobar",
		Name:           "i-141421",
	}

	resp, err := s.autoscaling.CreateLaunchConfiguration(&options)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], DeepEquals, []string{"CreateLaunchConfiguration"})
	c.Assert(req.Form["InstanceType"], DeepEquals, []string{"m1.small"})
	c.Assert(req.Form["SecurityGroups.member.1"], DeepEquals, []string{"sg-1111"})
	c.Assert(err, IsNil)
	c.Assert(resp.RequestId, Equals, "7c6e177f-f082-11e1-ac58-3714bEXAMPLE")
}
