package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rad03i2/metadata-cleaner/cleaner"
)

const version = "1.0.0"

func main(){ os.Exit(run()) }
func run() int {
	out:=flag.String("output","","output path (default: <name>.clean.<ext>)")
	overwrite:=flag.Bool("overwrite",false,"allow replacing an existing output file")
	asJSON:=flag.Bool("json",false,"print machine-readable JSON")
	showVersion:=flag.Bool("version",false,"print version")
	flag.Usage=func(){fmt.Fprintf(flag.CommandLine.Output(),"Metadata Cleaner %s — privacy-focused JPEG/PNG metadata removal\nAuthor: Radwan Abdulhadi Ahmed (@rad03i2)\n\nUsage: metadata-cleaner [options] <image>\n\nOptions:\n",version);flag.PrintDefaults()}
	flag.Parse()
	if *showVersion {fmt.Printf("metadata-cleaner %s — Radwan Abdulhadi Ahmed (@rad03i2)\n",version);return 0}
	if flag.NArg()!=1 {flag.Usage();return 2}
	input:=flag.Arg(0); output:=*out; if output==""{output=cleaner.SuggestedOutput(input)}
	result,err:=cleaner.CleanFile(input,output,*overwrite); if err!=nil{fmt.Fprintln(os.Stderr,"error:",err);return 1}
	if *asJSON { b,_:=json.MarshalIndent(result,"","  ");fmt.Println(string(b));return 0 }
	fmt.Printf("Cleaned %s -> %s\nFormat: %s | removed blocks: %d | bytes: %d -> %d\n",result.Input,result.Output,result.Format,result.RemovedBlocks,result.OriginalBytes,result.CleanBytes)
	return 0
}
