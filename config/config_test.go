package config_test

import (
	. "github.com/jh-bate/fantail-bot/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	var (
		myConfig map[string]interface{}
	)

	Describe("When none exists", func() {
		It("should not load anything", func() {
			Load(&myConfig, "does_not_exist")
			Expect(myConfig).To(BeNil())
		})
	})

	Describe("When it exists", func() {
		It("should load the asked for config", func() {
			Load(&myConfig, "chat.json")
			Expect(myConfig).To(Not(BeNil()))
			Expect(myConfig["QandA"]).To(Not(BeNil()))
		})
	})

})
