package gstats

import (
	"errors"
	"math"
	"net"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func helper_GetDriftFromTrace(readStr string, statsdPrefix string, traceIdentifier string) float64 {
	// match something like: "test.testing trace:1139292|ms"
	re := regexp.MustCompile(`(.*)\.(.*):([0-9]*)\|(.*)`)
	match := re.FindStringSubmatch(readStr)
	Expect(len(match)).To(Equal(5))
	Expect(match[1]).To(Equal(statsdPrefix))
	Expect(match[2]).To(Equal(traceIdentifier))
	Expect(match[4]).To(Equal("ms"))
	traceTime, err := strconv.ParseInt(match[3], 10, 64)
	Expect(err).NotTo(HaveOccurred())
	driftNanoseconds := math.Abs(float64(traceTime - int64(time.Millisecond)))
	return driftNanoseconds
}

func TestStatsd(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	g.Describe("Sanity", func() {
		g.It("tests should function in a sane environment", func() {
			Expect(1).To(Equal(1))
		})
	})
	g.Describe("Normalization (vendored import)", func() {
		g.It("should correctly cammel-case a simple happy-path error text", func() {
			err := errors.New("cammel case")
			normalizedErrorString := normalize(err)
			Expect(normalizedErrorString).To(Equal("CammelCase"))
		})
		g.It("should remove non-alpha-numeric characters when normalizing error text", func() {
			err := errors.New("This is a full sentence. Does it have punctuation?")
			normalizedErrorString := normalize(err)
			Expect(normalizedErrorString).To(Equal("ThisIsAFullSentenceDoesItHavePunctuation"))

			err = errors.New("Here is some !@#$%^&*()_+= Perl")
			normalizedErrorString = normalize(err)
			Expect(normalizedErrorString).To(Equal("HereIsSomePerl"))
		})
		g.It("should strip non-western UTF-8 characters", func() {
			err := errors.New("Here is some !@#$%^&*()_+= Perl Hello, 世界")
			normalizedErrorString := normalize(err)
			Expect(normalizedErrorString).To(Equal("HereIsSomePerlHello"))
		})
	})
	g.Describe("Environment", func() {
		g.BeforeEach(func() {
			os.Clearenv()
		})
		g.It("should fail fast without STATSD_ADDRESS defined", func() {
			os.Setenv("STATSD_PREFIX", "test")
			_, err := CreateStatsdClient()
			Expect(err).To(HaveOccurred())
		})
		g.It("should fail fast without STATSD_PREFIX defined", func() {
			os.Setenv("STATSD_ADDRESS", "127.0.0.1:31337")
			_, err := CreateStatsdClient()
			Expect(err).To(HaveOccurred())
		})
	})
	g.Describe("statsd", func() {
		buf := make([]byte, 1024)
		var addr *net.UDPAddr
		var sock *net.UDPConn
		var err error
		g.BeforeEach(func() {
			addr = nil
			sock = nil
			addr, err = net.ResolveUDPAddr("udp", "127.0.0.1:31337")
			sock, err = net.ListenUDP("udp", addr)
			os.Setenv("STATSD_ADDRESS", "127.0.0.1:31337")
			os.Setenv("STATSD_PREFIX", "test")
		})
		g.AfterEach(func() {
			if sock != nil {
				sock.Close()
			}
		})
		g.It("should correctly initialize", func() {
			client, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			Expect(client).NotTo(Equal(nil))
		})
		g.It("should send increment", func() {
			client, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 5; i++ {
				err = client.Inc("teststat")
				Expect(err).NotTo(HaveOccurred())
				readLength, _, err := sock.ReadFromUDP(buf)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(buf[:readLength])).To(Equal("test.teststat.count:1|c"))
			}
		})
		g.It("should send accurate timing", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			// Note: usually, we'd handle tracing in a call that looks like this:
			// defer stats.End(Trace("testing trace"))
			// but because we're trying to get results from the UDP listener
			// we've stubbed in instead of a real statsd service, we need more
			// method body after the "do some work for a while" part of our function.
			// In order to accomodate this, we need to break things up in a way
			// that won't be used in real code very often
			traceIdentifier, timestamp, incrementBy := Trace("testing trace")
			// do some work for a while
			time.Sleep(1 * time.Millisecond)
			stats.End(traceIdentifier, timestamp, incrementBy)
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr := string(buf[:readLength])
			driftNanoseconds := helper_GetDriftFromTrace(readStr, "test", "testing trace")
			Expect(driftNanoseconds).Should(BeNumerically("<", 10*time.Millisecond))
		})
		g.It("should send accurate timing even with BufferedEnd", func() {
			stats, err := _CreateStatsdClient(100 * time.Millisecond)
			Expect(err).NotTo(HaveOccurred())
			// Note: usually, we'd handle tracing in a call that looks like this:
			// defer stats.End(Trace("testing trace"))
			// but because we're trying to get results from the UDP listener
			// we've stubbed in instead of a real statsd service, we need more
			// method body after the "do some work for a while" part of our function.
			// In order to accomodate this, we need to break things up in a way
			// that won't be used in real code very often
			traceIdentifier, timestamp, incrementBy := Trace("testing trace")
			// do some work for a while
			time.Sleep(1 * time.Millisecond)
			stats.BufferedEnd(traceIdentifier, timestamp, incrementBy)
			time.Sleep(105 * time.Millisecond)
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr := string(buf[:readLength])
			driftNanoseconds := helper_GetDriftFromTrace(readStr, "test", "testing trace")
			Expect(driftNanoseconds).Should(BeNumerically("<", 10*time.Millisecond))
		})
		g.It("should trace and count", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			traceIdentifier, timestamp, incrementBy := TraceAndIncrement("testing trace")
			// do some work for a while
			time.Sleep(1 * time.Millisecond)
			stats.End(traceIdentifier, timestamp, incrementBy)

			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr := string(buf[:readLength])
			Expect(string(buf[:readLength])).To(Equal("test.testing trace.count:1|c"))

			readLength, _, err = sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr = string(buf[:readLength])
			driftNanoseconds := helper_GetDriftFromTrace(readStr, "test", "testing trace")
			Expect(driftNanoseconds).Should(BeNumerically("<", 10*time.Millisecond))
		})
		g.It("should trace and count even with BufferedEnd", func() {
			stats, err := _CreateStatsdClient(100 * time.Millisecond)
			Expect(err).NotTo(HaveOccurred())
			traceIdentifier, timestamp, incrementBy := TraceAndIncrement("testing trace")
			// do some work for a while
			time.Sleep(1 * time.Millisecond)
			stats.BufferedEnd(traceIdentifier, timestamp, incrementBy)
			time.Sleep(105 * time.Millisecond)

			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr := string(buf[:readLength])
			driftNanoseconds := helper_GetDriftFromTrace(readStr, "test", "testing trace")
			Expect(driftNanoseconds).Should(BeNumerically("<", 10*time.Millisecond))

			readLength, _, err = sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr = string(buf[:readLength])
			Expect(string(buf[:readLength])).To(Equal("test.testing trace.count:1|c"))
		})
		g.It("should stat normalized error text", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			err = errors.New("Custom !@#$ error &((@*# including. Punctuation & Perl")
			err = stats.IncErr("myStat", err)
			Expect(err).NotTo(HaveOccurred())

			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(buf[:readLength])).To(Equal("test.myStat.CustomErrorIncludingPunctuationPerl.count:1|c"))
		})
		g.It("should increment by the requested amount", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			err = stats.IncrementBy("testing IncrementBy", int64(20))
			Expect(err).NotTo(HaveOccurred())
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(readLength).To(BeNumerically(">", 0))
			Expect(string(buf[:readLength])).To(Equal("test.testing IncrementBy:20|c"))
		})
		g.It("should buffer incrementBy calls until it gets to 100", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 5; i++ {
				err = stats.BufferedIncrementBy("testing IncrementBy", int64(20))
				Expect(err).NotTo(HaveOccurred())
			}
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(readLength).To(BeNumerically(">", 0))
			Expect(string(buf[:readLength])).To(Equal("test.testing IncrementBy:100|c"))
		})
		g.It("should buffer incrementBy calls until it gets to 100 even if it does not get there in even numbers", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 6; i++ {
				err = stats.BufferedIncrementBy("testing IncrementBy", int64(19))
				Expect(err).NotTo(HaveOccurred())
			}
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(readLength).To(BeNumerically(">", 0))
			Expect(string(buf[:readLength])).To(Equal("test.testing IncrementBy:114|c"))
		})
		g.It("should buffer separate stats seperately", func() {
			stats, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 4; i++ {
				err = stats.BufferedIncrementBy("testing IncrementBy 1", int64(20))
				err = stats.BufferedIncrementBy("testing IncrementBy 2", int64(20))
				Expect(err).NotTo(HaveOccurred())
			}
			err = stats.BufferedIncrementBy("testing IncrementBy 1", int64(20))
			Expect(err).NotTo(HaveOccurred())
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(readLength).To(BeNumerically(">", 0))
			Expect(string(buf[:readLength])).To(Equal("test.testing IncrementBy 1:100|c"))
			err = stats.BufferedIncrementBy("testing IncrementBy 2", int64(20))
			Expect(err).NotTo(HaveOccurred())
			readLength, _, err = sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			Expect(readLength).To(BeNumerically(">", 0))
			Expect(string(buf[:readLength])).To(Equal("test.testing IncrementBy 2:100|c"))
		})
		g.It("should flush stats periodically", func() {
			stats, err := _CreateStatsdClient(100 * time.Millisecond)
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 5; i++ {
				err = stats.BufferedIncrementBy("testing IncrementBy", int64(19))
				Expect(err).NotTo(HaveOccurred())
			}
			time.Sleep(105 * time.Millisecond)
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(readLength).To(BeNumerically(">", 0))
			Expect(string(buf[:readLength])).To(Equal("test.testing IncrementBy:95|c"))
		})
	})
}
