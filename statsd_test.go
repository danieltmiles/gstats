package stats

import (
	"fmt"
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

func TestStatsd(t *testing.T) {
	g := Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	g.Describe("Sanity", func() {
		g.It("tests should function in a sane environment", func() {
			Expect(1).To(Equal(1))
		})
	})
	g.Describe("Environment", func() {
		g.BeforeEach(func() {
			os.Clearenv()
		})
		g.It("should fail fast without STATSD_ADDRESS defined", func(){
			os.Setenv("STATSD_PREFIX", "test")
			_, err := CreateStatsdClient()
			Expect(err).To(HaveOccurred())
		})
		g.It("should fail fast without STATSD_PREFIX defined", func(){
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
		g.It("should send and track increment", func() {
			client, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			for i := 0; i < 5; i++ {
				err = client.Inc("teststat")
				Expect(err).NotTo(HaveOccurred())
				readLength, _, err := sock.ReadFromUDP(buf)
				Expect(err).NotTo(HaveOccurred())
				Expect(readLength).To(Equal(17))
				Expect(string(buf[:readLength])).To(Equal(fmt.Sprintf("test.teststat:%d|c", i+1)))
			}
		})
		g.It("should send accurate timing", func() {
			client, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			err = client.StartTrace("testing trace")
			Expect(err).NotTo(HaveOccurred())
			time.Sleep(1 * time.Millisecond)
			err = client.EndTrace("testing trace")
			Expect(err).NotTo(HaveOccurred())
			readLength, _, err := sock.ReadFromUDP(buf)
			Expect(err).NotTo(HaveOccurred())
			readStr := string(buf[:readLength])
			// match something like: "test.testing trace:1139292|ms"
			re := regexp.MustCompile(`(.*)\.(.*):([0-9]*)\|(.*)`)
			match := re.FindStringSubmatch(readStr)
			Expect(len(match)).To(Equal(5))
			Expect(match[1]).To(Equal("test"))
			Expect(match[2]).To(Equal("testing trace"))
			Expect(match[4]).To(Equal("ms"))
			traceTime, err := strconv.ParseInt(match[3], 10, 64)
			Expect(err).NotTo(HaveOccurred())
			driftNanoseconds := math.Abs(float64(traceTime - int64(time.Millisecond)))
			Expect(driftNanoseconds).Should(BeNumerically("<", 10*time.Millisecond))
		})
		g.It("should provide correct error when ending a trace without starting it", func() {
			client, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			err = client.EndTrace("testing trace")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Error, must start trace \"testing trace\" before ending it"))
		})
		g.It("should clean up completed traces", func() {
			client, err := CreateStatsdClient()
			Expect(err).NotTo(HaveOccurred())
			err = client.StartTrace("testing trace")
			Expect(err).NotTo(HaveOccurred())
			err = client.EndTrace("testing trace")
			Expect(err).NotTo(HaveOccurred())
			_, foundTestingTrace := client.TimeableStats["testing trace"]
			Expect(foundTestingTrace).To(Equal(false))
		})
	})
}
