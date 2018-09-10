package sewansdk

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		ID                  int
		InputToken          string
		InputURL            string
		InputEnterpriseSlug string
		OutputAPI           API
	}{
		{1,
			wrongAPIToken,
			rightAPIURL,
			unitTestEnterprise,
			API{wrongAPIToken, rightAPIURL, unitTestEnterprise, APIMeta{}, nil},
		},
		{2,
			rightAPIToken,
			wrongAPIURL,
			unitTestEnterprise,
			API{rightAPIToken, wrongAPIURL, unitTestEnterprise, APIMeta{}, nil},
		},
		{3,
			wrongAPIToken,
			wrongAPIURL,
			unitTestEnterprise,
			API{wrongAPIToken, wrongAPIURL, unitTestEnterprise, APIMeta{}, nil},
		},
		{4,
			rightAPIToken,
			rightAPIURL,
			unitTestEnterprise,
			API{rightAPIToken, rightAPIURL, unitTestEnterprise, APIMeta{}, nil},
		},
	}
	fakeAPItools := APITooler{
		APIImplementer: FakeAirDrumResourceAPIer{},
		Initialyser:    Initialyser{},
	}
	for _, testCase := range testCases {
		api := fakeAPItools.Initialyser.New(
			testCase.InputToken,
			testCase.InputURL,
			testCase.InputEnterpriseSlug,
		)
		switch {
		case api.Token != testCase.OutputAPI.Token:
			t.Errorf("\n\nTC %d : API token was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, api.Token, testCase.OutputAPI.Token)
		case api.URL != testCase.OutputAPI.URL:
			t.Errorf("\n\nTC %d : API URL was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, api.URL, testCase.OutputAPI.URL)
		case api.Enterprise != testCase.OutputAPI.Enterprise:
			t.Errorf("\n\nTC %d : API enterprise was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, api.Enterprise, testCase.OutputAPI.Enterprise)
		}
	}
}

func TestCheckCloudDcStatus(t *testing.T) {
	testCases := []struct {
		ID               int
		InputAPI         *API
		TCResourceTooler Resourceer
		Err              error
	}{
		{1,
			&API{
				wrongAPIToken,
				rightAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongTokenError),
		},
		{2,
			&API{
				rightAPIToken,
				wrongAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongAPIURLError),
		},
		{3,
			&API{
				wrongAPIToken,
				wrongAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			FakeResourceResourceer{},
			errors.New(wrongAPIURLError),
		},
		{4,
			&API{
				rightAPIToken,
				rightAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			FakeResourceResourceer{},
			nil,
		},
	}
	fakeAPItools := APITooler{}
	fakeAPItools.Initialyser = Initialyser{}
	fakeClientTooler := &ClientTooler{
		Client: HTTPClienterDummy{},
	}
	fakeResourceTooler := &ResourceTooler{}
	for _, testCase := range testCases {
		fakeAPItools.APIImplementer = FakeAirDrumResourceAPIer{}
		fakeResourceTooler.Resource = testCase.TCResourceTooler
		err := fakeAPItools.Initialyser.CheckCloudDcStatus(testCase.InputAPI,
			fakeClientTooler,
			fakeResourceTooler)
		switch {
		case err == nil || testCase.Err == nil:
			if !(err == nil && testCase.Err == nil) {
				t.Errorf("\n\nTC %d : Check API error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err, testCase.Err)
			}
		case err.Error() != testCase.Err.Error():
			t.Errorf("\n\nTC %d : Check API error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.Err.Error())
		}
	}
}

func TestGetClouddcEnvMeta(t *testing.T) {
	testCases := []struct {
		ID            int
		InputAPI      *API
		TcClienter    Clienter
		OutputAPIMeta *APIMeta
		Err           error
	}{
		{
			1,
			&API{
				wrongAPIToken,
				rightAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			HTTPClienterDummy{},
			&APIMeta{},
			nil,
		},
		{
			2,
			&API{
				wrongAPIToken,
				rightAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			getListSuccessHTTPClienterFake{},
			&APIMeta{
				nonCriticalResourceMetaDataList,
				criticalResourceMetaDataList,
				otherResourceMetaDataList},
			nil,
		},
		{
			3,
			&API{
				wrongAPIToken,
				rightAPIURL,
				unitTestEnterprise,
				APIMeta{},
				&http.Client{},
			},
			getJSONListFailureHTTPClienterFake{},
			nil,
			errEmptyResourcesList,
		},
	}
	fakeAPItools := &APITooler{}
	fakeAPItools.Initialyser = Initialyser{}
	fakeClientTooler := &ClientTooler{}
	for _, testCase := range testCases {
		fakeClientTooler.Client = testCase.TcClienter
		apiMeta, err := fakeAPItools.Initialyser.GetClouddcEnvMeta(testCase.InputAPI,
			fakeClientTooler)
		diffs := cmp.Diff(apiMeta, testCase.OutputAPIMeta)
		switch {
		case (err == nil || testCase.Err == nil):
			if !(err == nil && testCase.Err == nil) {
				t.Errorf("\n\nTC %d : GetClouddcEnvMeta error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err, testCase.Err)
			} else {
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong GetClouddcEnvMeta returned structure (-got +want) \n%s",
						testCase.ID, diffs)
				}
			}
		case err.Error() != testCase.Err.Error():
			t.Errorf("\n\nTC %d : GetClouddcEnvMeta error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.Err.Error())
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong GetClouddcEnvMeta returned structure (-got +want) \n%s",
				testCase.ID, diffs)
		}
	}
}

func TestCreateResource(t *testing.T) {
	resourceName := "Unit test resource creation"
	testCases := []struct {
		ID              int
		TcClienter      Clienter
		ResourceType    string
		CreationErr     error
		CreatedResource map[string]interface{}
	}{
		{
			1,
			VMCreationSuccessHTTPClienterFake{},
			VMResourceType,
			nil,
			noTemplateVMMap,
		},
		{
			2,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			map[string]interface{}{},
		},
		{
			3,
			ResourceCreationFailureHTTPClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(creationOperation,
				resourceName,
				errEmptyResp),
			map[string]interface{}{},
		},
		{
			4,
			HandleRespErrHTTPClienterFake{},
			VMResourceType,
			errHandleResponse,
			map[string]interface{}{},
		},
	}
	apier := AirDrumResourcesAPI{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplatesTooler := TemplatesTooler{
		TemplatesTools: TemplateTemplater{},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		respCreationMap, err := apier.CreateResource(d,
			&fakeClientTooler,
			&fakeTemplatesTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		diffs := cmp.Diff(testCase.CreatedResource, respCreationMap)
		switch {
		case err == nil || testCase.CreationErr == nil:
			if !(err == nil && testCase.CreationErr == nil) {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.CreationErr)
			} else {
				switch {
				case diffs != "":
					t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.CreationErr != nil:
			if err.Error() != testCase.CreationErr.Error() {
				t.Errorf("\n\nTC %d : resource creation error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.CreationErr.Error())
			}
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong created resource map (-got +want) \n%s",
				testCase.ID, diffs)
		}
	}
}

func TestReadResource(t *testing.T) {
	resourceName := "Unit test resource read"
	testCases := []struct {
		ID           int
		TcClienter   Clienter
		TcAPI        *API
		ResourceType string
		ReadError    error
		ReadResource map[string]interface{}
	}{
		{
			1,
			VMReadSuccessHTTPClienterFake{},
			&API{
				Token: rightAPIToken,
				URL:   rightAPIURL,
				Meta: APIMeta{
					NonCriticalResourceList: nonCriticalResourceMetaDataList,
					CriticalResourceList:    criticalResourceMetaDataList,
					OtherResourceList:       otherResourceMetaDataList,
				},
				Client: &http.Client{},
			},
			VMResourceType,
			nil,
			noTemplateVMMap,
		},
		{
			2,
			VdcReadSuccessHTTPClienterFake{},
			&API{
				Token: rightAPIToken,
				URL:   rightAPIURL,
				Meta: APIMeta{
					NonCriticalResourceList: nonCriticalResourceMetaDataList,
					CriticalResourceList:    criticalResourceMetaDataList,
					OtherResourceList:       otherResourceMetaDataList,
				},
				Client: &http.Client{},
			},
			VdcResourceType,
			nil,
			vdcReadResponseMap,
		},
		{
			3,
			HTTPClienterDummy{},
			&API{
				Token: rightAPIToken,
				URL:   rightAPIURL,
				Meta: APIMeta{
					NonCriticalResourceList: nonCriticalResourceMetaDataList,
					CriticalResourceList:    criticalResourceMetaDataList,
					OtherResourceList:       otherResourceMetaDataList,
				},
				Client: &http.Client{},
			},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
			map[string]interface{}{},
		},
		{
			4,
			ResourceReadFailureHTTPClienterFake,
			&API{
				Token: rightAPIToken,
				URL:   rightAPIURL,
				Meta: APIMeta{
					NonCriticalResourceList: nonCriticalResourceMetaDataList,
					CriticalResourceList:    criticalResourceMetaDataList,
					OtherResourceList:       otherResourceMetaDataList,
				},
				Client: &http.Client{},
			},
			VdcResourceType,
			errDoCrudRequestsBuilder(readOperation,
				resourceName,
				errEmptyResp),
			map[string]interface{}{},
		},
		{
			5,
			Error404HTTPClienterFake{},
			&API{
				Token: rightAPIToken,
				URL:   rightAPIURL,
				Meta: APIMeta{
					NonCriticalResourceList: nonCriticalResourceMetaDataList,
					CriticalResourceList:    criticalResourceMetaDataList,
					OtherResourceList:       otherResourceMetaDataList,
				},
				Client: &http.Client{},
			},
			VdcResourceType,
			ErrResourceNotExist,
			map[string]interface{}{},
		},
		{
			6,
			HandleRespErrHTTPClienterFake{},
			&API{
				Token: rightAPIToken,
				URL:   rightAPIURL,
				Meta: APIMeta{
					NonCriticalResourceList: nonCriticalResourceMetaDataList,
					CriticalResourceList:    criticalResourceMetaDataList,
					OtherResourceList:       otherResourceMetaDataList,
				},
				Client: &http.Client{},
			},
			VMResourceType,
			errHandleResponse,
			map[string]interface{}{},
		},
		{
			7,
			VdcReadSuccessHTTPClienterFake{},
			&API{
				Token:  rightAPIToken,
				URL:    rightAPIURL,
				Meta:   APIMeta{},
				Client: &http.Client{},
			},
			VdcResourceType,
			errResourceNotExist(RAMField),
			map[string]interface{}{},
		},
	}
	APIImplementerer := AirDrumResourcesAPI{}
	fakeClientTooler := ClientTooler{}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		respCreationMap, err := APIImplementerer.ReadResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			testCase.TcAPI)
		diffs := cmp.Diff(respCreationMap, testCase.ReadResource)
		switch {
		case err == nil || testCase.ReadError == nil:
			if !((err == nil) && (testCase.ReadError == nil)) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.ReadError)
			} else {
				if diffs != "" {
					t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
						testCase.ID, diffs)
				}
			}
		case err != nil && testCase.ReadError != nil:
			if cmp.Diff(respCreationMap, map[string]interface{}{}) != "" {
				t.Errorf("\n\nTC %d : Wrong created resource map,"+
					" it should be nil as error is not nil,"+
					"\n\rgot map: \n\r\"%s\"\n\rwant map: \n\r\"%s\"\n\r",
					testCase.ID, respCreationMap, testCase.ReadResource)
			}
			if err.Error() != testCase.ReadError.Error() {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"",
					testCase.ID, err.Error(), testCase.ReadError.Error())
			}
		case diffs != "":
			t.Errorf("\n\nTC %d : Wrong resource read resource map (-got +want) :\n%s",
				testCase.ID, diffs)
		}
	}
}

func TestUpdateResource(t *testing.T) {
	resourceName := "Unit test resource update"
	testCases := []struct {
		ID           int
		TcClienter   Clienter
		ResourceType string
		UpdateErr    error
	}{
		{
			1,
			ResourceUpdateSuccessHTTPClienterFake,
			VMResourceType,
			nil,
		},
		{
			2,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
		},
		{
			3,
			ResourceUpdateFailureHTTPClienterFake,
			VdcResourceType,
			errDoCrudRequestsBuilder(updateOperation,
				resourceName,
				errEmptyResp),
		},
	}
	APIImplementerer := AirDrumResourcesAPI{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeTemplatesTooler := TemplatesTooler{
		TemplatesTools: TemplateTemplater{},
	}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		err := APIImplementerer.UpdateResource(d,
			&fakeClientTooler,
			&fakeTemplatesTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		switch {
		case err == nil || testCase.UpdateErr == nil:
			if !(err == nil && testCase.UpdateErr == nil) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.UpdateErr)
			}
		case err.Error() != testCase.UpdateErr.Error():
			t.Errorf("\n\nTC %d : resource read error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.UpdateErr.Error())
		}
	}
}

func TestDeleteResource(t *testing.T) {
	resourceName := "Unit test resource deletion"
	testCases := []struct {
		ID           int
		TcClienter   Clienter
		ResourceType string
		DeleteErr    error
	}{
		{
			1,
			ResourceDeletionSuccessHTTPClienterFake,
			VMResourceType,
			nil,
		},
		{
			2,
			HTTPClienterDummy{},
			wrongResourceType,
			errWrongResourceTypeBuilder(wrongResourceType),
		},
		{
			3,
			ResourceDeletionFailureHTTPClienterFake{},
			VdcResourceType,
			errDoCrudRequestsBuilder(deleteOperation,
				resourceName,
				errEmptyResp),
		},
	}
	APIImplementerer := AirDrumResourcesAPI{}
	sewan := &API{Token: "42", URL: "42", Client: &http.Client{}}
	fakeClientTooler := ClientTooler{}
	fakeResourceTooler := ResourceTooler{
		Resource: ResourceResourceer{},
	}
	for _, testCase := range testCases {
		resourceResponse := resource(testCase.ResourceType)
		d := resourceResponse.TestResourceData()
		d.SetId("UnitTest resource1")
		d.Set(NameField, resourceName)
		fakeClientTooler.Client = testCase.TcClienter
		err := APIImplementerer.DeleteResource(d,
			&fakeClientTooler,
			&fakeResourceTooler,
			testCase.ResourceType,
			sewan)
		switch {
		case err == nil || testCase.DeleteErr == nil:
			if !(err == nil && testCase.DeleteErr == nil) {
				t.Errorf("\n\nTC %d : resource read error was incorrect,"+
					"\n\rgot: \"%s\"\n\rwant: \"%s\"", testCase.ID, err, testCase.DeleteErr)
			}
		case err.Error() != testCase.DeleteErr.Error():
			t.Errorf("\n\nTC %d : resource read error was incorrect,"+
				"\n\rgot: \"%s\"\n\rwant: \"%s\"",
				testCase.ID, err.Error(), testCase.DeleteErr.Error())
		}
	}
}
