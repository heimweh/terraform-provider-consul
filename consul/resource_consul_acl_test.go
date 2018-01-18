package consul

import (
	"fmt"
	"os"
	"testing"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccConsulACL_basic(t *testing.T) {
	if os.Getenv("CONSUL_TOKEN") == "" && os.Getenv("CONSUL_HTTP_TOKEN") == "" {
		t.Skip("Environment variable CONSUL_TOKEN or CONSUL_HTTP_TOKEN is not set")
	}

	resourceName := "consul_acl.foo"
	name := fmt.Sprintf("tf-%s", acctest.RandString(5))
	nameUpdated := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckConsulACLDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccConsulACLConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulACLExists(resourceName),
				),
			},
			resource.TestStep{
				Config: testAccConsulACLConfig_Update(nameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConsulACLExists(resourceName),
				),
			},
		},
	})
}

func testAccCheckConsulACLDestroy(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "consul_acl.foo" {
			continue
		}

		client := testAccProvider.Meta().(*consulapi.Client)
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}

		acl, _, err := client.ACL().Info(r.Primary.ID, opts)
		if err == nil {
			return fmt.Errorf("%s still exists", acl.Name)
		}
	}

	return nil
}

func testAccCheckConsulACLExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*consulapi.Client)
		opts := &consulapi.QueryOptions{Datacenter: "dc1"}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if _, _, err := client.ACL().Info(rs.Primary.ID, opts); err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckConsulACLValue(n, attr, val string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rn, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found")
		}
		out, ok := rn.Primary.Attributes["var."+attr]
		if !ok {
			return fmt.Errorf("Attribute '%s' not found: %#v", attr, rn.Primary.Attributes)
		}
		if val != "<any>" && out != val {
			return fmt.Errorf("Attribute '%s' value '%s' != '%s'", attr, out, val)
		}
		if val == "<any>" && out == "" {
			return fmt.Errorf("Attribute '%s' value '%s'", attr, out)
		}
		return nil
	}
}

func testAccConsulACLConfig(name string) string {
	return fmt.Sprintf(`
resource "consul_acl" "foo" {
  name = "%[1]v"
  type = "client"
  rules = <<EOF
key "" {
  policy = "deny"
}

key "foo/private/" {
  policy = "read"
}
EOF
}
`, name)
}

func testAccConsulACLConfig_Update(name string) string {
	return fmt.Sprintf(`
resource "consul_acl" "foo" {
  name = "%[1]v"
  type = "client"
  rules = <<EOF
key "" {
  policy = "read"
}

key "foo/private/" {
  policy = "deny"
}
EOF
}
`, name)
}
