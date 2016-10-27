package xdgdirs_test

import (
	"github.com/redforks/hal"
	"github.com/redforks/testing/reset"
	. "github.com/redforks/xdgdirs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Xdgdirs", func() {

	var (
		envs map[string]string

		Getenv = func(key string) string {
			return envs[key]
		}

		homeDir = ""
	)

	BeforeEach(func() {
		reset.Enable()
		hal.Getenv = Getenv
		homeDir = "/user/foo"
		envs = map[string]string{
			"HOME": homeDir,
		}
	})

	AfterEach(func() {
		reset.Disable()
	})

	It("Home", func() {
		Ω(Home()).Should(Equal(homeDir))
	})

	Context("DataHome", func() {

		It("environment defined", func() {
			envs["XDG_DATA_HOME"] = "/usr/share/data"
			Ω(DataHome()).Should(Equal("/usr/share/data"))
		})

		It("environment not defined", func() {
			Ω(DataHome()).Should(Equal("/user/foo/.local/share"))
		})
	})

	Context("ConfigHome", func() {
		It("environment defined", func() {
			envs["XDG_CONFIG_HOME"] = "/foo"
			Ω(ConfigHome()).Should(Equal("/foo"))
		})

		It("environment not defined", func() {
			Ω(ConfigHome()).Should(Equal("/user/foo/.config"))
		})
	})

	Context("CacheHome", func() {
		It("environment defined", func() {
			envs["XDG_CACHE_HOME"] = "/foo"
			Ω(CacheHome()).Should(Equal("/foo"))
		})

		It("environment not defined", func() {
			Ω(CacheHome()).Should(Equal("/user/foo/.cache"))
		})
	})

	Context("DataDirs", func() {
		It("environment not defined", func() {
			Ω(DataDirs()).Should(Equal([]string{
				"/user/foo/.local/share",
				"/usr/local/share",
				"/usr/share",
			}))
		})

		It("environment defined with one entry", func() {
			envs["XDG_DATA_DIRS"] = "/share/data"
			Ω(DataDirs()).Should(Equal([]string{
				"/user/foo/.local/share",
				"/share/data",
			}))
		})

		It("environment defined with multi entries", func() {
			envs["XDG_DATA_DIRS"] = "/share/data:/data:/share"
			Ω(DataDirs()).Should(Equal([]string{
				"/user/foo/.local/share",
				"/share/data",
				"/data",
				"/share",
			}))
		})
	})

	Context("ConfigDirs", func() {
		It("environment not defined", func() {
			Ω(ConfigDirs()).Should(Equal([]string{
				"/user/foo/.config",
				"/etc/xdg",
			}))
		})

		It("environment defined with multi entries", func() {
			envs["XDG_CONFIG_DIRS"] = "/share/data:/data:/share"
			Ω(ConfigDirs()).Should(Equal([]string{
				"/user/foo/.config",
				"/share/data",
				"/data",
				"/share",
			}))
		})
	})
})
