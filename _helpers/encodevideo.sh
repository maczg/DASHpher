#!/bin/bash

#TODO keep only track 0 (video)
#ffmpeg -i video.mp4 -map 0:v:0 \
#  -b:v 4300k -maxrate:v 4300k -bufsize 8600k -c libx264 -filter:v:0 "scale=-2:1080" -movflags faststart -profile:v main -preset fast -an video_4300k.mp4 \
#  -b:v 1050k -maxrate:v 1050k -bufsize 8600k -c libx264 -filter:v:1 "scale=-2:480" -movflags faststart -profile:v main -preset fast -an video_1050k.mp4

#MP4Box -dash 4000 -frag 4000 -profile full -rap -segment-name %s_segment -fps 24 \
#          video_1050k.mp4#video:id=480p \
#          video_4300k.mp4#video:id=1080p \
#          -url-template \
#          -out /videofiles/$1/video.mpd
#
#          ##$1 video_id
#tar -cvf /videofiles/$1.tar -C /videofiles/ $1/

if [ ! -f /video/$1/video.mp4 ]; then
    exit 1
fi
mkdir -p /videofiles/$1/

if [ "$IS_CLOUD" = 'true' ]; then

ffmpeg -y -i /video/$1/video.mp4 -c:v libx264 \
 -r 24 -x264opts 'keyint=48:min-keyint=48:no-scenecut' \
 -vf scale=-2:1080 -b:v 4300k -maxrate 4300k \
 -movflags faststart -bufsize 8600k \
 -profile:v main -preset fast -an /video/$1/video_4300k.mp4

ffmpeg -y -i /video/$1/video.mp4 -c:v libx264 \
 -r 24 -x264opts 'keyint=48:min-keyint=48:no-scenecut' \
 -vf scale=-2:480 -b:v 1050k -maxrate 1050k \
 -movflags faststart -bufsize 8600k \
  -profile:v main -preset fast -an /video/$1/video_1050k.mp4

MP4Box -dash 4000 -frag 4000 -profile full -rap -segment-name %s_segment -fps 24 \
        /video/$1/video_1050k.mp4#video:id=480p \
        -url-template \
        -out /videofiles/$1/video.mpd


tar -cvf /videofiles/$1.tar -C /videofiles/ $1/

else
  echo "TODO (Offline/Online Encoding)"
fi


ffmpeg -i video.mp4 -map 0:v -c:v libx264 \
  -b:v 4300k -maxrate:v 4300k -bufsize 8600k -c libx264 -filter:v:0 "scale=-2:1080" -movflags faststart -profile:v main -preset fast -an video_4300k.mp4 \
  -b:v 1050k -maxrate:v 1050k -bufsize 8600k -c libx264 -filter:v:1 "scale=-2:480" -movflags faststart -profile:v main -preset fast -an video_1050k.mp4