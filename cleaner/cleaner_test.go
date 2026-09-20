package cleaner

import (
	"encoding/binary"
	"hash/crc32"
	"os"
	"path/filepath"
	"testing"
)

func pngChunk(typ string, payload []byte) []byte { b:=make([]byte,12+len(payload));binary.BigEndian.PutUint32(b[:4],uint32(len(payload)));copy(b[4:8],typ);copy(b[8:],payload);binary.BigEndian.PutUint32(b[8+len(payload):],crc32.ChecksumIEEE(b[4:8+len(payload)]));return b }
func samplePNG() []byte { b:=[]byte{137,80,78,71,13,10,26,10}; b=append(b,pngChunk("IHDR",make([]byte,13))...);b=append(b,pngChunk("tEXt",[]byte("Author\x00Private"))...);b=append(b,pngChunk("IDAT",[]byte{1,2,3})...);b=append(b,pngChunk("IEND",nil)...);return b }
func sampleJPEG() []byte { return []byte{0xff,0xd8, 0xff,0xe1,0x00,0x08,'E','x','i','f',0,0, 0xff,0xfe,0x00,0x05,'h','i','!', 0xff,0xda,0x00,0x08,1,2,3,4,5,6, 9,8,7, 0xff,0xd9} }

func TestCleanPNGRemovesTextAndPreservesImageChunks(t *testing.T){ clean,n,err:=cleanPNG(samplePNG());if err!=nil{t.Fatal(err)};if n!=1{t.Fatalf("removed=%d",n)};if len(clean)>=len(samplePNG()){t.Fatal("metadata was not removed")};if string(clean[12:16])!="IHDR"{t.Fatal("IHDR not preserved")} }
func TestCleanJPEGRemovesEXIFAndComment(t *testing.T){ clean,n,err:=cleanJPEG(sampleJPEG());if err!=nil{t.Fatal(err)};if n!=2{t.Fatalf("removed=%d",n)};if len(clean)>=len(sampleJPEG()){t.Fatal("metadata was not removed")} }
func TestCleanFileRefusesOverwrite(t *testing.T){ d:=t.TempDir();in:=filepath.Join(d,"in.png");out:=filepath.Join(d,"out.png");os.WriteFile(in,samplePNG(),0600);os.WriteFile(out,[]byte("keep"),0600);if _,err:=CleanFile(in,out,false);err==nil{t.Fatal("expected overwrite refusal")};got,_:=os.ReadFile(out);if string(got)!="keep"{t.Fatal("existing output changed")} }
func TestCleanFileWritesOutput(t *testing.T){ d:=t.TempDir();in:=filepath.Join(d,"in.png");out:=filepath.Join(d,"nested","out.png");os.WriteFile(in,samplePNG(),0600);r,err:=CleanFile(in,out,false);if err!=nil{t.Fatal(err)};if r.RemovedBlocks!=1||r.Format!="png"{t.Fatalf("unexpected result: %+v",r)};if _,err:=os.Stat(out);err!=nil{t.Fatal(err)} }
func TestRejectsInPlace(t *testing.T){ d:=t.TempDir();p:=filepath.Join(d,"x.png");os.WriteFile(p,samplePNG(),0600);if _,err:=CleanFile(p,p,true);err==nil{t.Fatal("expected in-place refusal")} }
func TestRejectsUnsupported(t *testing.T){ d:=t.TempDir();in:=filepath.Join(d,"x.txt");os.WriteFile(in,[]byte("hello"),0600);if _,err:=CleanFile(in,filepath.Join(d,"o.txt"),false);err==nil{t.Fatal("expected unsupported error")} }
