package utils

func Run() error {
	return CopyDir("./testdir/test1", "./testdir/test_out")
}
