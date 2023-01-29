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

	// orgQuery(organization)
	entQuery, error := enterpriseQuery(enterprise)

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println(entQuery.Enterprise.OwnerInfo)

	orgQuery, error := organinizationQuery(organization)

	if error != nil {
		log.Fatal(error)
	}

	fmt.Println(orgQuery.Organization)

	comparison := compare(orgQuery, entQuery)

	fmt.Println(comparison)

	// createCSV(orgQuery, entQuery, compare(orgQuery, entQuery))

	CreatePDF()

}

type OrganinizationQuery struct {
	Organization struct {
		IpAllowListEnabledSetting string
		IpAllowListEntries        struct {
			Edges struct {
				Node struct {
					AllowListValue string
				}
			}
		} `graphql:"ipAllowListEntries(first: $first)"`
		IpAllowListForInstalledAppsEnabledSetting string
		MembersCanForkPrivateRepositories         bool
		RequiresTwoFactorAuthentication           bool
		SamlIdentityProvider                      struct {
			Id string
		}
	} `graphql:"organization(login: $login)"`
}

func organinizationQuery(org string) (*OrganinizationQuery, error) {
	fmt.Println("Organization: ", org)

	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	query := new(OrganinizationQuery)

	variables := map[string]interface{}{
		"login": graphql.String(org),
		"first": graphql.Int(10),
	}

	err = client.Query("Organization", &query, variables)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query)

	return query, err
}

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
	// fmt.Println(query)

	return query, err
}

type Compare struct {
	PrivateRepositoryForking
	DefaultRepositoryPermissionSetting
	SamlIdentityProvider
	TwoFactorAuthenticationSetting
}

type PrivateRepositoryForking struct {
	Policy   string `default:"Private Repository Forking"`
	Category string `default:"repository"`
	Comment  string
	Status   string
}

type DefaultRepositoryPermissionSetting struct {
	Policy   string `default:"Default Repository Permission Setting"`
	Category string `default:"repository"`
	Comment  string
	Status   string
}

type SamlIdentityProvider struct {
	Policy   string `default:"SAML Identity Provider"`
	Category string `default:"account"`
	Comment  string
	Status   string
}

type TwoFactorAuthenticationSetting struct {
	Policy   string `default:"Two Factor Authentication Setting"`
	Category string `default:"account"`
	Comment  string
	Status   string
}

func compare(org *OrganinizationQuery, ent *EnterpriseQuery) *Compare {
	fmt.Println("Comparing Organization and Enterprise Policies")

	compare := new(Compare)

	fmt.Println(`Compare:`, compare)

	compare = comparePrivateRepositoryForking(compare, org, ent)

	fmt.Println(`Compare 1:`, compare)

	compare = compareSamlIdentityProvider(compare, org, ent)

	fmt.Println(`Compare 2:`, compare)

	compare = compareTwoFactorAuthentication(compare, org, ent)

	fmt.Println(`Compare 3:`, compare)

	return compare
}

func comparePrivateRepositoryForking(compare *Compare, org *OrganinizationQuery, ent *EnterpriseQuery) *Compare {
	fmt.Println("Comparing Private Repository Forking Policies")

	if ent.Enterprise.OwnerInfo.AllowPrivateRepositoryForkingSetting == "NO_POLICY" {
		compare.PrivateRepositoryForking.Comment = "There is no Enterprise policy."
		compare.PrivateRepositoryForking.Status = "✓"
		return compare
	}

	return compare
}

func compareTwoFactorAuthentication(compare *Compare, org *OrganinizationQuery, ent *EnterpriseQuery) *Compare {
	fmt.Println("Comparing Two Factor Authentication Policies")

	if ent.Enterprise.OwnerInfo.TwoFactorRequiredSetting == "NO_POLICY" {
		compare.TwoFactorAuthenticationSetting.Comment = "There is no Enterprise policy."
		compare.TwoFactorAuthenticationSetting.Status = "✓"
		return compare
	}

	if ent.Enterprise.OwnerInfo.TwoFactorRequiredSetting == "ENABLED" {
		if org.Organization.RequiresTwoFactorAuthentication {
			compare.TwoFactorAuthenticationSetting.Comment = "The Enterprise and the Organization both have Two Factor Authentication enabled."
			compare.TwoFactorAuthenticationSetting.Status = "✓"
		}

		if !org.Organization.RequiresTwoFactorAuthentication {
			compare.TwoFactorAuthenticationSetting.Comment = "The Enterprise two factor authentication setting will apply to the Organization. Members who do not have two factor authentication enabled will not be removed from the Organization."
			compare.TwoFactorAuthenticationSetting.Status = "✗"
		}

		return compare

	}

	return compare
}

func compareSamlIdentityProvider(compare *Compare, org *OrganinizationQuery, ent *EnterpriseQuery) *Compare {
	fmt.Println("Comparing SAML Identity Provider Policies")

	if ent.Enterprise.OwnerInfo.SamlIdentityProvider.Id == "" {
		compare.TwoFactorAuthenticationSetting.Comment = "SAML Single Sign On is not enabled at the Enterprise level."
		compare.TwoFactorAuthenticationSetting.Status = "✓"
		return compare
	}

	if ent.Enterprise.OwnerInfo.SamlIdentityProvider.Id != "" {
		if org.Organization.SamlIdentityProvider.Id == "" {
			compare.TwoFactorAuthenticationSetting.Comment = "SAML Single Sign On is enabled at the Enterprise level, but not at the Organization level."
			compare.TwoFactorAuthenticationSetting.Status = "✓"
			return compare
		}

		if org.Organization.SamlIdentityProvider.Id != "" {
			if org.Organization.SamlIdentityProvider.Id != ent.Enterprise.OwnerInfo.SamlIdentityProvider.Id {
				compare.TwoFactorAuthenticationSetting.Comment = "SAML Single Sign On is enabled at the Enterprise level and the Organization level. The Enterprise SAML Single Sign On provider will override the Organization's SAML Single Sign On."
				compare.TwoFactorAuthenticationSetting.Status = "✗"
				return compare
			}
		}

		return compare
	}

	return compare
}

// func createCSV(org *OrganinizationQuery, ent *EnterpriseQuery, compare *Compare) {
// 	fmt.Println("Creating CSV")
// 	csvWriter := createCSVFile()

// 	fmt.Println(csvWriter)

// 	csvWriter.Write([]string{"Policy", "Category", "Organization", "Enterprise", "Comment", "Status"})

// 	csvWriter.Write([]string{"Allow Private Repository Forking", "Repository", org.MembersCanForkPrivateRepositories, ent.OwnerInfo.AllowPrivateRepositoryForkingSetting, compare.PrivateRepositoryForking})

// 	csvWriter.Flush()
// }

// func createCSVFile() *csv.Writer {
// 	csvFile, err := os.Create("output.csv")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	csvWriter := csv.NewWriter(csvFile)
// 	return csvWriter
// }
