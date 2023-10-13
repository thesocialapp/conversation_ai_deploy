package api

import (
	"fmt"
	"io"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/rs/zerolog/log"
)

type PeerConnection struct {
	// The peer connection
	peerConn *webrtc.PeerConnection
}

// Close connection
func (pc *PeerConnection) Close() error {
	return pc.peerConn.Close()
}

func (pc *PeerConnection) handleAudioTrack(track *webrtc.TrackRemote, dc *webrtc.DataChannel) error {
	decoder, err := newDecoder()
	if err != nil {
		return err
	}
	errs := make(chan error, 2)
	audioStream := make(chan []byte)
	response := make(chan bool)
	timer := time.NewTimer(5 * time.Second)

	go func() {
		for {
			packet, _, err := track.ReadRTP()
			timer.Reset(1 * time.Second)
			if err != nil {
				timer.Stop()
				if err == io.EOF {
					log.Info().Msg("track has ended")
					close(audioStream)
					return
				}
				errs <- err
				return
			}
			audioStream <- packet.Payload
			<-response
		}
	}()

	err = nil
	for {
		select {
		case audioChunk := <-audioStream:
			// TODO: Send Payload to the transcriber
			_, err := decoder.decode(audioChunk)
			if err != nil {
				return err
			}
			/// Send payload to our trascriber stream
			log.Info().Msg("sending audio chunk")
		case <-timer.C:
			return fmt.Errorf("Read operation timed out")
		case err = <-errs:
			log.Info().Msgf("error reading track: %s %s", track.ID(), err.Error())
			return err
		}

	}
}

func (pi *PeerConnection) NewPeerConnection() (*PeerConnection, error) {
	pcConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}
	pc, err := webrtc.NewPeerConnection(pcConfig)
	if err != nil {
		return nil, err
	}

	dataChan := make(chan *webrtc.DataChannel)
	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dataChan <- dc
	})

	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		if track.Codec().MimeType == "audio/opus" {
			log.Info().Msgf("Received audio (%s) track, id = %s\n", track.Codec().MimeType, track.ID())
			err := pi.handleAudioTrack(track, <-dataChan)
			if err != nil {
				log.Error().Err(err).Msg("cannot handle audio track")
			}
		}
	})

	_, err = pc.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio, webrtc.RTPTransceiverInit{
		Direction: webrtc.RTPTransceiverDirectionRecvonly,
	})

	if err != nil {
		return nil, err
	}

	return &PeerConnection{
		peerConn: pc,
	}, nil
}
