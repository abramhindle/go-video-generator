#!/bin/bash
WAVFILE="$1"
BW=`basename "$WAVFILE"`
echo "BW: $BW"
shift 1
OUTPREFIX=`date +"%s.$RANDOM-out.$BW" | sed -e 's/ /_/g'`
OUT=${OUTPREFIX}.mkv
OUTA=${OUTPREFIX}.audio.mkv
echo $OUT
find $* |rand| XARGS go run govid3.go -out $OUT -frames `bash getwavsize.sh "$WAVFILE"` && \
ffmpeg -i $OUT -i "$WAVFILE" -vcodec copy $OUTA  && rm $OUT
echo $OUTA
