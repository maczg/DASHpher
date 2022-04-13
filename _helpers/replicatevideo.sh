while read -r video_id; do
    echo "$video_id"
   # cp /var/cloud/videos/encoded_video/video.mpd /var/cloud/videofiles/"$video_id"
  #  cp -r /var/cloud/videos/encoded_video/*segment* /var/cloud/videofiles/"$video_id"
    tar -cvf /videofiles/"video_id".tar -C /videofiles/ "video_id"/
    #ln -s /videofiles/"$video_id".tar.gz /videofiles/"$video_id.tar"
done < /var/cloud/videofiles/folders.txt