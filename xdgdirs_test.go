package xdgdirs_test

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/redforks/osutil"

	"github.com/redforks/hal"
	"github.com/redforks/testing/iotest"
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

		setXDGDataDirs = func(value string) {
			envs["XDG_DATA_DIRS"] = value
		}
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

	Context("RuntimeHome", func() {
		It("environment defined", func() {
			envs["XDG_RUNTIME_DIR"] = "/foo"
			Ω(RuntimeHome()).Should(Equal("/foo"))
		})

		It("default for root", func() {
			hal.CurrentUser = func() (*user.User, error) {
				return &user.User{
					Uid:  "0",
					Name: "root",
				}, nil
			}
			Ω(RuntimeHome()).Should(Equal("/run"))
		})

		It("default for non-root", func() {
			hal.CurrentUser = func() (*user.User, error) {
				return &user.User{
					Uid:  "100",
					Name: "user",
				}, nil
			}
			Ω(RuntimeHome()).Should(Equal("/tmp/user"))
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
			setXDGDataDirs("/share/data")
			Ω(DataDirs()).Should(Equal([]string{
				"/user/foo/.local/share",
				"/share/data",
			}))
		})

		It("environment defined with multi entries", func() {
			setXDGDataDirs("/share/data:/data:/share")
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

	Context("Resolve file path", func() {

		BeforeEach(func() {
			homeDir = iotest.NewTempTestDir().Dir()
			envs["HOME"] = homeDir
			setXDGDataDirs("")
		})

		Context("ResolveDataFile", func() {
			It("Found in 1st dir", func() {
				fooFile := filepath.Join(homeDir, ".local/share/foo")
				Ω(osutil.WriteFile(fooFile, nil, 0700, 0600)).Should(Succeed())
				Ω(ResolveDataFile("foo")).Should(Equal(fooFile))
			})

			It("Found in 2nd dir", func() {
				path2 := filepath.Join(homeDir, "data")
				fooFile := filepath.Join(path2, "foo")
				setXDGDataDirs(path2)
				Ω(osutil.WriteFile(fooFile, nil, 0700, 0600)).Should(Succeed())
				Ω(ResolveDataFile("foo")).Should(Equal(fooFile))
			})

			It("File not found", func() {
				_, err := ResolveDataFile("foo")
				Ω(err).Should(MatchError("[xdgdirs] Can not found data file: foo"))
			})

			It("Exist but not regular file", func() {
				Ω(os.MkdirAll(filepath.Join(homeDir, ".local/share/foo"), 0700)).Should(Succeed())
				_, err := ResolveDataFile("foo")
				Ω(err).ShouldNot(BeNil())
				Ω(err.Error()).Should(ContainSubstring("foo\" exist but not regular file"))
			})
		})

		It("ResolveConfigFile", func() {
			// ResolveConfigFile share the same implementation with ResolveDataFile(), so
			// no need doing detailed tests
			fooFile := filepath.Join(homeDir, ".config/foo")
			Ω(osutil.WriteFile(fooFile, nil, 0700, 0600)).Should(Succeed())
			Ω(ResolveConfigFile("foo")).Should(Equal(fooFile))
		})
	})

})
