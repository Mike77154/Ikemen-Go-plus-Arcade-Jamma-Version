package main

import (
	"encoding/binary"
	"fmt"
	"github.com/jfreymuth/go-vorbis/ogg/vorbis"
	"github.com/timshannon/go-openal/openal"
	"github.com/sqweek/fluidsynth"
	"io"
	"math"
	"os"
)

const (
	audioOutLen    = 2048
	audioFrequency = 48000
)

type AudioSource struct {
	Src  openal.Source
	bufs openal.Buffers
}

func NewAudioSource() (s *AudioSource) {
	s = &AudioSource{Src: openal.NewSource(), bufs: openal.NewBuffers(2)}
	for i := range s.bufs {
		s.bufs[i].SetDataInt16(openal.FormatStereo16, sys.nullSndBuf[:],
			audioFrequency)
	}
	s.Src.QueueBuffers(s.bufs)
	if err := openal.Err(); err != nil {
		println(err.Error())
	}
	return
}
func (s *AudioSource) Delete() {
	for s.Src.BuffersQueued() > 0 {
		s.Src.UnqueueBuffer()
	}
	s.bufs.Delete()
	s.Src.Delete()
}

type Mixer struct {
	buf        [audioOutLen * 2]float32
	sendBuf    []int16
	out        chan []int16
	normalizer *Normalizer
}

func newMixer() *Mixer {
	return &Mixer{out: make(chan []int16, 1), normalizer: NewNormalizer()}
}
func (m *Mixer) bufClear() {
	for i := range m.buf {
		m.buf[i] = 0
	}
}
func (m *Mixer) write() bool {
	if m.sendBuf == nil {
		m.sendBuf = make([]int16, len(m.buf))
		for i := 0; i <= len(m.sendBuf)-2; i += 2 {
			l, r := m.normalizer.Process(m.buf[i], m.buf[i+1])
			m.sendBuf[i] = int16(32767 * l)
			m.sendBuf[i+1] = int16(32767 * r)
		}
	}
	select {
	case m.out <- m.sendBuf:
	default:
		return false
	}
	m.sendBuf = nil
	m.bufClear()
	return true
}
func (m *Mixer) Mix(wav []byte, fidx float64, bytesPerSample, channels int,
	sampleRate float64, loop bool, volume float32) float64 {
	fidxadd := sampleRate / audioFrequency
	if fidxadd > 0 {
		switch bytesPerSample {
		case 1:
			switch channels {
			case 1:
				for i := 0; i <= len(m.buf)-2; i += 2 {
					iidx := int(fidx)
					if iidx >= len(wav) {
						if !loop {
							break
						}
						iidx, fidx = 0, 0
					}
					sam := volume * (float32(wav[iidx]) - 128) / 128
					m.buf[i] += sam
					m.buf[i+1] += sam
					fidx += fidxadd
				}
				return fidx
			case 2:
				for i := 0; i <= len(m.buf)-2; i += 2 {
					iidx := 2 * int(fidx)
					if iidx > len(wav)-2 {
						if !loop {
							break
						}
						iidx, fidx = 0, 0
					}
					m.buf[i] += volume * (float32(wav[iidx]) - 128) / 128
					m.buf[i+1] += volume * (float32(wav[iidx+1]) - 128) / 128
					fidx += fidxadd
				}
				return fidx
			}
		case 2:
			switch channels {
			case 1:
				for i := 0; i <= len(m.buf)-2; i += 2 {
					iidx := 2 * int(fidx)
					if iidx > len(wav)-2 {
						if !loop {
							break
						}
						iidx, fidx = 0, 0
					}
					sam := volume *
						float32(int(wav[iidx])|int(int8(wav[iidx+1]))<<8) / (1 << 15)
					m.buf[i] += sam
					m.buf[i+1] += sam
					fidx += fidxadd
				}
				return fidx
			case 2:
				for i := 0; i <= len(m.buf)-2; i += 2 {
					iidx := 4 * int(fidx)
					if iidx > len(wav)-4 {
						if !loop {
							break
						}
						iidx, fidx = 0, 0
					}
					m.buf[i] += volume *
						float32(int(wav[iidx])|int(int8(wav[iidx+1]))<<8) / (1 << 15)
					m.buf[i+1] += volume *
						float32(int(wav[iidx+2])|int(int8(wav[iidx+3]))<<8) / (1 << 15)
					fidx += fidxadd
				}
				return fidx
			}
		}
	}
	return float64(len(wav))
}

