package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/cli/go-gh/pkg/tableprinter"
	graphql "github.com/cli/shurcooL-graphql"
)

func main() {
	if err := cli(); err != nil {
		fmt.Fprintf(os.Stderr, "gh-org-transfer-audit failed: %s\n", err.Error())
		os.Exit(1)
	}
}

func cli() error {
	var organization string
	var enterprise string
	var repo repository.Repository
	var err error
	// isTerminal := term.IsTerminal(os.Stdout)

	flag.StringVar(&organization, "organization", "", "organization")
	flag.StringVar(&enterprise, "enterprise", "", "enterprise")

	flag.Parse()

	fmt.Println("Organization:", organization)
	fmt.Println("Enterprise:", enterprise)

	// if organization and enterprise are not provided, get the current repository
	if organization == "" && enterprise == "" {
		fmt.Println("No organization or enterprise provided. Retrieving policies current organization of current repository.")
		repo, err = gh.CurrentRepository()

		if err != nil {
			return fmt.Errorf("could not determine what org to use: %w", err)
		}

		organization = repo.Owner()

		orgPolicies, err := getOrganizationGQLPolicies(organization)

		if err != nil {
			return fmt.Errorf("the owner of the current repository is not an organization %w", err)
		}

		fmt.Println(orgPolicies.Organization)

	}

	var entPolicies *EnterprisePolicies
	var error error

	// if only enterprise is provided, get the enterprise policies
	if organization == "" && enterprise != "" {
		fmt.Println("No organization provided. Retrieving policies for enterprise.")
		entPolicies, error = getEnterprisePolicies(enterprise)

		if error != nil {
			log.Fatal(error)
		}

		fmt.Println(entPolicies.Enterprise)

		tablePrintEntPolicies(*entPolicies)

	}

	var orgGQLPolicies *OrganizationGQLPolicies
	var orgRESTPolicies *OrganizationRESTPolicies

	// if only organization is provided, get the organization policies
	if organization != "" && enterprise == "" {
		fmt.Println("No enterprise provided. Retrieving policies for organization.")

		// Get organization GraphQL policies
		orgGQLPolicies, error = getOrganizationGQLPolicies(organization)

		if error != nil {
			log.Fatal(error)
		}

		fmt.Println("Organization GQL Policies", orgGQLPolicies)

		orgRESTPolicies, error = getOrganizationRESTPolicies(organization)

		if error != nil {
			log.Fatal(error)
		}

		fmt.Println("Organization REST Policies", orgRESTPolicies)

		// Create a new orgPolicies object with the REST and GQL policies
		orgPolicies := OrganizationPolicies{
			GQL:  *orgGQLPolicies,
			REST: *orgRESTPolicies,
		}

		tablePrintOrgPolicies(orgPolicies)
	}

	// if both are provided, get the both policies and compare them
	if organization != "" && enterprise != "" {
		fmt.Println("Both enterprise and organization provided. Retrieving policies for both and comparing them.")
		// Get organization policies
		orgGQLPolicies, error = getOrganizationGQLPolicies(organization)

		if error != nil {
			log.Fatal(error)
		}

		// Get enterprise policies
		entPolicies, error = getEnterprisePolicies(enterprise)

		if error != nil {
			log.Fatal(error)
		}

		comparison := comparePolicies(orgGQLPolicies, entPolicies)

		fmt.Println(comparison)
	}
	// createCSV(orgPolicies, entPolicies, comparePolicies(orgPolicies, entPolicies))

	return nil
}

