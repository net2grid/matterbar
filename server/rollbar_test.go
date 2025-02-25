package main

import (
	"encoding/json"
	"testing"
)

func TestRollbarCustomerTimestamp(t *testing.T) {
	for name, test := range map[string]struct {
		TestFile                  string
		ExpectedCustomerTimestamp string
	}{
		"ok - customer timestamp int": {
			TestFile:                  "new_item.json",
			ExpectedCustomerTimestamp: "1541827726",
		},
		"ok - java customer timestamp decimal": {
			TestFile:                  "new_item_java.json",
			ExpectedCustomerTimestamp: "1545202894.312",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var rollbar Rollbar
			data := loadJsonFile(t, test.TestFile)
			err := json.Unmarshal([]byte(data), &rollbar)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	for name, test := range map[string]struct {
		TestFile       string
		ExpectedReturn string
	}{
		"valid - deploy": {
			TestFile:       "deploy.json",
			ExpectedReturn: "",
		},
		"valid - exp repeat item": {
			TestFile:       "exp_repeat_item.json",
			ExpectedReturn: "",
		},
		"valid - high velocity item": {
			TestFile:       "item_velocity.json",
			ExpectedReturn: "",
		},
		"valid - new item": {
			TestFile:       "new_item.json",
			ExpectedReturn: "",
		},
		"valid - every occurrence": {
			TestFile:       "occurrence.json",
			ExpectedReturn: "",
		},
		"valid - reactivated item": {
			TestFile:       "reactivated_item.json",
			ExpectedReturn: "",
		},
		"valid - reopened item": {
			TestFile:       "reopened_item.json",
			ExpectedReturn: "",
		},
		"valid - resolved item": {
			TestFile:       "resolved_item.json",
			ExpectedReturn: "",
		},
		"valid - test": {
			TestFile:       "test.json",
			ExpectedReturn: "",
		},
		"invalid - empty data": {
			TestFile:       "invalid/empty.json",
			ExpectedReturn: "Missing rollbar.event_name",
		},
		"invalid - unsupported event name": {
			TestFile:       "invalid/unsupported_event_name.json",
			ExpectedReturn: "Unsupported rollbar.event_name",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var rollbar Rollbar
			data := loadJsonFile(t, test.TestFile)
			err := json.Unmarshal([]byte(data), &rollbar)
			if err != nil {
				t.Fatal(err)
			}
			actualReturn := rollbar.isValid()
			if actualReturn != test.ExpectedReturn {
				t.Errorf("Expected: %s\nActual: %s", test.ExpectedReturn, actualReturn)
			}
		})
	}
}

func TestEventNameToTitle(t *testing.T) {
	for name, test := range map[string]struct {
		TestFile      string
		ExpectedTitle string
	}{
		"ok - new item": {
			TestFile:      "new_item.json",
			ExpectedTitle: "New Error",
		},
		"ok - reactivated item": {
			TestFile:      "reactivated_item.json",
			ExpectedTitle: "Reactivated Error",
		},
		"ok - exp repeat item": {
			TestFile:      "exp_repeat_item.json",
			ExpectedTitle: "10th Error",
		},
		"ok - reopened": {
			TestFile:      "reopened_item.json",
			ExpectedTitle: "Reopened Error",
		},
		"ok - resolved": {
			TestFile:      "resolved_item.json",
			ExpectedTitle: "Resolved Error",
		},
		"ok - occurrence": {
			TestFile:      "occurrence.json",
			ExpectedTitle: "Occurrence - Error",
		},
		"ok - high velocity": {
			TestFile:      "item_velocity.json",
			ExpectedTitle: "5 occurrences in 5 minutes",
		},
		"ok - deploy": {
			TestFile:      "deploy.json",
			ExpectedTitle: "Deploy",
		},
		"ok - test": {
			TestFile:      "test.json",
			ExpectedTitle: "",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var rollbar Rollbar
			data := loadJsonFile(t, test.TestFile)
			err := json.Unmarshal([]byte(data), &rollbar)
			if err != nil {
				t.Fatal(err)
			}
			actualTitle := rollbar.eventNameToTitle()
			if actualTitle != test.ExpectedTitle {
				t.Errorf("Expected: %s\nActual: %s", test.ExpectedTitle, actualTitle)
			}
		})
	}
}

func TestEventText(t *testing.T) {
	for name, test := range map[string]struct {
		TestFile     string
		ExpectedText string
	}{
		"ok - new item exception data under last_occurrence": {
			TestFile:     "new_item.json",
			ExpectedText: "TypeError: unsupported operand type(s) for +=: 'int' and 'str'",
		},
		"ok - new item log message no traceback": {
			TestFile:     "new_item_log_message.json",
			ExpectedText: "User 8563892 is missing permissions",
		},
		"ok - occurrence exception data under occurrence": {
			TestFile:     "occurrence.json",
			ExpectedText: "TypeError: 'NoneType' object has no attribute '__getitem__'",
		},
		"ok - high velocity missing occurrence data": {
			TestFile:     "item_velocity.json",
			ExpectedText: "",
		},
		"ok - reactivated item ios crash report": {
			TestFile:     "reactivated_item_ios_crash_report.json",
			ExpectedText: "iOS crash report",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var rollbar Rollbar
			data := loadJsonFile(t, test.TestFile)
			err := json.Unmarshal([]byte(data), &rollbar)
			if err != nil {
				t.Fatal(err)
			}
			actualText := rollbar.eventText()
			if test.ExpectedText != actualText {
				t.Errorf("Expected: %s\nActual: %s", test.ExpectedText, actualText)
			}
		})
	}
}

func TestDeployUser(t *testing.T) {
	for name, test := range map[string]struct {
		TestFile         string
		ExpectedUsername string
	}{
		"deploy with no local username": {
			TestFile:         "deploy_no_username.json",
			ExpectedUsername: "unknown user",
		},
		"deploy with local username": {
			TestFile:         "deploy.json",
			ExpectedUsername: "dliu",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var rollbar Rollbar
			data := loadJsonFile(t, test.TestFile)
			err := json.Unmarshal([]byte(data), &rollbar)
			if err != nil {
				t.Fatal(err)
			}
			actualUsername := rollbar.deployUser()
			if test.ExpectedUsername != actualUsername {
				t.Errorf("Expected: %s\nActual: %s", test.ExpectedUsername, actualUsername)
			}
		})
	}
}

func TestDeployDateTime(t *testing.T) {
	for name, test := range map[string]struct {
		TestFile         string
		ExpectedDateTime string
	}{
		"deploy": {
			TestFile:         "deploy.json",
			ExpectedDateTime: "2019-10-20 12:45:58 PDT-0700",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var rollbar Rollbar
			data := loadJsonFile(t, test.TestFile)
			err := json.Unmarshal([]byte(data), &rollbar)
			if err != nil {
				t.Fatal(err)
			}
			actualDateTime := rollbar.deployDateTime()
			if test.ExpectedDateTime != actualDateTime {
				t.Errorf("Expected: %s\nActual: %s", test.ExpectedDateTime, actualDateTime)
			}
		})
	}
}
