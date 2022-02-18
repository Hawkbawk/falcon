package darwin

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Darwin", func() {
	var argList []string
	var err error
	var cmdRunner = func(cmd string) error {
		argList = append(argList, cmd)
		return err
	}

	BeforeEach(func() {
		argList = make([]string, 0)
		err = nil
	})


	Describe("addDockerResolver", func() {
		It("tries to run the addResolver command", func() {
			Expect(addDockerResolver(cmdRunner)).Should(Succeed())
			Expect(argList[0]).Should(Equal(addResolverCmd))
		})

		Describe("an error occurs", func() {
			BeforeEach(func() {
				err = fmt.Errorf("didn't work :(")
			})

			It("returns that error", func() {
				Expect(addDockerResolver(cmdRunner)).Should(Equal(err))
			})
		})
	})

	Describe("addLoopbackAddress", func() {
		It("tries to run the addLoopbackAddress command", func() {
			Expect(addLoopbackAddress(cmdRunner)).Should(Succeed())
			Expect(argList[0]).Should(Equal(addLoopbackAddressCmd))
		})

		Describe("an error occurs", func() {
			BeforeEach(func() {
				err = fmt.Errorf("didn't work :(")
			})

			It("returns that error", func() {
				Expect(addLoopbackAddress(cmdRunner)).Should(Equal(err))
			})
		})
	})

	Describe("removeDockerResolver", func() {
		It("tries to run the removeResolver command", func() {
			Expect(removeDockerResolver(cmdRunner)).Should(Succeed())
			Expect(argList[0]).Should(Equal(removeResolverCmd))
		})

		Describe("an error occurs", func() {
			BeforeEach(func() {
				err = fmt.Errorf("didn't work :(")
			})

			It("returns that error", func() {
				Expect(removeDockerResolver(cmdRunner)).Should(Equal(err))
			})
		})
	})

	Describe("removeLoopbackAddress", func() {
		It("tries to run the removeLoopbackAddress command", func() {
			Expect(removeLoopbackAddress(cmdRunner)).Should(Succeed())
			Expect(argList[0]).Should(Equal(removeLoopbackAddressCmd))
		})

		Describe("an error occurs", func() {
			BeforeEach(func() {
				err = fmt.Errorf("didn't work :(")
			})

			It("returns that error", func() {
				Expect(removeLoopbackAddress(cmdRunner)).Should(Equal(err))
			})
		})
	})
})
