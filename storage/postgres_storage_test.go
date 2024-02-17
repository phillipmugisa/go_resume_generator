package storage

import (
	"fmt"
	"testing"
	"github.com/phillipmugisa/go_resume_generator/data"
)

type UserTestData struct {
	user data.User
	error error
}

var database_test_data = map[string]any {
	"users": []UserTestData{
		UserTestData{data.NewUser("phillip","mugisa","phillipmugisa","test@gmail.com","testpassword","+256782047612","This is a test bio","Uganda","01/01/2019",), nil},
		UserTestData{data.NewUser("alex","mark","alexmark","alexmark@gmail.com","testpassword","+256782047612","This is a test bio 2","Uganda","01/01/2020",), nil},
		UserTestData{data.NewUser("alex","mark","alexmark","alexmark@gmail.com","testpassword","+256782047612","This is a test bio 2","Uganda","01/01/2020",), "error"},
	},
}

func TestDB(t *testing.T) {
	_, err := NewPostgresStorage()
	if err != nil {
		t.Error(err)
	}


	// test user db functionality
	for key, val := range database_test_data {
		if key == "users" {
			// loop through users
			fmt.Println(val)
			// for _, data := range val {
			// 	err := s.CreateUser(data[0])
			// 	if err != data[1] {
			// 		t.Errorf("Got %v, Expeected: %v", err, data[1])
			// 	}
			// }
		}
	}
	
}
