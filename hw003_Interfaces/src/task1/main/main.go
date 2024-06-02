package main

import "fmt"

type Formatter interface {
	Format(s string) string
}
type BaseTextFormatter struct{}
type BoldFormatter struct{}
type ItalicFormatter struct{}
type CodeFormatter struct{}
type ChainFormatter struct {
	formatters []Formatter
}

func (t BaseTextFormatter) Format(text string) string {
	return text
}
func (b BoldFormatter) Format(text string) string {
	return fmt.Sprintf("**%s**", text)
}
func (i ItalicFormatter) Format(text string) string {
	return fmt.Sprintf("_%s_", text)
}
func (c CodeFormatter) Format(text string) string {
	return fmt.Sprintf("`%s`", text)
}
func (chain *ChainFormatter) AddFormatter(formatters ...Formatter) {
	for _, formatter := range formatters {
		chain.formatters = append(chain.formatters, formatter)
	}
}
func (chain ChainFormatter) Format(text string) string {
	for _, formatter := range chain.formatters {
		text = formatter.Format(text)
	}
	return text
}

func main() {}