type Normalizer struct {
	mul  float64
	l, r *NormalizerLR
}

func NewNormalizer() *Normalizer {
	return &Normalizer{mul: 4, l: &NormalizerLR{1, 0, 1, 1 / 32.0, 0, 0},
		r: &NormalizerLR{1, 0, 1, 1 / 32.0, 0, 0}}
}
func (n *Normalizer) Process(l, r float32) (float32, float32) {
	lmul := n.l.process(n.mul, &l)
	rmul := n.r.process(n.mul, &r)
	if lmul < rmul {
		n.mul = lmul
	} else {
		n.mul = rmul
	}
	if n.mul > 16 {
		n.mul = 16
	}
	return l, r
}

type NormalizerLR struct {
	heri, herihenka, fue, heikin, katayori, katayori2 float64
}

func (n *NormalizerLR) process(bai float64, sam *float32) float64 {
	n.katayori = (n.katayori*audioFrequency/110 + float64(*sam)) /
		(audioFrequency/110.0 + 1)
	n.katayori2 = (n.katayori2*audioFrequency/112640 + float64(*sam)) /
		(audioFrequency/112640.0 + 1)
	s := (n.katayori2 - n.katayori) * bai
	if math.Abs(s) > 1 {
		bai *= math.Pow(1/math.Abs(s), n.heri)
		n.herihenka += 32 * (1 - n.heri) / float64(audioFrequency+32)
		if s < 0 {
			s = -1
		} else {
			s = 1
		}
	} else {
		tmp := (1 - math.Pow(1-math.Abs(s), 64)) * math.Pow(0.5-math.Abs(s), 3)
		bai += bai * (n.heri*(1/32.0-n.heikin)/n.fue + tmp*n.fue*(1-n.heri)/32) /
			(audioFrequency*2/8.0 + 1)
		n.herihenka -= (0.5 - n.heikin) * n.heri / (audioFrequency * 2)
	}
	n.fue += (32*n.fue*(1/n.fue-math.Abs(s)) - n.fue) /
		(32 * audioFrequency * 2)
	n.heikin += (math.Abs(s) - n.heikin) / (audioFrequency * 2)
	n.heri += n.herihenka
	if n.heri < 0 {
		n.heri = 0
	} else if n.heri > 0 {
		n.heri = 1
	}
	*sam = float32(s)
	return bai
}

type Vorbis struct {
	dec        *vorbis.Vorbis
	fp         *os.File
	buf        []int16
	bufi       float64
	openReq    chan string
	normalizer *Normalizer
}

