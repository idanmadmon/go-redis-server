package cmd

func exitOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
