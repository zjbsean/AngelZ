package main

import "tests/reflect-test/testlogic"

func makeTeacher(name string) *testlogic.Teacher {
	te := &testlogic.Teacher{
		Name: name,
	}
	return te
}

func Test1() {
	te := makeTeacher("newyear")
	var tea testlogic.Speaker = te
	testlogic.TestRef(tea)
}

func Test2() {
	te := makeTeacher("newyear")
	var tea testlogic.Speaker = te
	testlogic.TestRef2(tea)
}

func Test3() {
	te := makeTeacher("newyear")
	testlogic.TestRef3(te)
}

func Test4() {
	testlogic.TestRef4()
}

func Test5() {
	testlogic.TestRef5()
}

func Test6() {
	testlogic.TestRef6()
}

func main() {
	testlogic.MakeFuncTest()
	return
}
