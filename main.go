package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/Kretech/xgo/version"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var opt Opt

type Opt struct {
	ShowHelp bool
	Verbose  bool

	Now  bool
	List bool

	Branch string

	Patch   bool
	Feature bool

	Path string
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func InitFlags(opt *Opt) {
	flag.BoolVar(&opt.ShowHelp, `h`, false, `show help`)
	flag.BoolVar(&opt.Verbose, `v`, false, `verbose`)

	flag.BoolVar(&opt.Now, `s`, false, `show latest tag`)
	flag.BoolVar(&opt.List, `l`, false, `list all tags`)
	flag.StringVar(&opt.Branch, `b`, `master`, `use which branch`)
	flag.BoolVar(&opt.Patch, `fix`, false, `add a fix tag`)
	flag.BoolVar(&opt.Feature, `ft`, false, `add a feature tag`)
	flag.StringVar(&opt.Path, `p`, `.`, `project path`)

	flag.Parse()

}

func main() {

	InitFlags(&opt)

	if opt.ShowHelp || len(os.Args) < 1 {
		flag.PrintDefaults()
		return
	}

	repo, err := git.PlainOpen(opt.Path)
	if err != nil {
		log.Printf("open repo(%s) err: %v", opt.Path, err)
		return
	}

	head, _ := repo.Head()

	tags := make([]*plumbing.Reference, 0, 32)

	tagIter, err := repo.Tags()
	{
		if err != nil {
			log.Println(err)
			return
		}

		err = tagIter.ForEach(func(r *plumbing.Reference) error {
			tags = append(tags, r)
			return nil
		})

		sort.Slice(tags, func(i, j int) bool {
			return version.GreaterThan(tags[i].Name().Short(), tags[j].Name().Short())
		})
	}

	svs := ``
	if len(tags) > 0 {
		svs = tags[0].Name().Short()
	}
	sv := version.Parse(svs)

	switch {
	case opt.List:
		for _, tag := range tags {
			fmt.Println(tag.Name().Short())
		}

	case opt.Patch, opt.Feature:
		var next *version.SemVer
		if opt.Patch {
			next = sv.NextPatch()
		} else if opt.Feature {
			next = sv.NextMinor()
		}

		ref, err := repo.CreateTag(next.String(), head.Hash(), nil)
		if err != nil {
			log.Fatalln(err)
			return
		}

		log.Println(`create tag success:`, ref)
	}
}
