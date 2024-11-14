package main

import (
	"github.com/GnoCheckBot/demo/client"
	c "github.com/GnoCheckBot/demo/condition"
	r "github.com/GnoCheckBot/demo/requirement"
)

type automaticCheck struct {
	Description string
	If          c.Condition
	Then        r.Requirement
}

type manualCheck struct {
	Description string
	If          c.Condition
	Teams       []string
}

func config(gh *client.GitHub) ([]automaticCheck, []manualCheck) {
	auto := []automaticCheck{
		{
			Description: "Changes on 'foo' file should be reviewed by at least one [Bot PR](https://github.com/gnolang/gno/pull/3037) reviewer",
			If: c.And(
				c.BaseBranch("main"),
				c.Or(
					c.FileChanged(gh, "foo"),
					c.FileChanged(gh, "bar"),
					c.FileChanged(gh, "baz"),
				),
			),
			Then: r.And(
				r.Author("aeddi"),
				r.Or(
					// r.ReviewByUser(gh, "thehowl"), // Stop bothering Morgan
					r.ReviewByUser(gh, "ltzmaxwell"),
					r.ReviewByUser(gh, "zivkovicmilos"),
				),
			),
		},
		{
			Description: "Maintainer must be able to edit this pull request",
			If:          c.Always(),
			Then:        r.MaintainerCanModify(),
		},
		{
			Description: "Pull request head branch must be up to date with its base",
			If:          c.Always(), // Or only if c.BaseBranch("main") ?
			Then:        r.UpToDateWith(gh, r.PR_BASE),
		},
		{
			Description: "Label bug is applied for no other reason than testing",
			If:          c.HeadBranch("demo-pr"),
			Then:        r.Label(gh, "bug"),
		},
	}

	manual := []manualCheck{
		{
			Description: "Demo manual check",
			If:          c.Always(),
		},
		{
			Description: "Demo manual check with lot of (useless) details",
			If: c.Or(
				c.Always(),
				c.And(
					c.Label("bug"),
					c.BaseBranch("main"),
					c.Or(
						c.FileChanged(gh, "misc/deployments"),
						c.FileChanged(gh, `misc/docker-\.*`),
						c.FileChanged(gh, "tm2/pkg/p2p"),
						c.FileChanged(gh, "contribs"),
					),
				),
			),
			Teams: []string{"foo", "bar", "baz"},
		},
		{
			Description: "Determine if infra needs to be updated",
			If: c.And(
				c.BaseBranch("main"),
				c.Or(
					c.FileChanged(gh, "misc/deployments"),
					c.FileChanged(gh, `misc/docker-\.*`),
					c.FileChanged(gh, "tm2/pkg/p2p"),
				),
			),
			Teams: []string{"tech-staff"},
		},
		{
			Description: "The code style is satisfactory",
			If: c.And(
				c.BaseBranch("main"),
				c.Or(
					c.FileChanged(gh, `.*\.go`),
					c.FileChanged(gh, `.*\.js`),
				),
			),
			Teams: []string{"tech-staff"},
		},
		{
			Description: "The documentation is accurate and relevant",
			If:          c.FileChanged(gh, `.*\.md`),
			Teams: []string{
				"tech-staff",
				"devrels",
			},
		},
	}

	// Check for duplicates in manual rule descriptions
	// (need to be unique for the bot operations)
	unique := make(map[string]struct{})
	for _, rule := range manual {
		if _, exists := unique[rule.Description]; exists {
			gh.Logger.Fatalf("Manual rule description must be unique (duplicate : %s)", rule.Description)
		}
		unique[rule.Description] = struct{}{}
	}

	return auto, manual
}
