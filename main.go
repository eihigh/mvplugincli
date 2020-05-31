package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
	"unicode"

	"github.com/eihigh/filetest"
	"github.com/manifoldco/promptui"
)

func main() {
	arg := struct {
		Plugin, Author, PluginJP, AuthorJP, Dir, Date string
	}{}
	var err error

	arg.Plugin, err = ask("プラグインの名前（英字のみ）", identValidate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("// %s.js v1.0.0\n", arg.Plugin)

	arg.Author, err = ask("作者の名前（英字のみ）", asciiValidate)
	if err != nil {
		log.Fatal(err)
	}
	year := time.Now().Year()
	fmt.Printf("// (C) %d \"%s\"\n", year, arg.Author)

	arg.PluginJP, err = ask("プラグインの日本語名", emptyValidate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("// * @plugindesc %s\n", arg.PluginJP)

	arg.AuthorJP, err = ask("作者の日本語名", emptyValidate)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("// * @author %s\n", arg.AuthorJP)

	prompt := promptui.Prompt{
		Label:    "プラグインを作成するフォルダ（デフォルト: js/plugins）",
		Default:  "js/plugins",
		Validate: dirValidate,
	}
	arg.Dir, err = prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	arg.Date = time.Now().Format("2006-01-02")

	t, err := template.ParseFiles("plugin.template.js")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(filepath.Join(arg.Dir, arg.Plugin+".js"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := t.Execute(f, arg); err != nil {
		log.Fatal(err)
	}

	fmt.Println("プラグインの作成に成功しました！")
}

func ask(label string, val promptui.ValidateFunc) (string, error) {
	prompt := promptui.Prompt{Label: label, Validate: val}
	return prompt.Run()
}

func identValidate(s string) error {
	if s == "" {
		return fmt.Errorf("空欄は入力できません")
	}
	for _, c := range s {
		if 'a' <= c && c <= 'z' {
			continue
		}
		if 'A' <= c && c <= 'Z' {
			continue
		}
		if c == '_' {
			continue
		}
		return fmt.Errorf("%s: 英字以外が含まれています", s)
	}
	return nil
}

func asciiValidate(s string) error {
	if s == "" {
		return fmt.Errorf("空欄は入力できません")
	}
	for _, c := range s {
		if c > unicode.MaxASCII {
			return fmt.Errorf("%s: 英字以外が含まれています", s)
		}
	}
	return nil
}

func emptyValidate(s string) error {
	if s == "" {
		return fmt.Errorf("空欄は入力できません")
	}
	return nil
}

func dirValidate(s string) error {
	if s == "" {
		return fmt.Errorf("空欄は入力できません（現在のフォルダに作成する場合は . を入力）")
	}
	if !filetest.IsDir(s) {
		return fmt.Errorf("%s はディレクトリではありません", s)
	}
	return nil
}
