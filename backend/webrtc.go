// credit to https://github.com/poi5305/go-yuv2webRTC/blob/master/webrtc/webrtc.go
package backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pion/webrtc/v2"
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

func getV8PlayLoadType(sdb string) (int, int) {
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

func main() {
}
