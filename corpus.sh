for file in /opt/hindle1/Videos/mp4s/candidates/*corpus; do go run govid2.go $file/*png; mv out.mkv `basename $file`; done
#X=VID_20130627_130944.mp4.corpus VID_20130629_193251.mp4.corpus VID_20130812_211005.mp4.corpus VID_20130812_211038.mp4.corpus VID_20130824_191904.mp4.corpus VID_20130824_191958.mp4.corpus VID_20130824_195928.mp4.corpus
#for file in $X; do go run govid2.go  /opt/hindle1/Videos/mp4s/candidates/$file/*png; mv out.mkv `basename $file`; done