func tablePrintOrgPolicies(orgPolicies OrganizationPolicies) {
	// have to actually get isTerminal
	tp := tableprinter.New(os.Stdout, true, 100)

	tp.AddField("Policy Name", tableprinter.WithColor(bold))
	tp.AddField("Policy Value")
	tp.EndRow()
	tp.AddField("HasOrganizationProjects")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Has_organization_projects))
	tp.EndRow()
	tp.AddField("HasRepositoryProjects")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Has_repository_projects))
	tp.EndRow()
	tp.AddField("DefaultRepositoryPermission")
	tp.AddField(orgPolicies.REST.Default_repository_permission)
	tp.EndRow()
	tp.AddField("MembersCanCreateRepositories")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Members_can_create_repositories))
	tp.EndRow()
	tp.AddField("TwoFactorRequirementEnabled")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Two_factor_requirement_enabled))
	tp.EndRow()
	tp.AddField("MembersAllowedRepositoryCreationType")
	tp.AddField(orgPolicies.REST.Members_allowed_repository_creation_type)
	tp.EndRow()
	tp.AddField("MembersCanCreatePublicRepositories")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Members_can_create_public_repositories))
	tp.EndRow()
	tp.AddField("MembersCanCreatePrivateRepositories")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Members_can_create_private_repositories))
	tp.EndRow()
	tp.AddField("MembersCanCreateInternalRepositories")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Members_can_create_internal_repositories))
	tp.EndRow()
	tp.AddField("MembersCanCreatePages")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Members_can_create_pages))
	tp.EndRow()
	tp.AddField("MembersCanForkPrivateRepositoriesREST")
	tp.AddField(strconv.FormatBool(orgPolicies.REST.Members_can_fork_private_repositories))
	tp.EndRow()
	tp.AddField("IpAllowListEnabledSetting")
	tp.AddField(orgPolicies.GQL.Organization.IpAllowListEnabledSetting, tableprinter.WithColor(red))
	tp.EndRow()
	tp.AddField("IpAllowListEntries")
	tp.AddField(orgPolicies.GQL.Organization.IpAllowListEntries.Edges.Node.AllowListValue, tableprinter.WithColor(red))
	tp.EndRow()
	tp.AddField("IpAllowListForInstalledAppsEnabledSetting")
	tp.AddField(orgPolicies.GQL.Organization.IpAllowListForInstalledAppsEnabledSetting, tableprinter.WithColor(green))
	tp.EndRow()
	tp.AddField("MembersCanForkPrivateRepositories")
	tp.AddField(strconv.FormatBool(orgPolicies.GQL.Organization.MembersCanForkPrivateRepositories))
	tp.EndRow()
	tp.AddField("RequiresTwoFactorAuthentication")
	tp.AddField(strconv.FormatBool(orgPolicies.GQL.Organization.RequiresTwoFactorAuthentication))
	tp.EndRow()
	tp.AddField("SamlIdentityProvider")
	tp.AddField(orgPolicies.GQL.Organization.SamlIdentityProvider.Id)
	tp.EndRow()

	tp.Render()
}

