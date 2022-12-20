package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/cli/go-gh"
	graphql "github.com/cli/shurcooL-graphql"
)

func main() {
	var organization string
	var enterprise string

	flag.StringVar(&organization, "organization", "dylanrinkertestorg3", "organization")
	flag.StringVar(&enterprise, "enterprise", "mr-magoriums-wunderbar-emporium", "enterprise")

	flag.Parse()

	orgQuery(organization)
	entQuery, error := enterpriseQuery(enterprise)

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println(entQuery.Enterprise.OwnerInfo)

}

func orgQuery(org string) {
	fmt.Println("Organization: ", org)

	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}
	var query struct {
		Organization struct {
			MembersWithRole struct {
				Edges []struct {
					Node struct {
						Id    string
						Login string
					}
				}
			} `graphql:"membersWithRole(first: $first)"`
		} `graphql:"organization(login: $login)"`
	}

	variables := map[string]interface{}{
		"first": graphql.Int(10),
		"login": graphql.String(org),
	}

	err = client.Query("OrgMembers", &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query)
}

type OrganinizationQuery struct {
	Organization struct {
    IpAllowListEnabledSetting
    IpAllowListEntries {
      Edges {
        Node {
          AllowListValue
        }
      }
    } `graphql:"ipAllowListEntries(first: $first)"`
    IpAllowListForInstalledAppsEnabledSetting
    MembersCanForkPrivateRepositories
    RequiresTwoFactorAuthentication
    SamlIdentityProvider {
      Id
    }
  } `graphql:"organization(login: $login)"`


// create a type of EnterpriseQuery

type EnterpriseQuery struct {
	Enterprise struct {
		OwnerInfo struct {
			AllowPrivateRepositoryForkingSetting            string
			AllowPrivateRepositoryForkingSettingPolicyValue string
			DefaultRepositoryPermissionSetting              string
			IpAllowListEnabledSetting                       string
			IpAllowListEntries                              struct {
				Edges struct {
					Node struct {
						AllowListValue string
					}
				}
			} `graphql:"ipAllowListEntries(first: $first)"`
			IpAllowListForInstalledAppsEnabledSetting     string
			MembersCanChangeRepositoryVisibilitySetting   string
			MembersCanCreateRepositoriesSetting           string
			MembersCanDeleteIssuesSetting                 string
			MembersCanDeleteRepositoriesSetting           string
			MembersCanInviteCollaboratorsSetting          string
			MembersCanMakePurchasesSetting                string
			MembersCanUpdateProtectedBranchesSetting      string
			MembersCanViewDependencyInsightsSetting       string
			NotificationDeliveryRestrictionEnabledSetting string
			OrganizationProjectsSetting                   string
			RepositoryProjectsSetting                     string
			SamlIdentityProvider                          struct {
				Id string
			}
			TeamDiscussionsSetting   string
			TwoFactorRequiredSetting string
		}
	} `graphql:"enterprise(slug: $slug)"`
}

func enterpriseQuery(ent string) (*EnterpriseQuery, error) {
	fmt.Println("Enterprise: ", ent)

	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	query := new(EnterpriseQuery)

	variables := map[string]interface{}{
		"slug":  graphql.String(ent),
		"first": graphql.Int(10),
	}

	err = client.Query("Enterprise", &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query)

	return query, err
}
