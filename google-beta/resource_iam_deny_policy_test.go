package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIAM2DenyPolicy_iamDenyPolicyUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
		"random_suffix":   randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersOiCS,
		CheckDestroy: testAccCheckIAM2DenyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyUpdate(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyUpdate2(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyUpdate(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
		},
	})
}

func TestAccIAM2DenyPolicy_iamDenyPolicyFolderParent(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProvidersOiCS,
		CheckDestroy: testAccCheckIAM2DenyPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccIAM2DenyPolicy_iamDenyPolicyFolder(context),
			},
			{
				ResourceName:            "google_iam_deny_policy.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "parent"},
			},
		},
	})
}

func testAccIAM2DenyPolicy_iamDenyPolicyUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "project" {
  provider        = google-beta
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_iam_deny_policy" "example" {
  provider = google-beta
  parent   = urlencode("cloudresourcemanager.googleapis.com/projects/${google_project.project.project_id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "First rule"
    deny_rule {
      denied_principals = ["principal://iam.googleapis.com/projects/-/serviceAccounts/${google_service_account.test-account.email}"]
      denial_condition {
        title = "Some expr"
        expression = "!resource.matchTag('12345678/env', 'test')"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.delete"]
    }
  }
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denial_condition {
        title = "Some expr"
        expression = "!resource.matchTag('12345678/env', 'test')"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.delete"]
      exception_principals = ["principal://iam.googleapis.com/projects/-/serviceAccounts/${google_service_account.test-account.email}"]
    }
  }
}

resource "google_service_account" "test-account" {
  provider = google-beta
  account_id   = "tf-test-deny-account%{random_suffix}"
  display_name = "Test Service Account"
  project      = google_project.project.project_id
}
`, context)
}

func testAccIAM2DenyPolicy_iamDenyPolicyUpdate2(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "project" {
  provider        = google-beta
  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_iam_deny_policy" "example" {
  provider = google-beta
  parent   = urlencode("cloudresourcemanager.googleapis.com/projects/${google_project.project.project_id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denial_condition {
        title = "Some other expr"
        expression = "!resource.matchTag('87654321/env', 'test')"
        location = "/some/file"
        description = "A denial condition"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.delete"]
    }
  }
}

resource "google_service_account" "test-account" {
  provider = google-beta
  account_id   = "tf-test-deny-account%{random_suffix}"
  display_name = "Test Service Account"
  project      = google_project.project.project_id
}
`, context)
}

func testAccIAM2DenyPolicy_iamDenyPolicyFolder(context map[string]interface{}) string {
	return Nprintf(`
resource "google_iam_deny_policy" "example" {
  provider = google-beta
  parent   = urlencode("cloudresourcemanager.googleapis.com/${google_folder.folder.id}")
  name     = "tf-test-my-deny-policy%{random_suffix}"
  display_name = "A deny rule"
  rules {
    description = "Second rule"
    deny_rule {
      denied_principals = ["principalSet://goog/public:all"]
      denial_condition {
        title = "Some expr"
        expression = "!resource.matchTag('12345678/env', 'test')"
      }
      denied_permissions = ["cloudresourcemanager.googleapis.com/projects.delete"]
    }
  }
}

resource "google_folder" "folder" {
  provider = google-beta
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
}
`, context)
}