func tablePrintEntPolicies(entPolicies EnterprisePolicies) {
	// have to actually get isTerminal
	tp := tableprinter.New(os.Stdout, true, 100)

	tp.AddField("Policy Name")
	tp.AddField("Policy Value")
	tp.EndRow()
	tp.AddField("AllowPrivateRepositoryForkingSettingPolicyValue")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.AllowPrivateRepositoryForkingSetting)
	// tp.AddField(entPolicies.Enterprise.OwnerInfo.AllowPrivateRepositoryForkingSettingPolicyValue)
	tp.EndRow()
	tp.AddField("DefaultRepositoryPermissionSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.DefaultRepositoryPermissionSetting)
	tp.EndRow()
	tp.AddField("IpAllowListEnabledSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.IpAllowListEnabledSetting)
	tp.EndRow()
	tp.AddField("IpAllowListEntries")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.IpAllowListEntries.Edges.Node.AllowListValue)
	tp.EndRow()
	tp.AddField("IpAllowListForInstalledAppsEnabledSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.IpAllowListForInstalledAppsEnabledSetting)
	tp.EndRow()
	tp.AddField("MembersCanChangeRepositoryVisibilitySetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanChangeRepositoryVisibilitySetting)
	tp.EndRow()
	tp.AddField("MembersCanCreateRepositoriesSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanCreateRepositoriesSetting)
	tp.EndRow()
	tp.AddField("MembersCanDeleteIssuesSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanDeleteIssuesSetting)
	tp.EndRow()
	tp.AddField("MembersCanDeleteRepositoriesSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanDeleteRepositoriesSetting)
	tp.EndRow()
	tp.AddField("MembersCanInviteCollaboratorsSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanInviteCollaboratorsSetting)
	tp.EndRow()
	tp.AddField("MembersCanMakePurchasesSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanMakePurchasesSetting)
	tp.EndRow()
	tp.AddField("MembersCanUpdateProtectedBranchesSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanUpdateProtectedBranchesSetting)
	tp.EndRow()
	tp.AddField("MembersCanViewDependencyInsightsSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.MembersCanViewDependencyInsightsSetting)
	tp.EndRow()
	tp.AddField("OrganizationProjectsSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.OrganizationProjectsSetting)
	tp.EndRow()
	tp.AddField("RepositoryProjectsSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.RepositoryProjectsSetting)
	tp.EndRow()
	tp.AddField("SamlIdentityProvider")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.SamlIdentityProvider.Id)
	tp.EndRow()
	tp.AddField("TeamDiscussionsSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.TeamDiscussionsSetting)
	tp.EndRow()
	tp.AddField("TwoFactorRequiredSetting")
	tp.AddField(entPolicies.Enterprise.OwnerInfo.TwoFactorRequiredSetting)
	tp.EndRow()

	tp.Render()
}

// function that takes in a string and returns that string color red
func red(s string) string {
	return fmt.Sprintf("\033[31m%s\033[0m", s)
}

func green(s string) string {
	return fmt.Sprintf("\033[32m%s\033[0m", s)
}

func bold(s string) string {
	return fmt.Sprintf("\u001b[1m%s\u001b[0m", s)
}

type OrganizationPolicies struct {
	GQL  OrganizationGQLPolicies
	REST OrganizationRESTPolicies
}

type OrganizationGQLPolicies struct {
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

// OrganizationRESTPolicies is a struct that contains the REST response for an organization

// type OrganizationRESTPolicies struct {

type OrganizationRESTPolicies struct {
	Has_organization_projects                bool
	Has_repository_projects                  bool
	Default_repository_permission            string
	Members_can_create_repositories          bool
	Two_factor_requirement_enabled           bool
	Members_allowed_repository_creation_type string
	Members_can_create_public_repositories   bool
	Members_can_create_private_repositories  bool
	Members_can_create_internal_repositories bool
	Members_can_create_pages                 bool
	Members_can_fork_private_repositories    bool
}

func getOrganizationRESTPolicies(org string) (*OrganizationRESTPolicies, error) {
	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	response := new(OrganizationRESTPolicies)

	err = client.Get(fmt.Sprintf("orgs/%s", org), &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
	return response, err
}

func getOrganizationGQLPolicies(org string) (*OrganizationGQLPolicies, error) {
	fmt.Println("Organization: ", org)

	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	query := new(OrganizationGQLPolicies)

	variables := map[string]interface{}{
		"login": graphql.String(org),
		"first": graphql.Int(10),
	}

	err = client.Query("Organization", &query, variables)
	if err != nil {
		log.Fatal(err)
	}

	return query, err
}

// create a type of getEnterprisePolicies

type EnterprisePolicies struct {
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

func getEnterprisePolicies(ent string) (*EnterprisePolicies, error) {
	fmt.Println("Enterprise: ", ent)

	client, err := gh.GQLClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	query := new(EnterprisePolicies)

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

func comparePolicies(org *OrganizationGQLPolicies, ent *EnterprisePolicies) *Compare {
	fmt.Println("Comparing Organization and Enterprise Policies")

	compare := new(Compare)

	// fmt.Println(`Compare:`, compare)

	compare = comparePrivateRepositoryForking(compare, org, ent)

	// fmt.Println(`Compare 1:`, compare)

	compare = compareTwoFactorAuthentication(compare, org, ent)

	// fmt.Println(`Compare 2:`, compare)

	compare = compareSamlIdentityProvider(compare, org, ent)

	// fmt.Println(`Compare 3:`, compare)

	return compare
}

func comparePrivateRepositoryForking(compare *Compare, org *OrganizationGQLPolicies, ent *EnterprisePolicies) *Compare {
	fmt.Println("Comparing Private Repository Forking Policies")

	if ent.Enterprise.OwnerInfo.AllowPrivateRepositoryForkingSetting == "NO_POLICY" {
		compare.PrivateRepositoryForking.Comment = "There is no Enterprise policy."
		compare.PrivateRepositoryForking.Status = "✓"
		return compare
	}

	return compare
}

func compareTwoFactorAuthentication(compare *Compare, org *OrganizationGQLPolicies, ent *EnterprisePolicies) *Compare {
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

func compareSamlIdentityProvider(compare *Compare, org *OrganizationGQLPolicies, ent *EnterprisePolicies) *Compare {
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
				compare.TwoFactorAuthenticationSetting.Comment = "SAML Single Sign On is enabled at the Enterprise level and the Organization level. The Enterprise SAML Single Sign On provider will apply to the Organization."
				compare.TwoFactorAuthenticationSetting.Status = "✗"
				return compare
			}
		}

		return compare
	}

	return compare
}

// func createCSV(org *OrganizationGQLPolicies, ent *getEnterprisePolicies, compare *Compare) {
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
