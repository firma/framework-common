package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return
}

func ConvertBytesToReader(data []byte) io.Reader {
	return &byteReader{data: data}
}

func ConvertOggToPCM(data []byte) ([]byte, error) {

	inStream := ConvertBytesToReader(data)
	outFormat := &AudioFormat{
		AudioCodec:   "pcm_s16le",
		SampleFormat: "s16le",
		Channel:      "1",
		SampleRate:   16000,
	}

	cmd := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-acodec", "mp3", "-ac", outFormat.Channel, "-ar", fmt.Sprintf("%d", outFormat.SampleRate), "-")
	cmd.Stdin = inStream
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if _, err := io.Copy(&out, stdout); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	if out.Len() == 0 {
		return nil, errors.New("转换后的数据为空")
	}

	return out.Bytes(), nil
}

type AudioFormat struct {
	AudioCodec   string
	SampleFormat string
	Channel      string
	SampleRate   int
}

func (f *AudioFormat) String() string {
	return fmt.Sprintf("%s, %s, %s, %d Hz", f.AudioCodec, f.SampleFormat, f.Channel, f.SampleRate)
}

//func main() {
//	inFile := "in.ogg"
//	formats := []*AudioFormat{
//		{
//			AudioCodec:   "pcm_s16le",
//			SampleFormat: "s16le",
//			Channel:      "1",
//			SampleRate:   8000,
//		},
//		{
//			AudioCodec:   "pcm_s16le",
//			SampleFormat: "s16le",
//			Channel:      "1",
//			SampleRate:   16000,
//		},
//	}
//
//	for _, outFormat := range formats {
//		data, err := convertAudioToPCM(inFile, outFormat)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Printf("转换后的 PCM 数据: %d bytes, 格式：%s\n", len(data), outFormat)
//		// 处理数据...
//	}
//}

func OgaToPcm(ogaData []byte) ([]byte, error) {
	// OGA文件头部信息
	const headerSize = 44
	const sampleRateOffset = 24
	const numChannelsOffset = 22
	const bitsPerSampleOffset = 34
	const dataSizeOffset = 40

	// 读取OGA文件头部信息
	sampleRate := binary.LittleEndian.Uint32(ogaData[sampleRateOffset : sampleRateOffset+4])
	numChannels := binary.LittleEndian.Uint16(ogaData[numChannelsOffset : numChannelsOffset+2])
	bitsPerSample := binary.LittleEndian.Uint16(ogaData[bitsPerSampleOffset : bitsPerSampleOffset+2])
	dataSize := binary.LittleEndian.Uint32(ogaData[dataSizeOffset : dataSizeOffset+4])

	// 计算PCM文件头部信息
	pcmDataSize := dataSize * 2
	pcmFileSize := pcmDataSize + headerSize - 8
	pcmSampleRate := sampleRate
	pcmNumChannels := numChannels
	pcmBitsPerSample := bitsPerSample

	// 构建PCM文件头部信息
	pcmHeader := new(bytes.Buffer)
	pcmHeader.WriteString("RIFF")
	binary.Write(pcmHeader, binary.LittleEndian, pcmFileSize)
	pcmHeader.WriteString("WAVEfmt ")
	binary.Write(pcmHeader, binary.LittleEndian, uint32(16))
	binary.Write(pcmHeader, binary.LittleEndian, uint16(1))
	binary.Write(pcmHeader, binary.LittleEndian, pcmNumChannels)
	binary.Write(pcmHeader, binary.LittleEndian, pcmSampleRate)
	binary.Write(pcmHeader, binary.LittleEndian, uint32(pcmSampleRate*uint32(pcmNumChannels*pcmBitsPerSample/8)))
	binary.Write(pcmHeader, binary.LittleEndian, uint16(pcmNumChannels*pcmBitsPerSample/8))
	binary.Write(pcmHeader, binary.LittleEndian, pcmBitsPerSample)
	pcmHeader.WriteString("data")
	binary.Write(pcmHeader, binary.LittleEndian, pcmDataSize)

	// 将OGA数据流转换为PCM数据流
	pcmData := make([]byte, pcmDataSize)
	for i := 0; i < int(dataSize); i++ {
		for j := 0; j < int(numChannels); j++ {
			pcmData[i*int(numChannels)*2+j*2] = ogaData[i*int(numChannels)*2+j*4]
			pcmData[i*int(numChannels)*2+j*2+1] = ogaData[i*int(numChannels)*2+j*4+1]
		}
	}

	// 将PCM头部信息和数据流合并
	pcm := make([]byte, len(pcmHeader.Bytes())+len(pcmData))
	copy(pcm, pcmHeader.Bytes())
	copy(pcm[len(pcmHeader.Bytes()):], pcmData)

	return pcm, nil
}
