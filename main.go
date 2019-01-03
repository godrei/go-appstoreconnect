package main

import (
	"fmt"
	"os"

	"github.com/bitrise-tools/xcode-project/pretty"
	"github.com/godrei/go-appstoreconnect/appstoreconnect"
)

func main() {
	// log.SetEnableDebugLog(true)

	// Authentication
	issuer := os.Getenv("ISSUER")
	kid := os.Getenv("PRIVATE_KEY_ID")
	signingKeyPth := os.Getenv("PRIVATE_KEY_PATH")

	c, err := appstoreconnect.NewClient(kid, issuer, signingKeyPth)
	if err != nil {
		panic(err)
	}

	// Search for app by bundle id
	r, _, err := c.TestFlight.Apps(&appstoreconnect.AppsOptions{
		BundleIDFilter: "com.bitrise.Bitrise-iTunesConnectBetaTest",
		Limit:          1,
	})
	if err != nil {
		panic(err)
	}

	app := r.Apps[0]
	fmt.Println("app:")
	fmt.Println(pretty.Object(app))
	fmt.Println("---")

	// List all builds
	opt := &appstoreconnect.BuildsOptions{
		AppFilter: app.ID,
		Limit:     20,
	}
	var allBuilds []appstoreconnect.Build
	var i int
	for {
		resp, _, err := c.TestFlight.Builds(opt)
		if err != nil {
			panic(err)
		}
		allBuilds = append(allBuilds, resp.Builds...)

		fmt.Printf("%d. turn\n", i)
		fmt.Println(pretty.Object(resp.Links))

		if resp.Links.Next == "" {
			break
		}

		opt.Next = resp.Links.Next
		i++
	}

	build := allBuilds[0]
	fmt.Printf("builds: %d\n", len(allBuilds))
	fmt.Println(pretty.Object(build))
	fmt.Println("---")

	// Submit the last build for Beta Review
	resp, err := c.TestFlight.BetaAppReviewSubmission(build.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(pretty.Object(resp))
}
