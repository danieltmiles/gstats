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
		var stats Statser
		var mock MockStatser
		g.BeforeEach(func() {
			mock = MockStatser{}
			stats = &mock
		})
		g.It("should record calls to Inc", func() {
			stats.Inc("testing inc")
			Expect(len(mock.CallsToInc)).To(Equal(1))
			Expect(mock.CallsToInc[0].IncVal).To(Equal("testing inc"))
		})
		g.It("should record calls to IncErr", func() {
			incErr := errors.New("this is an error to test things out")
			stats.IncErr("incing error", incErr)
			Expect(len(mock.CallsToIncErr)).To(Equal(1))
			Expect(mock.CallsToIncErr[0].IncVal).To(Equal("incing error"))
			Expect(mock.CallsToIncErr[0].Err).To(Equal(incErr))
		})
		g.It("should record calls to End", func() {
			endTime := time.Now()
			stats.End("calling end", endTime, 5)
			Expect(len(mock.CallsToEnd)).To(Equal(1))
			Expect(mock.CallsToEnd[0].Str).To(Equal("calling end"))
			Expect(mock.CallsToEnd[0].Tim).To(Equal(endTime))
			Expect(mock.CallsToEnd[0].Num).To(Equal(int64(5)))
		})
		g.It("should record calls to BufferedEnd", func() {
			endTime := time.Now()
			stats.BufferedEnd("calling end", endTime, 5)
			Expect(len(mock.CallsToBufferedEnd)).To(Equal(1))
			Expect(mock.CallsToBufferedEnd[0].Str).To(Equal("calling end"))
			Expect(mock.CallsToBufferedEnd[0].Tim).To(Equal(endTime))
			Expect(mock.CallsToBufferedEnd[0].Num).To(Equal(int64(5)))
		})
		g.It("should record calls to IncrementBy", func() {
			stats.IncrementBy("some stat", 20)
			Expect(len(mock.CallsToIncrementBy)).To(Equal(1))
			Expect(mock.CallsToIncrementBy[0].Str).To(Equal("some stat"))
			Expect(mock.CallsToIncrementBy[0].Num).To(Equal(int64(20)))
		})
		g.It("should record calls to BufferedIncrementBy", func() {
			stats.BufferedIncrementBy("some stat", 20)
			Expect(len(mock.CallsToBufferedIncrementBy)).To(Equal(1))
			Expect(mock.CallsToBufferedIncrementBy[0].Str).To(Equal("some stat"))
			Expect(mock.CallsToBufferedIncrementBy[0].Num).To(Equal(int64(20)))
		})
		g.It("should record calls to Gauge", func() {
			stats.Gauge("some stat", 38)
			Expect(len(mock.CallsToGauge)).To(Equal(1))
			Expect(mock.CallsToGauge[0].Str).To(Equal("some stat"))
			Expect(mock.CallsToGauge[0].Num).To(Equal(int64(38)))
		})
	})
}
