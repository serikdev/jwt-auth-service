package usecases

import "fmt"

type MockEmailSender struct{}

func (m *MockEmailSender) SendWarning(userGUID, newIP string) error {
	fmt.Printf("Sending email to user %s about IP change to %s\n", userGUID, newIP)
	return nil
}
