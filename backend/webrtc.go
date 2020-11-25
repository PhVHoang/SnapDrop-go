// credit to https://github.com/poi5305/go-yuv2webRTC/blob/master/webrtc/webrtc.go
package backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
	vpxEncoder "github.com/poi5305/go-yuv2webRTC/vpx-encoder"
)

var config = webrtc.Configuration{
	BundlePolicy: webrtc.BundlePolicyMaxBundle,
	ICEServers: []webrtc.ICEServer{
		{URLs: []string{"stun:stun.l.google.com:19302"}},
	},
}

// Create a new WebRTC
func NewWebRTC() *WebRTC {
	w := &WebRTC{
		ImageChannel: make(chan []byte, 2),
	}
	return w
}

type WebRTC struct {
	connection   *webrtc.PeerConnection
	encoder      *vpxEncoder.VpxEncoder
	vp8Track     *webrtc.Track
	isConnected  bool
	ImageChannel chan []byte
}

func (w *WebRTC) StopClient() {
	fmt.Println("======Stop Client======")
	w.isConnected = false
	if w.encoder != nil {
		w.encoder.Release()
	}
	if w.connection != nil {
		w.connection.Close()
	}
	w.connection = nil
}

func (w *WebRTC) StartStreaming(vp8Track *webrtc.Track) {
	fmt.Println("Start Streaming")
	// send screenshots
	go func() {
		for w.isConnected {
			yuv := <-w.ImageChannel
			if len(w.encoder.Input) < cap(w.encoder.Input) {
				w.encoder.Input <- yuv
			}
		}
	}()

	// receive frame buffer
	go func() {
		for i := 0; w.isConnected; i++ {
			bs := <-w.encoder.Output
			if i%10 == 0 {
				fmt.Println("On frame ", len(bs), i)
			}
			w.vp8Track.WriteSample(media.Sample{Data: bs, Sample: 1})
		}
	}()
}

func getV8PlayLoadType(sdp string) (int, int) {
	lines := string.Split(sdp, "\n")
	payloadType := 96
	clockRate := 90000
	for _, line := range lines {
		if strings.Contains(line, "a=rtpmap:") && strings.Contains(line, "VP8/") {
			fmt.Sscanf(line, "a=rtpmap:%d VP8/%d", &payloadType, &clockRate)
			break
		}
	}
	return payloadType, clockRate
}

func Decode(in string, obj interface{}) {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		panic(err)
	}
}

func (w *WebRTC) IsConnected() bool {
	return w.isConnected
}

func isPlanB(sdp string) bool {
	lines := strings.Split(sdp, "\n")
	for _, line := range lines {
		if strings.Contains(line, "mid:video") || strings.Contains(line, "mid:audio") ||
			strings.Contains(line, "mid:data") {
			return true
		}
	}
	return false
}

func Encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (w *WebRTC) StartClient(remoteSession string, width, height int) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			w.StopClient()
		}
	}()

	// Reset client
	if w.isConnected {
		w.StopClient()
		time.Sleep(2 * time.Second)
	}

	offer := webrtc.SessionDescription{}
	Decode(remoteSession, &offer)

	isPlanB := isPlanB(offer.SDP)

	// Safari VP8 Payload default is not 96
	payloadType, clockRate := getV8PlayLoadType(offer.SDP)
	m := webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.NewRTPVP8Codec(uint8(payloadType), uint32(clockRate)))
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	if isPlanB {
		config.SDPSemantics = webrtc.SDPSemanticsPlanB
	} else {
		config.SDPSemantics = webrtc.SDPSemanticsUnifiedPlan
	}

	encoder, err := vpxEncoder.NewVpxEncoder(width, height, 20, 1200, 5)
	if err != nil {
		return "", err
	}
	w.encoder = encoder

	fmt.Println("======Start Client======")
	conn, err := api.NewPeerConnection(config)
	if err != nil {
		return "", err
	}

	vp8Track, err := w.connection.NewTrack(uint8(payloadType), rand.Uint32(), "video", "robotmon")
	if err != nil {
		fmt.Println("Error: new xRobotmonScreen vp8 Track", err)
		w.StopClient()
		return "", err
	}
	_, err = w.connection.AddTrack(vp8Track)

}
