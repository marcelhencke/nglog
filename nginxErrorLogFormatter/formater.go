package nginxErrorLogFormatter

import (
	"fmt"
	"io"
	"nglog/config"
	debug "nglog/debugPrint"
	"regexp"
	"strings"
	"sync"
)

type NginxErrorLogFormatter struct {
	bufferMutex           sync.RWMutex
	buffer                string
	ioWriter              io.Writer
	rLogLineCompleteRegex *regexp.Regexp
	rLogLineFastCGIRegex  *regexp.Regexp
	rVarRegexNg           *regexp.Regexp
	rVarRegexTs           *regexp.Regexp
}

func New(w io.Writer) *NginxErrorLogFormatter {
	ptr := &NginxErrorLogFormatter{
		buffer:   "",
		ioWriter: w,
	}
	ptr.rLogLineCompleteRegex, _ = regexp.Compile(config.Flags.LogLineCompleteRegex)
	ptr.rLogLineFastCGIRegex, _ = regexp.Compile(config.Flags.LogLineFastCGIRegex)
	ptr.rVarRegexNg, _ = regexp.Compile(config.Flags.NginxVarRegex)
	ptr.rVarRegexTs, _ = regexp.Compile("\\d{4}[\\/-]\\d{2}[\\/-]\\d{2} \\d{2}:\\d{2}:\\d{2}")
	return ptr
}

func (this *NginxErrorLogFormatter) ReadBufferLine(line string) error {
	debug.Println("ReadBufferLine: " + line)
	this.buffer += line
	if this.rLogLineCompleteRegex.MatchString(this.buffer) {
		debug.Println("ReadBufferLine: LINE is complete (matches LogLineCompleteRegex)")
		e := this.parseFastCGILogLine(this.buffer)
		this.buffer = ""
		return e
	}
	this.buffer += "\n"
	return nil
}

func (this *NginxErrorLogFormatter) parseFastCGILogLine(logLine string) error {
	if this.rLogLineFastCGIRegex.MatchString(logLine) {
		matches := this.rLogLineFastCGIRegex.FindStringSubmatch(logLine)
		debug.Println("FastCGI Parser match LogLineFastCGIRegex")
		debug.Println("- PREFIX: " + matches[1])
		debug.Println("- BODY: " + matches[2])
		debug.Println("- SUFFIX: " + matches[3])
		return this.splitPhpLine(matches[1], matches[2], matches[3])
	}
	debug.Println("FastCGI Parser could NOT match LogLineFastCGIRegex: " + logLine)
	return this.writeOut(logLine)
}

func (this *NginxErrorLogFormatter) splitPhpLine(prefix, phpMsgs, suffix string) error {
	phpMsgsArr := strings.Split(phpMsgs, config.Flags.LogLinePhpMsgSplit)
	for _, phpMsg := range phpMsgsArr {
		if phpMsg == "" {
			continue
		}
		debug.Println("PHP message after split: " + phpMsg)
		err := this.writeOut(this.reformatPhpLine(prefix, phpMsg, suffix))
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *NginxErrorLogFormatter) reformatPhpLine(prefix, phpMsg, suffix string) string {
	if config.Flags.Template == "%raw%" {
		return prefix + phpMsg + suffix
	}

	dataTemplate := config.Flags.Template
	debug.Println("PHP format: " + dataTemplate)
	data := map[string]string{}
	data["raw"] = prefix + phpMsg + suffix
	data["prefix"] = prefix
	data["suffix"] = suffix
	data["php"] = phpMsg

	data["ts"] = this.rVarRegexTs.FindString(prefix)
	debug.Println("PHP format ts match: " + data["ts"])

	ngVars := this.rVarRegexNg.FindAllStringSubmatch(suffix, -1)
	for _, ngVar := range ngVars {
		if len(ngVar) > 2 {
			debug.Println("PHP format ng_" + ngVar[1] + " match: " + ngVar[2])
			data["ng_"+ngVar[1]] = ngVar[2]
		}
	}

	for k, v := range data {
		dataTemplate = strings.ReplaceAll(dataTemplate, "%"+k+"%", v)
	}
	return dataTemplate
}

func (this *NginxErrorLogFormatter) writeOut(data string) error {
	_, e := fmt.Fprintln(this.ioWriter, data)
	return e
}
