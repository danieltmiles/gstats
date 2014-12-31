package gstats

import (
	"errors"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestMocks(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	g.Describe("MockStatser", func() {
		var stats MockStatser
		g.BeforeEach(func() {
			stats = MockStatser{}
		})
		g.It("should record calls to Inc", func() {
			stats.Inc("testing inc")
			Expect(len(stats.CallsToInc)).To(Equal(1))
			Expect(stats.CallsToInc[0].IncVal).To(Equal("testing inc"))
		})
		g.It("should record calls to IncErr", func() {
			incErr := errors.New("this is an error to test things out")
			stats.IncErr("incing error", incErr)
			Expect(len(stats.CallsToIncErr)).To(Equal(1))
			Expect(stats.CallsToIncErr[0].IncVal).To(Equal("incing error"))
			Expect(stats.CallsToIncErr[0].Err).To(Equal(incErr))
		})
		g.It("should record calls to End", func() {
			endTime := time.Now()
			stats.End("calling end", endTime, 5)
			Expect(len(stats.CallsToEnd)).To(Equal(1))
			Expect(stats.CallsToEnd[0].Str).To(Equal("calling end"))
			Expect(stats.CallsToEnd[0].Tim).To(Equal(endTime))
			Expect(stats.CallsToEnd[0].Num).To(Equal(int64(5)))
		})
		g.It("should record calls to BufferedEnd", func() {
			endTime := time.Now()
			stats.BufferedEnd("calling end", endTime, 5)
			Expect(len(stats.CallsToBufferedEnd)).To(Equal(1))
			Expect(stats.CallsToBufferedEnd[0].Str).To(Equal("calling end"))
			Expect(stats.CallsToBufferedEnd[0].Tim).To(Equal(endTime))
			Expect(stats.CallsToBufferedEnd[0].Num).To(Equal(int64(5)))
		})
		g.It("should record calls to IncrementBy", func() {
			stats.IncrementBy("some stat", 20)
			Expect(len(stats.CallsToIncrementBy)).To(Equal(1))
			Expect(stats.CallsToIncrementBy[0].Str).To(Equal("some stat"))
			Expect(stats.CallsToIncrementBy[0].Num).To(Equal(int64(20)))
		})
		g.It("should record calls to BufferedIncrementBy", func() {
			stats.BufferedIncrementBy("some stat", 20)
			Expect(len(stats.CallsToBufferedIncrementBy)).To(Equal(1))
			Expect(stats.CallsToBufferedIncrementBy[0].Str).To(Equal("some stat"))
			Expect(stats.CallsToBufferedIncrementBy[0].Num).To(Equal(int64(20)))
		})
		g.It("should record calls to Gauge", func() {
			stats.Gauge("some stat", 38)
			Expect(len(stats.CallsToGauge)).To(Equal(1))
			Expect(stats.CallsToGauge[0].Str).To(Equal("some stat"))
			Expect(stats.CallsToGauge[0].Num).To(Equal(int64(38)))
		})
	})
}
