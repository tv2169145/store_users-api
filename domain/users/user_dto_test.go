package users_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tv2169145/store_users-api/domain/users"
)

var _ = Describe("user_dto", func() {
	var (
		usr *users.User

		setupData = func() {
			usr = &users.User{
				FirstName: "aaaa",
				LastName:  "bbbb",
				Email:     "cccc@gmail.com",
				Password: "1234",
			}
		}
	)

	BeforeEach(func() {
		setupData()
	})

	Describe("Validate", func() {
		It("should fail because email is empty", func() {
			usr.Email = ""
			err := usr.Validate()
			Ω(err).NotTo(BeNil())
			Ω(err.Message).To(Equal("invalid email address"))
		})

		It("success", func() {
			err := usr.Validate()
			Ω(err).To(BeNil())
		})
	})
})
