package voice

var supportedFileType = [...]string{"wav", "raw", "voc", "au"}

// AudioFormat , the audio format
type AudioFormat struct {
	sampleRate     int64
	numChannels    int64
	bytesPerSample int64
}

// CalcBytesPerSecond , calculate the bytes per second
func (afmt AudioFormat) CalcBytesPerSecond() int64 {
	return afmt.sampleRate * afmt.numChannels * afmt.bytesPerSample
}

// func WaveFileGetFormat() {}

// func WaveFileSetFormat() {}
