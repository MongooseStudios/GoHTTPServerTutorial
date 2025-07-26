package users

import (
	"errors"
	"net/mail"
	"reflect"
	"testing"
)

func TestAddUser(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	if len(testManager.users) != 1 {
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 1, len(testManager.users))
		if len(testManager.users) < 1 {
			t.Fatal()
		}
	}

	expectedUser := User{
		FirstName: testFirstName,
		LastName:  testLastName,
		Email:     *testEmail,
	}

	foundUser := testManager.users[0]

	if !reflect.DeepEqual(expectedUser, foundUser) {
		t.Errorf("added user data is not correct\nwanted: %+v\ngot: %+v\n", expectedUser, foundUser)
	}
}

func TestAddUserInvalidEmail(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail := "foobar"

	err := testManager.AddUser(testFirstName, testLastName, testEmail)
	if err == nil {
		t.Error("no error returned for invalid email")
	} else {
		expectedErr := "invalid email: foobar"
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err)
		}
	}

	if len(testManager.users) > 0 {
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 0, len(testManager.users))
	}
}

func TestAddUserEmptyFirstName(t *testing.T) {
	testManager := NewManager()

	testFirstName := ""
	testLastName := "Userman"
	// TODO: make this email valid
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Error("no error returned for invalid email")
	} else {
		expectedErr := "invalid first name: \"\""
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err)
		}
	}

	if len(testManager.users) > 0 {
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 0, len(testManager.users))
	}
}

func TestAddUserEmptyLastName(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := ""
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Error("no error returned for invalid email")
	} else {
		expectedErr := "invalid last name: \"\""
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err)
		}
	}

	if len(testManager.users) > 0 {
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 0, len(testManager.users))
	}
}

func TestAddUserDuplicateName(t *testing.T) {
	testManager := NewManager()

	testFirstName := "Test"
	testLastName := "Userman"
	testEmail, err := mail.ParseAddress("foo@bar.com")
	if err != nil {
		t.Fatalf("error parsing test email address: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err != nil {
		t.Fatalf("error creating user: %v", err)
	}

	err = testManager.AddUser(testFirstName, testLastName, testEmail.String())
	if err == nil {
		t.Error("no error returned for duplicate user", err)
	} else {
		expectedErr := "user with this name already exists"
		if err.Error() != expectedErr {
			t.Errorf("bad error text, wanted: %s, got: %s", expectedErr, err)
		}
	}

	if len(testManager.users) != 1 {
		t.Errorf("bad test manager user count, wanted: %d, got: %d", 1, len(testManager.users))
	}
}

func TestGetUserByName(t *testing.T) {
	// TODO: replace this with NewManager
	testManager := NewManager()

	// create some test users with overlapping names
	err := testManager.AddUser("foo", "bar", "f.bar@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}
	err = testManager.AddUser("bar", "baz", "bbaz@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}
	err = testManager.AddUser("foo", "baz", "fbaz@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}
	err = testManager.AddUser("baz", "foo", "bazf@example.com")
	if err != nil {
		t.Fatalf("error adding test user: %v", err)
	}

	tests := map[string]struct {
		first       string
		last        string
		expected    *User
		expectedErr error
	}{
		"simple lookup": {
			first:       "foo",
			last:        "bar",
			expected:    &testManager.users[0],
			expectedErr: nil,
		},
		"last element lookup": {
			first:       "baz",
			last:        "foo",
			expected:    &testManager.users[3],
			expectedErr: nil,
		},
		"similar name returns correct user": {
			first:       "foo",
			last:        "baz",
			expected:    &testManager.users[2],
			expectedErr: nil,
		},
		"no match lookup": {
			first:       "quux",
			last:        "quuz",
			expected:    nil,
			expectedErr: ErrNoResultsFound,
		},
		"partial match lookup": {
			first:       "foo",
			last:        "foo",
			expected:    nil,
			expectedErr: ErrNoResultsFound,
		},
		"empty first name": {
			first:       "",
			last:        "baz",
			expected:    nil,
			expectedErr: ErrNoResultsFound,
		},
		"empty last name": {
			first:       "foo",
			last:        "",
			expected:    nil,
			expectedErr: ErrNoResultsFound,
		},
	}

	for name, test := range tests {
		result, err := testManager.GetUserByName(test.first, test.last)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("%s: invalid result\ngot: %+v,\nwanted: %+v", name, result, test.expected)
		}
		if !errors.Is(err, test.expectedErr) {
			t.Errorf("%s: invalid error result\ngot: %v\nwanted: %v", name, result, test.expectedErr)
		}
	}
}
