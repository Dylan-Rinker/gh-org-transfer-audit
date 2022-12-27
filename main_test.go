package main

import (
	"fmt"
	"testing"
)

// type compareRequiresTwoFactorAuthenticationTest struct {
// 	org OrganinizationQuery
// 	ent EnterpriseQuery
// 	result Compare
// }

// orgNoPolicy := OrganinizationQuery{
// 	RequiresTwoFactorAuthentication: false,
// }

// var entNoPolicy := EnterpriseQuery{
// 	OwnerInfo: OwnerInfo{
// 		TwoFactorRequiredSetting: "NO_POLICY",
// 	},
// }

// var compareNoPolicies := Compare{
// 	TwoFactorAuthenticationSetting: TwoFactorAuthenticationSetting{
// 		Comment: "There is no Enterprise policy.",
// 		Status: "✓",
// 	},
// }

// var compareRequiresTwoFactorAuthenticationTests = []compareRequiresTwoFactorAuthenticationTest{
// 	compareRequiresTwoFactorAuthenticationTest{orgNoPolicy, entNoPolicy, compareNoPolicies}
// }

func TestCompareRequiresTwoFactorAuthentication(t *testing.T) {
	tests := map[string]struct {
		org    OrganinizationQuery
		ent    EnterpriseQuery
		result Compare
	}{
		"no policies": {
			org: OrganinizationQuery{
				// Organization{
				// 	RequiresTwoFactorAuthentication: false,
				// },
			},
			ent: EnterpriseQuery{
				// OwnerInfo: OwnerInfo{
				// 	TwoFactorRequiredSetting: "NO_POLICY",
				// },
			},
			result: Compare{
				TwoFactorAuthenticationSetting: TwoFactorAuthenticationSetting{
					Comment: "There is no Enterprise policy.",
					Status:  "✓",
				},
			},
		},
	}

	compare := new(Compare)

	for _, test := range tests {
		result := compareTwoFactorAuthentication(compare, &test.org, &test.ent)

		testing := test.ent.Enterprise.OwnerInfo.TwoFactorRequiredSetting

		fmt.Println("testing", testing)

		if testing == "" {
			fmt.Println("testing is empty")
		}

		fmt.Println(test.ent.Enterprise.OwnerInfo.TwoFactorRequiredSetting)
		if *result != test.result {
			t.Error("Expected", test.result, "got", result)
		}
	}
}

// func TestComparePrivateRepoForking(t *testing.T) {
// 	org := OrganinizationQuery{
// 		Organization {
// 			RequiresTwoFactorAuthentication: true,
// 	},
// }
// 	ent := Enterprise{
// 		MembersCanForkPrivateRepositories: false,
// 	}
// 	if comparePrivateRepoForking(org, ent) != "true" {
// 		t.Error("Expected true, got false")
// 	}

// 	// Instantiate a new OrganinizationQuery struct
// 	org = OrganinizationQuery{
// 		RequiresTwoFactorAuthentication: false,
// 	}

// }