func newVorbis() *Vorbis {
	return &Vorbis{openReq: make(chan string, 1), normalizer: NewNormalizer()}
}
func (v *Vorbis) Open(file string) {
	v.openReq <- file
}
func (v *Vorbis) openFile(file string) bool {
	v.clear()
	var err error
	if v.fp, err = os.Open(file); err != nil {
		return false
	}
	return v.restart()
}
func (v *Vorbis) restart() bool {
	if v.fp == nil {
		return false
	}
	_, err := v.fp.Seek(0, 0)
	chk(err)
	if v.dec, err = vorbis.Open(v.fp); err != nil {
		v.clear()
		return false
	}
	v.buf = nil
	return true
}
func (v *Vorbis) clear() {
	if v.dec != nil {
		v.dec = nil
	}
	if v.fp != nil {
		chk(v.fp.Close())
		v.fp = nil
	}
}
func (v *Vorbis) samToAudioOut(buf [][]float32) (out []int16) {
	var o1i int
	if len(buf) == 1 {
		o1i = 0
	} else {
		o1i = 1
	}
	sr := audioFrequency / float64(v.dec.SampleRate())
	out = make([]int16, 2*(int(float64(len(buf[0])-1)*sr)+1))
	oldbufi := -2
	for i := range buf[0] {
		for j := oldbufi + 2; j <= 2*int(v.bufi); j += 2 {
			l, r := v.normalizer.Process(buf[0][i], buf[o1i][i])
			out[j], out[j+1] = int16(25000*l), int16(25000*r)
		}
		oldbufi = 2 * int(v.bufi)
		v.bufi = sr * float64(i+1)
	}
	v.bufi -= float64(int(v.bufi))
	return
}
func (v *Vorbis) read() []int16 {
	select {
	case file := <-v.openReq:
		v.openFile(file)
	default:
	}
	for v.dec != nil {
		if len(v.buf) >= audioOutLen*2 {
			out := v.buf[:audioOutLen*2]
			v.buf = v.buf[audioOutLen*2:]
			return out
		}
		for len(v.buf) < audioOutLen*2 && v.dec != nil {
			sam, err := v.dec.DecodePacket()
			if err == io.EOF {
				v.restart()
				continue
			} else {
				chk(err)
			}
			v.buf = append(v.buf, v.samToAudioOut(sam)...)
		}
	}
	return sys.nullSndBuf[:]
}

type Wave struct {
	SamplesPerSec  uint32
	Channels       uint16
	BytesPerSample uint16
	Wav            []byte
}

func ReadWave(f *os.File, ofs int64) (*Wave, error) {
	buf := make([]byte, 4)
	n, err := f.Read(buf)
	if err != nil {
		return nil, err
	}
	if string(buf[:n]) != "RIFF" {
		return nil, Error("RIFFではありません")
	}
	read := func(x interface{}) error {
		return binary.Read(f, binary.LittleEndian, x)
	}
	var riffSize uint32
	if err := read(&riffSize); err != nil {
		return nil, err
	}
	riffSize += 8
	if n, err = f.Read(buf); err != nil {
		return nil, err
	}
	if string(buf[:n]) != "WAVE" {
		return &Wave{SamplesPerSec: 11025, Channels: 1, BytesPerSample: 1}, nil
	}
	fmtSize, dataSize := uint32(0), uint32(0)
	w := Wave{}
	riffend := ofs + 16 + int64(riffSize)
	ofs += 28
	for (fmtSize == 0 || dataSize == 0) && ofs < riffend {
		if n, err = f.Read(buf); err != nil {
			return nil, err
		}
		var size uint32
		if err := read(&size); err != nil {
			return nil, err
		}
		switch string(buf[:n]) {
		case "fmt ":
			fmtSize = size
			var fmtID uint16
			if err := read(&fmtID); err != nil {
				return nil, err
			}
			if fmtID != 1 {
				return nil, Error("リニアPCMではありません")
			}
			if err := read(&w.Channels); err != nil {
				return nil, err
			}
			if w.Channels < 1 || w.Channels > 2 {
				return nil, Error("チャンネル数が不正です")
			}
			if err := read(&w.SamplesPerSec); err != nil {
				return nil, err
			}
			if w.SamplesPerSec < 1 || w.SamplesPerSec >= 0xfffff {
				return nil, Error(fmt.Sprintf("周波数が不正です %d", w.SamplesPerSec))
			}
			var musi uint32
			if err := read(&musi); err != nil {
				return nil, err
			}
			var mushi uint16
			if err := read(&mushi); err != nil {
				return nil, err
			}
			if err := read(&w.BytesPerSample); err != nil {
				return nil, err
			}
			if w.BytesPerSample != 8 && w.BytesPerSample != 16 {
				return nil, Error("bit数が不正です")
			}
			w.BytesPerSample >>= 3
		case "data":
			dataSize = size
			w.Wav = make([]byte, dataSize)
			if err := binary.Read(f, binary.LittleEndian, w.Wav); err != nil {
				return nil, err
			}
		}
		ofs += int64(size) + 8
		f.Seek(ofs, 0)
	}
	if fmtSize == 0 {
		if dataSize > 0 {
			return nil, Error("fmt がありません")
		}
		return nil, nil
	}
	return &w, nil
}

