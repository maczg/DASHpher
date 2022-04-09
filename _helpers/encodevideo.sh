#!/bin/bash

ffmpeg -i video.mp4 -map 0:v:0 -map 0:a\?:0 -map 0:v:0 -map 0:a\?:0 \
  -b:v:0 4300k -maxrate:v:0 4300k -bufsize:v:0 8600k -c:v:0 libx264 -filter:v:0 "scale=-2:1080" -movflags faststart -profile:v:0 main -preset fast -an video1_4300k.mp4 \
  -b:v:1 1050k -maxrate:v:1 1050k -bufsize:v:1 8600k -c:v:1 libx264 -filter:v:1 "scale=-2:480" -movflags faststart -profile:v:1 main -preset fast -an video1_1050k.mp4

MP4Box -dash 4000 -frag 4000 -profile full -rap -segment-name %s_segment -fps 24 \
          video_1050k.mp4#video:id=480p \
          video_1050k.mp4#video:id=480p \
          -url-template \
          -out video.mpd


