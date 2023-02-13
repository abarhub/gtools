package utils

import "fmt"

func Run(args []string) error {
	if len(args) > 0 {
		cmd := args[0]
		if cmd == "copy" {
			if len(args) == 3 {
				return cmdCopy(args[1], args[2])
			} else {
				return fmt.Errorf("command %v  need 2 arguments", cmd)
			}
		} else {
			return fmt.Errorf("command %v invalid", cmd)
		}
	} else {
		return fmt.Errorf("syntax CMD param1 param2 ...")
	}

}

func cmdCopy(src string, dest string) error {
	return CopyDir(src, dest)
}