type Snd struct {
	table     map[[2]int32]*Wave
	ver, ver2 uint16
}

func newSnd() *Snd { return &Snd{table: make(map[[2]int32]*Wave)} }
func LoadSnd(filename string) (*Snd, error) {
	s := newSnd()
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() { chk(f.Close()) }()
	buf := make([]byte, 12)
	var n int
	if n, err = f.Read(buf); err != nil {
		return nil, err
	}
	if string(buf[:n]) != "ElecbyteSnd\x00" {
		return nil, Error("ElecbyteSndではありません")
	}
	read := func(x interface{}) error {
		return binary.Read(f, binary.LittleEndian, x)
	}
	if err := read(&s.ver); err != nil {
		return nil, err
	}
	if err := read(&s.ver2); err != nil {
		return nil, err
	}
	var numberOfSounds uint32
	if err := read(&numberOfSounds); err != nil {
		return nil, err
	}
	var subHeaderOffset uint32
	if err := read(&subHeaderOffset); err != nil {
		return nil, err
	}
	for i := uint32(0); i < numberOfSounds; i++ {
		f.Seek(int64(subHeaderOffset), 0)
		var nextSubHeaderOffset uint32
		if err := read(&nextSubHeaderOffset); err != nil {
			return nil, err
		}
		var subFileLenght uint32
		if err := read(&subFileLenght); err != nil {
			return nil, err
		}
		var num [2]int32
		if err := read(&num); err != nil {
			return nil, err
		}
		if num[0] >= 0 && num[1] >= 0 {
			_, ok := s.table[num]
			if !ok {
				tmp, err := ReadWave(f, int64(subHeaderOffset))
				if err != nil {
					return nil, err
				}
				s.table[num] = tmp
			}
		}
		subHeaderOffset = nextSubHeaderOffset
	}
	return s, nil
}
func (s *Snd) Get(gn [2]int32) *Wave {
	return s.table[gn]
}
func (s *Snd) play(gn [2]int32) bool {
	c := sys.sounds.GetChannel()
	if c == nil {
		return false
	}
	c.sound = s.Get(gn)
	return c.sound != nil
}

type Sound struct {
	sound   *Wave
	volume  int16
	loop    bool
	freqmul float32
	fidx    float64
}

func (s *Sound) mix() {
	if s.sound == nil {
		return
	}
	s.fidx = sys.mixer.Mix(s.sound.Wav, s.fidx,
		int(s.sound.BytesPerSample), int(s.sound.Channels),
		float64(s.sound.SamplesPerSec)*float64(s.freqmul), s.loop,
		float32(s.volume)/256)
	if int(s.fidx) >= len(s.sound.Wav)/
		int(s.sound.BytesPerSample*s.sound.Channels) {
		s.sound = nil
		s.fidx = 0
	}
}
func (s *Sound) SetVolume(vol int32) {
	if vol < 0 {
		s.volume = 0
	} else if vol > 512 {
		s.volume = 512
	} else {
		s.volume = int16(vol)
	}
}
func (s *Sound) SetPan(pan float32, offset *float32) {
	// 未実装
}

type Sounds []Sound

func newSounds(size int) (s Sounds) {
	s = make(Sounds, size)
	for i := range s {
		s[i] = Sound{volume: 256, freqmul: 1}
	}
	return
}
func (s Sounds) GetChannel() *Sound {
	for i := range s {
		if s[i].sound == nil {
			return &s[i]
		}
	}
	return nil
}
func (s Sounds) mixSounds() {
	for i := range s {
		s[i].mix()
	}
}
