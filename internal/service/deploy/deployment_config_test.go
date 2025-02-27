package deploy_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codedeploy"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccDeployDeploymentConfig_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var config1 codedeploy.DeploymentConfigInfo
	resourceName := "aws_codedeploy_deployment_config.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, codedeploy.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfigConfig_fleet(rName, 75),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config1),
					resource.TestCheckResourceAttr(resourceName, "deployment_config_name", rName),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Server"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDeployDeploymentConfig_fleetPercent(t *testing.T) {
	ctx := acctest.Context(t)
	var config1, config2 codedeploy.DeploymentConfigInfo
	resourceName := "aws_codedeploy_deployment_config.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, codedeploy.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfigConfig_fleet(rName, 75),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config1),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.type", "FLEET_PERCENT"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.value", "75"),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Server"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "0"),
				),
			},
			{
				Config: testAccDeploymentConfigConfig_fleet(rName, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config2),
					testAccCheckDeploymentConfigRecreated(&config1, &config2),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.type", "FLEET_PERCENT"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.value", "50"),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Server"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDeployDeploymentConfig_hostCount(t *testing.T) {
	ctx := acctest.Context(t)
	var config1, config2 codedeploy.DeploymentConfigInfo
	resourceName := "aws_codedeploy_deployment_config.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, codedeploy.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfigConfig_hostCount(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config1),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.type", "HOST_COUNT"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Server"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "0"),
				),
			},
			{
				Config: testAccDeploymentConfigConfig_hostCount(rName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config2),
					testAccCheckDeploymentConfigRecreated(&config1, &config2),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.type", "HOST_COUNT"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.0.value", "2"),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Server"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDeployDeploymentConfig_trafficCanary(t *testing.T) {
	ctx := acctest.Context(t)
	var config1, config2 codedeploy.DeploymentConfigInfo
	resourceName := "aws_codedeploy_deployment_config.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, codedeploy.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfigConfig_trafficCanary(rName, 10, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config1),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Lambda"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.type", "TimeBasedCanary"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.0.interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.0.percentage", "50"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "0"),
				),
			},
			{
				Config: testAccDeploymentConfigConfig_trafficCanary(rName, 3, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config2),
					testAccCheckDeploymentConfigRecreated(&config1, &config2),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Lambda"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.type", "TimeBasedCanary"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.0.interval", "3"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.0.percentage", "10"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDeployDeploymentConfig_trafficLinear(t *testing.T) {
	ctx := acctest.Context(t)
	var config1, config2 codedeploy.DeploymentConfigInfo
	resourceName := "aws_codedeploy_deployment_config.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, codedeploy.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDeploymentConfigDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentConfigConfig_trafficLinear(rName, 10, 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config1),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Lambda"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.type", "TimeBasedLinear"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.0.interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.0.percentage", "50"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "0"),
				),
			},
			{
				Config: testAccDeploymentConfigConfig_trafficLinear(rName, 3, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeploymentConfigExists(ctx, resourceName, &config2),
					testAccCheckDeploymentConfigRecreated(&config1, &config2),
					resource.TestCheckResourceAttr(resourceName, "compute_platform", "Lambda"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.type", "TimeBasedLinear"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.0.interval", "3"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_linear.0.percentage", "10"),
					resource.TestCheckResourceAttr(resourceName, "traffic_routing_config.0.time_based_canary.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_hosts.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDeploymentConfigDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DeployConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_codedeploy_deployment_config" {
				continue
			}

			resp, err := conn.GetDeploymentConfigWithContext(ctx, &codedeploy.GetDeploymentConfigInput{
				DeploymentConfigName: aws.String(rs.Primary.ID),
			})

			if tfawserr.ErrCodeEquals(err, codedeploy.ErrCodeDeploymentConfigDoesNotExistException) {
				continue
			}

			if err == nil {
				if resp.DeploymentConfigInfo != nil {
					return fmt.Errorf("CodeDeploy deployment config still exists:\n%#v", *resp.DeploymentConfigInfo.DeploymentConfigName)
				}
			}

			return err
		}

		return nil
	}
}

func testAccCheckDeploymentConfigExists(ctx context.Context, name string, config *codedeploy.DeploymentConfigInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DeployConn()

		resp, err := conn.GetDeploymentConfigWithContext(ctx, &codedeploy.GetDeploymentConfigInput{
			DeploymentConfigName: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return err
		}

		*config = *resp.DeploymentConfigInfo

		return nil
	}
}

func testAccCheckDeploymentConfigRecreated(i, j *codedeploy.DeploymentConfigInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if aws.TimeValue(i.CreateTime).Equal(aws.TimeValue(j.CreateTime)) {
			return errors.New("CodeDeploy Deployment Config was not recreated")
		}

		return nil
	}
}

func testAccDeploymentConfigConfig_fleet(rName string, value int) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_config" "test" {
  deployment_config_name = %q

  minimum_healthy_hosts {
    type  = "FLEET_PERCENT"
    value = %d
  }
}
`, rName, value)
}

func testAccDeploymentConfigConfig_hostCount(rName string, value int) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_config" "test" {
  deployment_config_name = %q

  minimum_healthy_hosts {
    type  = "HOST_COUNT"
    value = %d
  }
}
`, rName, value)
}

func testAccDeploymentConfigConfig_trafficCanary(rName string, interval, percentage int) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_config" "test" {
  deployment_config_name = %q
  compute_platform       = "Lambda"

  traffic_routing_config {
    type = "TimeBasedCanary"

    time_based_canary {
      interval   = %d
      percentage = %d
    }
  }
}
`, rName, interval, percentage)
}

func testAccDeploymentConfigConfig_trafficLinear(rName string, interval, percentage int) string {
	return fmt.Sprintf(`
resource "aws_codedeploy_deployment_config" "test" {
  deployment_config_name = %q
  compute_platform       = "Lambda"

  traffic_routing_config {
    type = "TimeBasedLinear"

    time_based_linear {
      interval   = %d
      percentage = %d
    }
  }
}
`, rName, interval, percentage)
}
