package main

import(
	"errors"
)

// this function returns an error
func retError1(a int) error {
	return errors.New("error 1"); 
}

func retError2(a, b int) (int, error) {
	if(a < b) {
		return a, nil
	}
	return 0, errors.New("error 2");
}

func retError3(a, b int) (int, error, int) {
	if(a < b) {
		return a, nil, b;
	}
	if(a == b) {
		return b, nil, a;
	}
	return 0, errors.New("error 3"), 0;
}

func retError4(a, b int) (int, int, error) {
	return a, b, errors.New("error 4");
}

func errorsTest() int {
	var err error;
	var a, b int;

	// good
	err = retError1(1);
	
	// bad
	retError1(1);

	// good
	a, err = retError2(0,1);

	// bad
	retError2(0,1);

	// good
	a, err, b = retError3(0,1);

	// bad
	a, _, b = retError3(0, 1);

	// good
	a, b, err = retError4(0,1);

	// bad
	a, b, _ = retError4(0,1);

	if err != nil {
		return 1
	} else {
		return a + b
	} 

}
