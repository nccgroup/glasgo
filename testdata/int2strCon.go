package main

func retInt() int {
    return 90;
}

func stringCon() string {
	var a int;
	a = 123;
	
	// bad
	b := string(a);

	// bad
	b = string(123)

	// bad
	c := string(80)

	// bad
	c = string(retInt())

	b = c

	return b
}
