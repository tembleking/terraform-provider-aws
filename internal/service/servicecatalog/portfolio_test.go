package servicecatalog_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicecatalog"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
)

func TestAccServiceCatalogPortfolio_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_servicecatalog_portfolio.test"
	name := sdkacctest.RandString(5)
	var dpo servicecatalog.DescribePortfolioOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicecatalog.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckServiceCatlaogPortfolioDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPortfolioConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortfolio(ctx, resourceName, &dpo),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "catalog", regexp.MustCompile(`portfolio/.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "created_time"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "test-2"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "test-3"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccServiceCatalogPortfolio_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	name := sdkacctest.RandString(5)
	resourceName := "aws_servicecatalog_portfolio.test"
	var dpo servicecatalog.DescribePortfolioOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicecatalog.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckServiceCatlaogPortfolioDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPortfolioConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortfolio(ctx, resourceName, &dpo),
					testAccCheckServiceCatlaogPortfolioDisappears(ctx, &dpo),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccServiceCatalogPortfolio_tags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_servicecatalog_portfolio.test"
	name := sdkacctest.RandString(5)
	var dpo servicecatalog.DescribePortfolioOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ErrorCheck:               acctest.ErrorCheck(t, servicecatalog.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckServiceCatlaogPortfolioDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPortfolioConfig_tags1(name, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortfolio(ctx, resourceName, &dpo),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPortfolioConfig_tags2(name, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortfolio(ctx, resourceName, &dpo),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccPortfolioConfig_tags1(name, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPortfolio(ctx, resourceName, &dpo),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckPortfolio(ctx context.Context, pr string, dpo *servicecatalog.DescribePortfolioOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ServiceCatalogConn()
		rs, ok := s.RootModule().Resources[pr]
		if !ok {
			return fmt.Errorf("Not found: %s", pr)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		input := servicecatalog.DescribePortfolioInput{}
		input.Id = aws.String(rs.Primary.ID)

		resp, err := conn.DescribePortfolioWithContext(ctx, &input)
		if err != nil {
			return err
		}

		*dpo = *resp
		return nil
	}
}

func testAccCheckServiceCatlaogPortfolioDisappears(ctx context.Context, dpo *servicecatalog.DescribePortfolioOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ServiceCatalogConn()

		input := servicecatalog.DeletePortfolioInput{}
		input.Id = dpo.PortfolioDetail.Id

		_, err := conn.DeletePortfolioWithContext(ctx, &input)
		return err
	}
}

func testAccCheckServiceCatlaogPortfolioDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).ServiceCatalogConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_servicecatalog_portfolio" {
				continue
			}
			input := servicecatalog.DescribePortfolioInput{}
			input.Id = aws.String(rs.Primary.ID)

			_, err := conn.DescribePortfolioWithContext(ctx, &input)
			if err == nil {
				return fmt.Errorf("Portfolio still exists")
			}
		}

		return nil
	}
}

func testAccPortfolioConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "aws_servicecatalog_portfolio" "test" {
  name          = "%s"
  description   = "test-2"
  provider_name = "test-3"
}
`, name)
}

func testAccPortfolioConfig_tags1(name, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_servicecatalog_portfolio" "test" {
  name          = %[1]q
  description   = "test-b"
  provider_name = "test-c"

  tags = {
    %[2]q = %[3]q
  }
}
`, name, tagKey1, tagValue1)
}

func testAccPortfolioConfig_tags2(name, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_servicecatalog_portfolio" "test" {
  name          = %[1]q
  description   = "test-only-change-me"
  provider_name = "test-c"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, name, tagKey1, tagValue1, tagKey2, tagValue2)
}
