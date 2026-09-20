// Package cleaner removes common privacy-sensitive metadata from JPEG and PNG images.
package cleaner

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Result struct { Input, Output, Format string; OriginalBytes, CleanBytes int64; RemovedBlocks int }

func CleanFile(input, output string, overwrite bool) (Result, error) {
	inAbs, err := filepath.Abs(input); if err != nil { return Result{}, err }
	outAbs, err := filepath.Abs(output); if err != nil { return Result{}, err }
	if inAbs == outAbs { return Result{}, errors.New("input and output must be different; in-place cleaning is intentionally disabled") }
	data, err := os.ReadFile(inAbs); if err != nil { return Result{}, fmt.Errorf("read input: %w", err) }
	var clean []byte; var removed int; var format string
	switch {
	case len(data)>=2 && data[0]==0xff && data[1]==0xd8:
		format="jpeg"; clean, removed, err = cleanJPEG(data)
	case len(data)>=8 && bytes.Equal(data[:8], []byte{137,80,78,71,13,10,26,10}):
		format="png"; clean, removed, err = cleanPNG(data)
	default: return Result{}, errors.New("unsupported format: only JPEG and PNG are supported")
	}
	if err != nil { return Result{}, err }
	if !overwrite { if _, e:=os.Stat(outAbs); e==nil { return Result{}, errors.New("output already exists; use --overwrite to replace it") } else if !os.IsNotExist(e) { return Result{}, e } }
	if err:=os.MkdirAll(filepath.Dir(outAbs),0755); err!=nil{return Result{},err}
	tmp, err:=os.CreateTemp(filepath.Dir(outAbs), ".metadata-cleaner-*"); if err!=nil{return Result{},err}; tmpName:=tmp.Name(); defer os.Remove(tmpName)
	w:=bufio.NewWriter(tmp); if _,err=w.Write(clean);err!=nil{tmp.Close();return Result{},err}; if err=w.Flush();err!=nil{tmp.Close();return Result{},err}; if err=tmp.Sync();err!=nil{tmp.Close();return Result{},err}; if err=tmp.Close();err!=nil{return Result{},err}
	if overwrite { _=os.Remove(outAbs) }; if err=os.Rename(tmpName,outAbs);err!=nil{return Result{},fmt.Errorf("commit output: %w",err)}
	return Result{inAbs,outAbs,format,int64(len(data)),int64(len(clean)),removed},nil
}

func cleanJPEG(data []byte)([]byte,int,error){
	out:=append([]byte{},data[:2]...); removed:=0; i:=2
	for i<len(data){
		if data[i]!=0xff{return nil,0,errors.New("malformed JPEG marker")}; start:=i; for i<len(data)&&data[i]==0xff{i++}; if i>=len(data){return nil,0,errors.New("truncated JPEG")}; marker:=data[i]; i++
		if marker==0xd9 { out=append(out,data[start:i]...); return out,removed,nil }
		if marker==0xda { if i+2>len(data){return nil,0,errors.New("truncated JPEG scan")}; l:=int(binary.BigEndian.Uint16(data[i:i+2])); if l<2||i+l>len(data){return nil,0,errors.New("invalid JPEG scan length")}; out=append(out,data[start:]...); return out,removed,nil }
		if marker==0x01 || (marker>=0xd0&&marker<=0xd7){out=append(out,data[start:i]...);continue}
		if i+2>len(data){return nil,0,errors.New("truncated JPEG segment")}; l:=int(binary.BigEndian.Uint16(data[i:i+2])); if l<2||i+l>len(data){return nil,0,errors.New("invalid JPEG segment length")}; end:=i+l
		// APP1 commonly stores EXIF/XMP; APP13 stores IPTC/Photoshop metadata; COM stores comments.
		if marker==0xe1||marker==0xed||marker==0xfe { removed++ } else { out=append(out,data[start:end]...) }; i=end
	}
	return nil,0,errors.New("JPEG ended unexpectedly")
}

func cleanPNG(data []byte)([]byte,int,error){
	out:=append([]byte{},data[:8]...); removed:=0; i:=8
	for i<len(data){
		if i+12>len(data){return nil,0,errors.New("truncated PNG chunk")}; n:=int(binary.BigEndian.Uint32(data[i:i+4])); if n<0||i+12+n>len(data){return nil,0,errors.New("invalid PNG chunk length")}; typ:=string(data[i+4:i+8]); end:=i+12+n
		remove:=typ=="tEXt"||typ=="zTXt"||typ=="iTXt"||typ=="eXIf"||typ=="tIME"
		if remove { removed++ } else { out=append(out,data[i:end]...) }; i=end
		if typ=="IEND" { if i!=len(data){return nil,0,errors.New("unexpected bytes after PNG IEND")}; return out,removed,nil }
	}
	return nil,0,errors.New("PNG missing IEND")
}

func SuggestedOutput(input string) string { ext:=filepath.Ext(input); base:=strings.TrimSuffix(input,ext); return base+".clean"+ext }
