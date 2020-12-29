import { Sdp } from './../shared/sdp';
import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute } from '@angular/router';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-peer-calling',
  templateUrl: './peer-calling.component.html',
  styleUrls: ['./peer-calling.component.scss']
})
export class PeerCallingComponent implements OnInit {

  pcSender: any;
  pcReceiver: any;
  meetingId: string;
  peerId: string;
  userId: string;

  constructor(
      private http: HttpClient,
      private route: ActivatedRoute
  ) { }

  ngOnInit() {
      // this.meetingId = this.route.snapshot.paramMap.get('meetingId');
      // this.peerId = this.route.snapshot.paramMap.get('peerId');
      // this.userId = this.route.snapshot.paramMap.get('userId');

      // Add dummy data
      this.meetingId = '07927fc8-af0a-11ea-b338-064f26a5f90a';
      this.peerId = 'bob';
      this.userId = 'alice';

      this.pcSender = new RTCPeerConnection({
          iceServers: [{
              urls: 'stun:stun.l.google.com:19302'
          }]
      });
      this.pcReceiver = new RTCPeerConnection({
          iceServers: [{
              urls: 'stun:stun.l.google.com:19302'
          }]
      });

      this.pcSender.onicecandidate = event => {
        if (event.candidate === null) {
          this.http.post<Sdp>(environment.basePath + '/webrtc/sdp/m/' + this.meetingId + "/c/"+ this.userId + "/p/" + this.peerId + "/s/" + true,
          {"sdp" : btoa(JSON.stringify(this.pcSender.localDescription))}).subscribe(response => {
            this.pcSender.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.Sdp))));
          });
        }
      };
      this.pcReceiver.onicecandidate = event => {
        if (event.candidate === null) {
            this.http.post<Sdp>(environment.basePath + '/webrtc/sdp/m/' + this.meetingId + "/c/"+ this.userId + "/p/" + this.peerId + "/s/" + false, 
            {"sdp" : btoa(JSON.stringify(this.pcReceiver.localDescription))}).subscribe(response => {
            this.pcReceiver.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(response.Sdp))));
          })
        }
      };
  }

  startCall() {
    // sender part of the call
      navigator.mediaDevices.getUserMedia({video: true, audio: true}).then((stream) =>{
          const senderVideo: any = document.getElementById('senderVideo');
          senderVideo.srcObject = stream;
          const tracks = stream.getTracks();
          for (let i = 0; i < tracks.length; i++) {
            this.pcSender.addTrack(stream.getTracks()[i]);
          }
          this.pcSender.createOffer().then(d => this.pcSender.setLocalDescription(d))
      });
      // you can use event listner so that you inform he is connected!
      this.pcSender.addEventListener('connectionstatechange', event => {
        if (this.pcSender.connectionState === 'connected') {
            console.log("horray!")
        }
      });

      // receiver part of the call
      this.pcReceiver.addTransceiver('video', {direction: 'recvonly'});

      this.pcReceiver.createOffer()
        .then(d => this.pcReceiver.setLocalDescription(d));

      this.pcReceiver.ontrack = (event) => {
        const receiverVideo: any = document.getElementById('receiverVideo');
        receiverVideo.srcObject = event.streams[0];
        receiverVideo.autoplay = true;
        receiverVideo.controls = true;
      };

  }

}
